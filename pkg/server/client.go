package server

import (
	"fmt"
	"github.com/gobwas/ws/wsutil"
	"github.com/spkaeros/rscgo/pkg/server/clients"
	"github.com/spkaeros/rscgo/pkg/server/db"
	"github.com/spkaeros/rscgo/pkg/server/errors"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"github.com/spkaeros/rscgo/pkg/server/packetbuilders"
	"github.com/spkaeros/rscgo/pkg/server/packethandlers"
	"github.com/spkaeros/rscgo/pkg/server/script"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

//Client Represents a single connecting client.
type Client struct {
	Kill              chan struct{}
	player            *world.Player
	IncomingPackets   chan *packet.Packet
	CacheBuffer       []byte
	Socket            net.Conn
	DataBuffer        []byte
	DataLock          sync.RWMutex
	destroyer, killer sync.Once
}

//Player returns the scene player that this client represents
func (c *Client) Player() *world.Player {
	return c.player
}

//SendPacket Queue a packet for sending to the client.
func (c *Client) SendPacket(p *packet.Packet) {
	c.Player().OutgoingPackets <- p
}

//Message Builds a new game packet to display a message in the client chat box with msg as its contents, and queues it in the outgoing packet queue.
func (c *Client) Message(msg string) {
	c.SendPacket(packetbuilders.ServerMessage(msg))
}

//UpdateStat Builds and queues for sending a new packet containing our players stat information for given skill ID
func (c *Client) UpdateStat(id int) {
	c.SendPacket(packetbuilders.PlayerStat(c.player, id))
}

func (c *Client) SendStats() {
	c.SendPacket(packetbuilders.PlayerStats(c.player))
}

func (c *Client) SendInventory() {
	c.SendPacket(packetbuilders.InventoryItems(c.player))
}

//TeleBubble Queues a new packet to create a teleport bubble at the given offsets relative to our player.
func (c *Client) TeleBubble(diffX, diffY int) {
	c.SendPacket(packetbuilders.TeleBubble(diffX, diffY))
}

//UpdatePlane Updates the client about the plane that its player is on.
func (c *Client) UpdatePlane() {
	c.SendPacket(packetbuilders.PlaneInfo(c.player))
}

//TradeOpen Opens a trade window for this client.  Must have this client's player's trade target set, or will cause a disconnect.
func (c *Client) TradeOpen() {
	c.SendPacket(packetbuilders.TradeOpen(c.player))
}

//OpenAppearanceChangePanel Sends a packet to open the player appearance changing panel on the client.
func (c *Client) OpenAppearanceChangePanel() {
	c.player.AddState(world.MSChangingAppearance)
	c.SendPacket(packetbuilders.ChangeAppearance)
}

//Teleport Moves the client's player to x,y in the game world, and sends a teleport bubble animation packet to all of the view-area client.
func (c *Client) Teleport(x, y int) {
	if !world.WithinWorld(x, y) {
		return
	}
	for _, nearbyPlayer := range c.player.NearbyPlayers() {
		nearbyPlayer.SendPacket(packetbuilders.TeleBubble(c.player.X()-nearbyPlayer.X(), c.player.Y()-nearbyPlayer.Y()))
	}
	c.TeleBubble(0, 0)
	oldPlane := c.player.Plane()
	c.player.SetLocation(world.NewLocation(x, y), true)
	if c.player.Plane() != oldPlane {
		c.UpdatePlane()
	}
}

//startReader Starts the client Socket reader goroutine.  Takes a waitgroup as an argument to facilitate synchronous destruction.
func (c *Client) startReader() {
	defer c.Destroy()
	for {
		select {
		default:
			p, err := c.ReadPacket()
			if err != nil {
				if err, ok := err.(errors.NetError); ok && err.Error() != "Connection closed." && err.Error() != "Connection timed out." {
					if err.Error() != "SHORT_DATA" {
						log.Warning.Printf("Rejected Packet from: %s\n", c)
						log.Warning.Println(err)
					}
					continue
				}
				c.Destroy()
				return
			}
			if !c.player.Connected() && p.Opcode != 32 && p.Opcode != 0 && p.Opcode != 2 && p.Opcode != 220 {
				log.Warning.Printf("Unauthorized packet[opcode:%v,len:%v] rejected from: %v\n", p.Opcode, len(p.Payload), c)
				c.Destroy()
				return
			}
			c.IncomingPackets <- p
		case <-c.Kill:
			return
		}
	}
}

//startWriter Starts the client Socket writer goroutine.
func (c *Client) startWriter() {
	defer c.Destroy()
	for {
		select {
		case p := <-c.Player().OutgoingPackets:
			if p == nil {
				return
			}
			c.WritePacket(*p)
		case <-c.Kill:
			return
		}
	}
}

//Destroy Wrapper around Client.destroy to prevent multiple channel closes causing a panic.
func (c *Client) Destroy() {
	c.killer.Do(func() {
		close(c.Kill)
	})
}

//destroy Safely tears down a client, saves it to the database, and removes it from server-wide clients.
func (c *Client) destroy(wg *sync.WaitGroup) {
	// Wait for network goroutines to finish.
	c.destroyer.Do(func() {
		(*wg).Wait()
		c.player.TransAttrs.UnsetVar("connected")
		close(c.Player().OutgoingPackets)
		close(c.Player().OptionMenuC)
		close(c.IncomingPackets)
		if err := c.Socket.Close(); err != nil {
			log.Error.Println("Couldn't close Socket:", err)
		}
		if _, ok := clients.FromUserHash(c.player.UserBase37); ok {
			// Always try to launch I/O-heavy functions in their own goroutine.
			// Goroutines are light-weight and made for this kind of thing.
			go db.SavePlayer(c.player)
			world.RemovePlayer(c.player)
			c.player.TransAttrs.SetVar("remove", true)
			clients.BroadcastLogin(c.player, false)
			clients.Remove(c)
			log.Info.Printf("Unregistered: %v\n", c)
		}
	})
}

//ResetUpdateFlags Resets the players movement updating synchronization variables.
func (c *Client) ResetUpdateFlags() {
	// TODO: Is shouldReset semantically correct?
	c.player.TransAttrs.SetVar("self", true)
	c.player.TransAttrs.UnsetVar("remove")
	c.player.TransAttrs.UnsetVar("moved")
	c.player.TransAttrs.UnsetVar("changed")
}

//UpdatePositions Updates the client about entities in it's view-area (16x16 tiles in the game world surrounding the player).  Should be run every game engine tick.
func (c *Client) UpdatePositions() {
	// Everything is updated relative to our player's position, so player position packet comes first
	if positions := packetbuilders.PlayerPositions(c.player); positions != nil {
		c.SendPacket(positions)
	}
	if appearances := packetbuilders.PlayerAppearances(c.player); appearances != nil {
		c.SendPacket(appearances)
	}
	if npcUpdates := packetbuilders.NPCPositions(c.player); npcUpdates != nil {
		c.SendPacket(npcUpdates)
	}
	/*
		if npcAppearances := packetbuilders.NpcAppearances(c.player); npcAppearances != nil {
			c.SendPacket(npcAppearances)
		}
	*/
	if itemUpdates := packetbuilders.ItemLocations(c.player); itemUpdates != nil {
		c.SendPacket(itemUpdates)
	}
	if objectUpdates := packetbuilders.ObjectLocations(c.player); objectUpdates != nil {
		c.SendPacket(objectUpdates)
	}
	if boundaryUpdates := packetbuilders.BoundaryLocations(c.player); boundaryUpdates != nil {
		c.SendPacket(boundaryUpdates)
	}
}

//StartNetworking Starts up 3 new goroutines; one for reading incoming data from the Socket, one for writing outgoing data to the Socket, and one for client state updates and parsing plus handling incoming packetbuilders.  When the client kill signal is sent through the kill channel, the state update and packet handling goroutine will wait for both the reader and writer goroutines to complete their operations before unregistering the client.
func (c *Client) StartNetworking() {
	var nwg sync.WaitGroup
	nwg.Add(2)
	go func() {
		defer nwg.Done()
		c.startReader()
	}()
	go func() {
		defer nwg.Done()
		c.startWriter()
	}()
	go func() {
		defer c.destroy(&nwg)
		for {
			select {
			case p := <-c.IncomingPackets:
				if p == nil {
					return
				}
				c.HandlePacket(p)
			case <-c.Kill:
				return
			}
		}
	}()
}

//HandlePacket Finds the mapped handler function for the specified packet, and calls it with the specified parameters.
func (c *Client) HandlePacket(p *packet.Packet) {
	handler := packethandlers.Get(p.Opcode)
	if handler == nil {
		log.Info.Printf("Unhandled Packet: {opcode:%d; length:%d};\n", p.Opcode, len(p.Payload))
		fmt.Printf("CONTENT: %v\n", p.Payload)
		return
	}

	handler(c, p)
}

//Initialize Adds client to server's Clients list, initializes player variables and adds to world, announces login to
// other client that care, and sends all of the player and world information to the client.  Upon first login,
// sends appearance change screen to setup player's appearance.
func (c *Client) Initialize() {
	for user := range c.player.FriendList {
		if clients.ContainsHash(user) {
			c.player.FriendList[user] = true
		}
	}
	world.AddPlayer(c.player)
	c.player.Change()
	c.player.SetConnected(true)
	if c.player.Skills().Experience(world.StatHits) < 10 {
		for i := 0; i < 18; i++ {
			level := 1
			exp := 0
			if i == 3 {
				level = 10
				exp = 1154
			}
			c.player.Skills().SetCur(i, level)
			c.player.Skills().SetMax(i, level)
			c.player.Skills().SetExp(i, exp)
		}
	}
	if s := time.Until(script.UpdateTime).Seconds(); s > 0 {
		c.SendPacket(packetbuilders.SystemUpdate(int(s)))
	}
	c.SendPacket(packetbuilders.PlaneInfo(c.player))
	c.SendPacket(packetbuilders.FriendList(c.player))
	c.SendPacket(packetbuilders.IgnoreList(c.player))
	if !c.player.Reconnecting() {
		// Reconnecting implies that the client has all of this data already, so as an optimization, we don't send it again
		c.SendPacket(packetbuilders.PlayerStats(c.player))
		c.SendPacket(packetbuilders.EquipmentStats(c.player))
		c.SendPacket(packetbuilders.Fatigue(c.player))
		c.SendPacket(packetbuilders.InventoryItems(c.player))
		// TODO: Not canonical RSC, but definitely good QoL update...
		//  c.SendPacket(packetbuilders.FightMode(c.player)
		c.SendPacket(packetbuilders.ClientSettings(c.player))
		c.SendPacket(packetbuilders.PrivacySettings(c.player))
		c.SendPacket(packetbuilders.WelcomeMessage)
		t, err := time.Parse(time.ANSIC, c.player.Attributes.VarString("lastLogin", time.Time{}.Format(time.ANSIC)))
		if err != nil {
			log.Info.Println(err)
			return
		}

		days := int(time.Since(t).Hours()/24)
		c.SendPacket(packetbuilders.LoginBox(days, c.player.Attributes.VarString("lastIP", "127.0.0.1")))
	}
	clients.BroadcastLogin(c.player, true)
	if c.player.FirstLogin() {
		c.player.SetFirstLogin(false)
		c.OpenAppearanceChangePanel()
	}
}

//HandleLogin This method will block until a byte is sent down the reply channel with the login response to send to the client, or if this doesn't occur, it will timeout after 10 seconds.
func (c *Client) HandleLogin(reply chan byte) {
	isValid := func(r byte) bool {
		valid := [...]byte{0, 1, 24, 25}
		for _, i := range valid {
			if i == r {
				return true
			}
		}
		return false
	}
	defer close(reply)
	select {
	case r := <-reply:
		c.SendPacket(packetbuilders.LoginResponse(int(r)))
		if isValid(r) {
			clients.Put(c)
			log.Info.Printf("Registered: %v\n", c)
			c.Initialize()
			return
		}
		log.Info.Printf("Denied Client: {IP:'%v', username:'%v', Response='%v'}\n", c.player.IP, c.player.Username, r)
		c.Destroy()
		return
	case <-time.After(time.Second * 10):
		c.SendPacket(packetbuilders.LoginResponse(-1))
		return
	}
}

//HandleRegister This method will block until a byte is sent down the reply channel with the registration response to send to the client, or if this doesn't occur, it will timeout after 10 seconds.
func (c *Client) HandleRegister(reply chan byte) {
	defer c.Destroy()
	defer close(reply)
	select {
	case r := <-reply:
		c.SendPacket(packetbuilders.LoginResponse(int(r)))
		return
	case <-time.After(time.Second * 10):
		c.SendPacket(packetbuilders.LoginResponse(0))
		return
	}
}

//NewClient Creates a new instance of a Client, launches goroutines to handle I/O for it, and returns a reference to it.
func NewClient(socket net.Conn, ws bool) *Client {
	c := &Client{Socket: socket, IncomingPackets: make(chan *packet.Packet, 20), Kill: make(chan struct{}), DataBuffer: make([]byte, 5000)}
	c.player = world.NewPlayer(clients.NextIndex(), strings.Split(socket.RemoteAddr().String(), ":")[0])
	c.player.Websocket = ws
	c.StartNetworking()
	return c
}

//String Returns a string populated with some of the more identifying fields from the receiver Client.
func (c *Client) String() string {
	return fmt.Sprintf("Client[%v] {username:'%v', IP:'%v'}", c.player.Index, c.player.Username, c.player.IP)
}

//Write Writes data to the client's Socket from `b`.  Returns the length of the written bytes.
func (c *Client) Write(src []byte) int {
	var err error
	var dataLen int
	if c.player.Websocket {
		err = wsutil.WriteServerBinary(c.Socket, src)
		dataLen = len(src)
	} else {
		dataLen, err = c.Socket.Write(src)
	}
	if err != nil {
		log.Error.Println("Problem writing to websocket client:", err)
		c.Destroy()
		return -1
	}
	return dataLen
}

//Read Reads data off of the client's Socket into 'dst'.  Returns length read into dst upon success.  Otherwise, returns -1 with a meaningful error message.
func (c *Client) Read(dst []byte) (int, error) {
	// Set the read deadline for the socket to 10 seconds from now.
	err := c.Socket.SetReadDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		return -1, errors.ConnDeadline
	}

	expectedLen := len(dst)
	// Unstash any overflow data from previous read calls.
	cacheLen := len(c.CacheBuffer)
	if cacheLen > 0 {
		copy(dst, c.CacheBuffer)
		if cacheLen > expectedLen {
			c.CacheBuffer = c.CacheBuffer[expectedLen:]
			return expectedLen, nil
		} else {
			c.CacheBuffer = []byte{}
			if cacheLen == expectedLen {
				return expectedLen, nil
			}
		}
	}

	// Mark length of data left to read from socket after unstashing anything from the buffer
	reqDataLen := expectedLen - cacheLen

	var dataLen int
	var data []byte
	if !c.player.Websocket {
		dataLen, err = c.Socket.Read(dst[cacheLen:])
	} else {
		data, err = wsutil.ReadClientBinary(c.Socket)
		dataLen = len(data)
	}
	if err != nil {
		if err == io.EOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
			return -1, errors.ConnClosed
		} else if e, ok := err.(net.Error); ok && e.Timeout() {
			return -1, errors.ConnTimedOut
		}
		return -1, err
	}
	if c.player.Websocket {
		copy(dst[cacheLen:], data)
	}

	if dataLen < reqDataLen {
		// We didn't have enough data.  In practice, this produces an error I believe, but just in case!
		c.CacheBuffer = dst[:dataLen+cacheLen]
	} else if dataLen > reqDataLen {
		// We read too much data.  Stash what is not required.
		if c.player.Websocket {
			// Cache the recv'd data starting right after the last needed byte, next Read will unstash as if it were new data
			c.CacheBuffer = data[reqDataLen:]
		} else {
			// I don't think this can happen with TCP sockets.  We have finer control over what we read with them.
			// Just in case, I'll handle it in a semantically correct way, but I doubt it will ever run.
			c.CacheBuffer = dst[cacheLen+dataLen:]
		}
	}
	return dataLen + cacheLen, nil
}

//ReadPacket Attempts to read and parse the next 3 bytes of incoming data for the 16-bit length and 8-bit opcode of the next packet frame the client is sending us.
func (c *Client) ReadPacket() (*packet.Packet, error) {
	header := make([]byte, 2)
	if l, err := c.Read(header); err != nil {
		return nil, err
	} else if l < 2 {
		return nil, errors.NewNetworkError("SHORT_DATA")
	}
	length := int(header[0])
	bigLength := length >= 160
	if bigLength {
		// length = (length-160)*256 + int(header[1])
		length = (length-160)<<8 + int(header[1])
	} else {
		// We have the final byte of frame data already, stored at header[1]
		length--
	}

	if length+2 >= 5000 || length+2 < 2 {
		log.Suspicious.Printf("Invalid packet length from [%v]: %d\n", c, length)
		log.Warning.Printf("Packet from [%v] length out of bounds; got %d, expected between 0 and 5000\n", c, length)
		return nil, errors.NewNetworkError("Packet length out of bounds; must be between 0 and 5000.")
	}

	payload := make([]byte, length)

	if length > 0 {
		if l, err := c.Read(payload); err != nil {
			return nil, err
		} else if l < length {
			return nil, errors.NewNetworkError("SHORT_DATA")
		}
	}

	if !bigLength {
		// If the length in the packet header used 1 byte, the 2nd byte in the header is the final byte of frame data
		payload = append(payload, header[1])
	}

	return packet.NewPacket(payload[0], payload[1:]), nil
}

//WritePacket This is a method to send a packet to the client.  If this is a bare packet, the packet payload will
// be written as-is.  If this is not a bare packet, the packet will have the first 3 bytes changed to the
// appropriate values for the client to parse the length and opcode for this packet.
func (c *Client) WritePacket(p packet.Packet) {
	if p.Bare {
		c.Write(p.Payload)
		return
	}
	frameLength := len(p.Payload)
	c.DataLock.Lock()
	header := c.DataBuffer[0:2]
	defer c.DataLock.Unlock()
	if frameLength >= 160 {
		//		header[0] = byte(frameLength/256+160)
		header[0] = byte(frameLength>>8 + 160)
		header[1] = byte(frameLength)
	} else {
		header[0] = byte(frameLength)
		header[1] = p.Payload[frameLength-1]
		p.Payload = p.Payload[:frameLength-1]
	}
	c.Write(append(header, p.Payload...))
	return
}
