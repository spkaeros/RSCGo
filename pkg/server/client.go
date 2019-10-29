package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/clients"
	"bitbucket.org/zlacki/rscgo/pkg/server/db"
	"bitbucket.org/zlacki/rscgo/pkg/server/errors"
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/packetbuilders"
	"bitbucket.org/zlacki/rscgo/pkg/server/packethandlers"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
	"fmt"
	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

//Client Represents a single connecting client.
type Client struct {
	Kill                             chan struct{}
	player                           *world.Player
	IncomingPackets, OutgoingPackets chan *packetbuilders.Packet
	PacketData                       []byte
	Socket                           net.Conn
}

func (c *Client) TypeName() string {
	return "Client"
}

func (c *Client) Equals(c1 objects.Object) bool {
	if c1, ok := c1.(*Client); ok {
		return c.Player().Index == c1.Player().Index && c1.Player().UserBase37 == c1.Player().UserBase37
	}

	return false
}

func (c *Client) Copy() objects.Object {
	return c
}

func (c *Client) BinaryOp(op token.Token, rhs objects.Object) (objects.Object, error) {
	return nil, objects.ErrInvalidOperator
}

func (c *Client) IsFalsy() bool {
	return !c.Player().TransAttrs.VarBool("connected", false)
}

func (c *Client) IndexGet(index objects.Object) (objects.Object, error) {
	switch index := index.(type) {
	case *objects.String:
		switch index.Value {
		case "index":
			return &objects.Int{Value: int64(c.player.Index)}, nil
		case "x":
			return &objects.Int{Value: int64(c.player.X.Load())}, nil
		case "y":
			return &objects.Int{Value: int64(c.player.Y.Load())}, nil
		case "username":
			return &objects.String{Value: c.player.Username}, nil
		case "level":
			return &objects.Int{Value: int64(c.player.Y.Load())}, nil
		case "curAttack":
			return &objects.Int{Value: int64(c.player.Skillset.Current[0])}, nil
		case "maxAttack":
			return &objects.Int{Value: int64(c.player.Skillset.Maximum[0])}, nil
		case "curDefense":
			return &objects.Int{Value: int64(c.player.Skillset.Current[1])}, nil
		case "maxDefense":
			return &objects.Int{Value: int64(c.player.Skillset.Maximum[1])}, nil
		case "curStrength":
			return &objects.Int{Value: int64(c.player.Skillset.Current[2])}, nil
		case "maxStrength":
			return &objects.Int{Value: int64(c.player.Skillset.Maximum[2])}, nil
		case "curHits":
			return &objects.Int{Value: int64(c.player.Skillset.Current[3])}, nil
		case "maxHits":
			return &objects.Int{Value: int64(c.player.Skillset.Maximum[3])}, nil
		case "curRanged":
			return &objects.Int{Value: int64(c.player.Skillset.Current[4])}, nil
		case "maxRanged":
			return &objects.Int{Value: int64(c.player.Skillset.Maximum[4])}, nil
		case "curPrayer":
			return &objects.Int{Value: int64(c.player.Skillset.Current[5])}, nil
		case "maxPrayer":
			return &objects.Int{Value: int64(c.player.Skillset.Maximum[5])}, nil
		case "curMagic":
			return &objects.Int{Value: int64(c.player.Skillset.Current[6])}, nil
		case "maxMagic":
			return &objects.Int{Value: int64(c.player.Skillset.Maximum[6])}, nil
		case "curCooking":
			return &objects.Int{Value: int64(c.player.Skillset.Current[7])}, nil
		case "maxCooking":
			return &objects.Int{Value: int64(c.player.Skillset.Maximum[7])}, nil
//		case "cur":
//			return &objects.Int{Value: int64(c.player.Skillset.Current[6])}, nil
//		case "max":
//			return &objects.Int{Value: int64(c.player.Skillset.Maximum[6])}, nil
		case "curSkill":
			return &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					if len(args) < 1 {
						return nil, objects.ErrWrongNumArguments
					}
					index, ok := objects.ToInt(args[0])
					if !ok {
						return nil, objects.ErrInvalidArgumentType{
							Name:     "index",
							Expected: "int",
							Found:    args[0].TypeName(),
						}
					}

					return &objects.Int{Value: int64(c.player.Skillset.Current[index])}, nil
				},
			}, nil
		case "maxSkill":
			return &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					if len(args) < 1 {
						return nil, objects.ErrWrongNumArguments
					}
					index, ok := objects.ToInt(args[0])
					if !ok {
						return nil, objects.ErrInvalidArgumentType{
							Name:     "index",
							Expected: "int",
							Found:    args[0].TypeName(),
						}
					}

					return &objects.Int{Value: int64(c.player.Skillset.Maximum[index])}, nil
				},
			}, nil
		case "teleport":
			return &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					ret = objects.UndefinedValue
					if len(args) != 2 {
						c.Message("teleport(x,y): Invalid argument count provided")
						return nil, objects.ErrWrongNumArguments
					}
					x, ok := objects.ToInt(args[0])
					if !ok {
						c.Message("teleport(x,y): Invalid argument type provided")
						return nil, objects.ErrInvalidArgumentType{
							Name:     "x",
							Expected: "int",
							Found:    args[0].TypeName(),
						}
					}
					y, ok := objects.ToInt(args[1])
					if !ok {
						c.Message("teleport(x,y): Invalid argument type provided")
						return nil, objects.ErrInvalidArgumentType{
							Name:     "y",
							Expected: "int",
							Found:    args[1].TypeName(),
						}
					}
					c.Player().Teleport(x, y)
					return
				},
			}, nil
		case "message":
			return &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					ret = objects.UndefinedValue

					message, ok := objects.ToString(args[0])
					if !ok {
						message = args[0].String()
					}

					c.Message(message)
					return
				},
			}, nil
		case "goUp":
			return &objects.UserFunction{
				Name: "goUp",
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					ret = objects.UndefinedValue
					if nextLocation := c.Player().Above(); !nextLocation.Equals(c.Player().Location) {
						c.Player().ResetPath()
						c.Player().SetLocation(&nextLocation)
						c.UpdatePlane()
					}
					return
				},
			}, nil
		case "goDown":
			return &objects.UserFunction{
				Name: "goDown",
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					ret = objects.UndefinedValue
					if nextLocation := c.Player().Below(); !nextLocation.Equals(c.Player().Location) {
						c.Player().ResetPath()
						c.Player().SetLocation(&nextLocation)
						c.UpdatePlane()
					}
					return
				},
			}, nil
		}
	}
	return nil, objects.ErrInvalidIndexType
}

//Player returns the scene player that this client represents
func (c *Client) Player() *world.Player {
	return c.player
}

//SendPacket Queue a packet for sending to the client.
func (c *Client) SendPacket(p *packetbuilders.Packet) {
	c.OutgoingPackets <- p
}

//Message Builds a new game packet to display a message in the client chat box with msg as its contents, and queues it in the outgoing packet queue.
func (c *Client) Message(msg string) {
	c.SendPacket(packetbuilders.ServerMessage(msg))
}

//UpdateStat Builds and queues for sending a new packet containing our players stat information for given skill ID
func (c *Client) UpdateStat(id int) {
	c.SendPacket(packetbuilders.PlayerStat(c.player, id))
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
	c.SendPacket(packetbuilders.ChangeAppearance)
}

//Teleport Moves the client's player to x,y in the game world, and sends a teleport bubble animation packet to all of the view-area client.
func (c *Client) Teleport(x, y int) {
	if !world.WithinWorld(x, y) {
		return
	}
	for _, nearbyPlayer := range c.player.NearbyPlayers() {
		if c1, ok := clients.FromIndex(nearbyPlayer.Index); ok {
			c1.TeleBubble(int(c.player.X.Load()-nearbyPlayer.X.Load()), int(c.player.Y.Load()-nearbyPlayer.Y.Load()))
		}
	}
	c.TeleBubble(0, 0)
	c.player.Teleport(x, y)
}

//startReader Starts the client Socket reader goroutine.  Takes a waitgroup as an argument to facilitate synchronous destruction.
func (c *Client) startReader() {
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
			if !c.player.TransAttrs.VarBool("connected", false) && p.Opcode != 32 && p.Opcode != 0 && p.Opcode != 2 && p.Opcode != 220 {
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
	for {
		select {
		case p := <-c.OutgoingPackets:
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
	if !c.player.TransAttrs.VarBool("destroying", false) {
		c.player.TransAttrs.SetVar("destroying", true)
		close(c.Kill)
	}
}

//destroy Safely tears down a client, saves it to the database, and removes it from server-wide clients.
func (c *Client) destroy(wg *sync.WaitGroup) {
	// Wait for network goroutines to finish.
	(*wg).Wait()
	c.player.TransAttrs.UnsetVar("connected")
	close(c.OutgoingPackets)
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
}

//ResetUpdateFlags Resets the players movement updating synchronization variables.
func (c *Client) ResetUpdateFlags() {
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
func (c *Client) HandlePacket(p *packetbuilders.Packet) {
	handler := packethandlers.GetPacketHandler(p.Opcode)
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
	c.player.TransAttrs.SetVar("changed", true)
	c.player.TransAttrs.SetVar("connected", true)
	for i := 0; i < 18; i++ {
		level := 1
		exp := 0
		if i == 3 {
			level = 10
			exp = 1154
		}
		c.player.Skillset.Current[i] = level
		c.player.Skillset.Maximum[i] = level
		c.player.Skillset.Experience[i] = exp
	}
	c.SendPacket(packetbuilders.PlaneInfo(c.player))
	if !c.player.Reconnecting() {
		// Reconnecting implies that the client has all of this data already, so as an optimization, we don't send it again
		c.SendPacket(packetbuilders.PlayerStats(c.player))
		c.SendPacket(packetbuilders.EquipmentStats(c.player))
		c.SendPacket(packetbuilders.Fatigue(c.player))
		c.SendPacket(packetbuilders.InventoryItems(c.player))
		// TODO: Not canonical RSC, but definitely good QoL update...
		//  c.SendPacket(packetbuilders.FightMode(c.player)
		c.SendPacket(packetbuilders.FriendList(c.player))
		c.SendPacket(packetbuilders.IgnoreList(c.player))
		c.SendPacket(packetbuilders.ClientSettings(c.player))
		c.SendPacket(packetbuilders.PrivacySettings(c.player))
		c.SendPacket(packetbuilders.WelcomeMessage)
		c.SendPacket(packetbuilders.LoginBox(0, c.player.IP))
	}
	clients.BroadcastLogin(c.player, true)
	if c.player.Attributes.VarBool("first_login", true) {
		c.player.Attributes.SetVar("first_login", false)
		c.OpenAppearanceChangePanel()
	}
}

//HandleLogin This method will block until a byte is sent down the reply channel with the login response to send to the client, or if this doesn't occur, it will timeout after 10 seconds.
func (c *Client) HandleLogin(reply chan byte) {
	defer close(reply)
	select {
	case r := <-reply:
		c.SendPacket(packetbuilders.LoginResponse(int(r)))
		if r == 0 || r == 1 || r == 25 || r == 24 {
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
func NewClient(socket net.Conn) *Client {
	c := &Client{Socket: socket, IncomingPackets: make(chan *packetbuilders.Packet, 20), OutgoingPackets: make(chan *packetbuilders.Packet, 20), Kill: make(chan struct{}), player: world.NewPlayer(clients.NextIndex(), strings.Split(socket.RemoteAddr().String(), ":")[0])}
	c.StartNetworking()
	return c
}

//String Returns a string populated with some of the more identifying fields from the receiver Client.
func (c *Client) String() string {
	return fmt.Sprintf("Client[%v] {username:'%v', IP:'%v'}", c.player.Index, c.player.Username, c.player.IP)
}

//Write Writes data to the client's Socket from `b`.  Returns the length of the written bytes.
func (c *Client) Write(b []byte) int {
	if !c.player.Websocket {
		l, err := c.Socket.Write(b)
		if err != nil {
			log.Error.Println("Could not write to client Socket.", err)
			c.Destroy()
		} else if l != len(b) {
			// Possibly non-fatal?
			log.Error.Printf("Wrong number of bytes written to Client Socket.  Expected %d, got %d.\n", len(b), l)
		}
		return l
	}
	w := wsutil.NewWriter(c.Socket, ws.StateServerSide, ws.OpBinary)
	l, err := w.Write(b)
	if err != nil {
		log.Error.Println("Could not write to client Socket.", err)
		c.Destroy()
	} else if l != len(b) {
		// Possibly non-fatal?
		log.Error.Printf("Wrong number of bytes written to Client Socket.  Expected %d, got %d.\n", len(b), l)
	}
	if err := w.Flush(); err != nil {
		log.Warning.Println("Error writing to Websocket:", err)
	}
	return l
}

//Read Reads data off of the client's Socket into 'dst'.  Returns length read into dst upon success.  Otherwise, returns -1 with a meaningful error message.
func (c *Client) Read(dst []byte) (int, error) {
	err := c.Socket.SetReadDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		return -1, errors.ConnDeadline
	}

	if len(c.PacketData) >= len(dst) {
		// If we have enough data to fill dst, fill it, stash the remaining leftovers
		copy(dst, c.PacketData)
		if len(c.PacketData) == len(dst) {
			c.PacketData = c.PacketData[:0]
		} else {
			c.PacketData = c.PacketData[len(dst):]
		}
		return len(dst), nil
	}

	var data []byte
	if c.player.Websocket {
		data, _, err = wsutil.ReadData(c.Socket, ws.StateServerSide)
	} else {
		_, err = c.Socket.Read(dst)
		data = dst
	}

	if err != nil {
		if err == io.EOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
			return -1, errors.ConnClosed
		} else if e, ok := err.(net.Error); ok && e.Timeout() {
			return -1, errors.ConnTimedOut
		}
		return -1, err
	}

	if len(c.PacketData) > 0 {
		// unstash extra data
		if c.player.Websocket {
			data = append(c.PacketData, data...)
		} else {
			dst = append(c.PacketData, dst...)
		}
		c.PacketData = c.PacketData[:0]
	}

	if c.player.Websocket {
		copy(dst, data)

		if len(data) > len(dst) {
			// stash extra data
			c.PacketData = data[len(dst):]
		}
	}
	return len(dst), nil
}

//ReadPacket Attempts to read and parse the next 3 bytes of incoming data for the 16-bit length and 8-bit opcode of the next packet frame the client is sending us.
func (c *Client) ReadPacket() (*packetbuilders.Packet, error) {
	// TODO: Is allocation overhead more expensive than mutex locks?
	header := make([]byte, 2)
	if l, err := c.Read(header); err != nil {
		return nil, err
	} else if l < 2 {
		return nil, errors.NewNetworkError("SHORT_DATA")
	}
	length := int(header[0])
	if length >= 160 {
		length = (length-160)*256 + int(header[1])
	} else {
		length--
	}

	if length >= 5000 || length < 0 {
		log.Suspicious.Printf("Invalid packet length from [%v]: %d\n", c, length)
		log.Warning.Printf("Packet from [%v] length out of bounds; got %d, expected between 4 and 5000\n", c, length+3)
		return nil, errors.NewNetworkError("Packet length out of bounds; must be between 4 and 5000.")
	}

	payload := make([]byte, length)

	if l, err := c.Read(payload); err != nil {
		return nil, err
	} else if l < length {
		return nil, errors.NewNetworkError("SHORT_DATA")
	}

	if length < 160 {
		payload = append(payload, header[1])
	}

	return packetbuilders.NewPacket(payload[0], payload[1:]), nil
}

//WritePacket This is a method to send a packet to the client.  If this is a bare packet, the packet payload will
// be written as-is.  If this is not a bare packet, the packet will have the first 3 bytes changed to the
// appropriate values for the client to parse the length and opcode for this packet.
func (c *Client) WritePacket(p packetbuilders.Packet) {
	var buf []byte
	if !p.Bare {
		if frameLength := len(p.Payload); frameLength >= 160 {
			buf = append(buf, byte(160+frameLength/256))
			buf = append(buf, byte(frameLength))
		} else {
			buf = append(buf, byte(frameLength))
			buf = append(buf, p.Payload[frameLength-1])
			p.Payload = p.Payload[:frameLength-1]
		}
	}
	buf = append(buf, p.Payload...)

	c.Write(buf)
}
