package client

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/collections"
	"bitbucket.org/zlacki/rscgo/pkg/server/db"
	"bitbucket.org/zlacki/rscgo/pkg/server/errors"
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/packetbuilders"
	"bitbucket.org/zlacki/rscgo/pkg/server/packethandlers"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
	"fmt"
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
		if c1, ok := collections.Clients.FromIndex(nearbyPlayer.Index); ok {
			c1.TeleBubble(int(c.player.X.Load()-nearbyPlayer.X.Load()), int(c.player.Y.Load()-nearbyPlayer.Y.Load()))
		}
	}
	c.TeleBubble(0, 0)
	c.player.Teleport(x, y)
}

//StartReader Starts the client Socket reader goroutine.  Takes a waitgroup as an argument to facilitate synchronous destruction.
func (c *Client) StartReader() {
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

//StartWriter Starts the client Socket writer goroutine.
func (c *Client) StartWriter() {
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

//destroy Safely tears down a client, saves it to the database, and removes it from server-wide collections.
func (c *Client) destroy(wg *sync.WaitGroup) {
	// Wait for network goroutines to finish.
	(*wg).Wait()
	c.player.TransAttrs.UnsetVar("connected")
	close(c.OutgoingPackets)
	close(c.IncomingPackets)
	if err := c.Socket.Close(); err != nil {
		log.Error.Println("Couldn't close Socket:", err)
	}
	if _, ok := collections.Clients.FromUserHash(c.player.UserBase37); ok {
		// Always try to launch I/O-heavy functions in their own goroutine.
		// Goroutines are light-weight and made for this kind of thing.
		go db.SavePlayer(c.player)
		world.RemovePlayer(c.player)
		c.player.TransAttrs.SetVar("remove", true)
		collections.BroadcastLogin(c.player, false)
		collections.Clients.Remove(c)
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
		c.StartReader()
	}()
	go func() {
		defer nwg.Done()
		c.StartWriter()
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
		if collections.Clients.ContainsHash(user) {
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
	collections.BroadcastLogin(c.player, true)
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
			collections.Clients.Put(c)
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
	c := &Client{Socket: socket, IncomingPackets: make(chan *packetbuilders.Packet, 20), OutgoingPackets: make(chan *packetbuilders.Packet, 20), Kill: make(chan struct{}), player: world.NewPlayer(collections.Clients.NextIndex(), strings.Split(socket.RemoteAddr().String(), ":")[0])}
	c.StartNetworking()
	return c
}

//String Returns a string populated with some of the more identifying fields from the receiver Client.
func (c *Client) String() string {
	return fmt.Sprintf("Client[%v] {username:'%v', IP:'%v'}", c.player.Index, c.player.Username, c.player.IP)
}
