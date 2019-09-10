package server

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/server/errors"
	"bitbucket.org/zlacki/rscgo/pkg/server/packets"
)

//Client Represents a single connecting client.
type Client struct {
	ip                               string
	uID                              uint8
	Index                            int
	isaacStream                      *IsaacStream
	Kill                             chan struct{}
	networkingGroup                  sync.WaitGroup
	player                           *entity.Player
	socket                           net.Conn
	incomingPackets, outgoingPackets chan *packets.Packet
	destroying, reconnecting         bool
	buffer                           []byte
}

//StartReader Starts the clients socket reader goroutine.  Takes a waitgroup as an argument to facilitate synchronous destruction.
func (c *Client) StartReader() {
	defer c.networkingGroup.Done()
	// 50ms for 20pps per client--is this too much?  Practically I don't think we need more than maybe 10.
	for range time.Tick(50 * time.Millisecond) {
		select {
		default:
			p, err := c.ReadPacket()
			if err != nil {
				if err, ok := err.(errors.NetError); ok && err.Error() != "Connection closed." {
					LogWarning.Printf("Rejected Packet from: %s\n", c)
					LogWarning.Println(err)
				}
				c.Destroy()
				return
			}
			c.incomingPackets <- p
		case <-c.Kill:
			return
		}
	}
}

//StartWriter Starts the clients socket writer goroutine.  Takes a waitgroup as an argument to facilitate synchronous destruction.
func (c *Client) StartWriter() {
	defer c.networkingGroup.Done()
	// 50ms for 20pps per client--is this too much?  Practically I don't think we need more than maybe 10.
	for range time.Tick(50 * time.Millisecond) {
		select {
		case p := <-c.outgoingPackets:
			if p == nil {
				return
			}
			c.WritePacket(p)
		case <-c.Kill:
			return
		}
	}
}

//Destroy Wrapper around Client.destroy to prevent multiple channel closes causing a panic.
func (c *Client) Destroy() {
	if !c.destroying {
		close(c.Kill)
		c.destroying = true
	}
}

//destroy Safely tears down a client, saves it to the database, and removes it from server-wide collections.
func (c *Client) destroy() {
	// Wait for network goroutines to finish.
	c.networkingGroup.Wait()
	c.player.Connected = false
	close(c.outgoingPackets)
	close(c.incomingPackets)
	c.buffer = []byte{} // try to collect this early it's 5KB
	if err := c.socket.Close(); err != nil {
		LogError.Println("Couldn't close socket:", err)
	}
	if _, ok := ClientsIdx[c.Index]; ok {
		delete(ClientsIdx, c.Index)
	}
	if _, ok := Clients[c.player.UserBase37]; ok {
		// Always try to launch I/O-heavy functions in their own goroutine.
		// Goroutines are light-weight and made for this kind of thing.
		go c.Save()
		entity.RemovePlayer(c.player)
		c.player.TransAttrs["plrremove"] = true
		BroadcastLogin(c.player.UserBase37, false)
		delete(Clients, c.player.UserBase37)
		LogInfo.Printf("Unregistered: %v\n", c)
	}
}

//ResetUpdateFlags Resets the players movement updating synchronization variables.
func (c *Client) ResetUpdateFlags() {
	c.player.TransAttrs["plrremove"] = false
	c.player.TransAttrs["plrmoved"] = false
	c.player.TransAttrs["plrchanged"] = false
	c.player.TransAttrs["plrself"] = true
}

//UpdatePositions Updates the client about entities in it's view-area (16x16 tiles in the game world surrounding the player).  Should be run every game engine tick.
func (c *Client) UpdatePositions() {
	if c.player.Location.Equals(entity.DeathSpot) {
		return
	}
	var localPlayers []*entity.Player
	var localAppearances []*entity.Player
	var removingPlayers []*entity.Player
	var localObjects []*entity.Object
	var removingObjects []*entity.Object
	for _, r := range entity.SurroundingRegions(c.player.X, c.player.Y) {
		for _, p := range r.Players.List {
			if p, ok := p.(*entity.Player); ok && p.Index != c.Index {
				if c.player.LongestDelta(p.Location) <= 15 {
					if !c.player.LocalPlayers.Contains(p) {
						localPlayers = append(localPlayers, p)
					}
				} else {
					if c.player.LocalPlayers.Contains(p) {
						removingPlayers = append(removingPlayers, p)
					}
				}
			}
		}
		for _, o := range r.Objects.List {
			if o, ok := o.(*entity.Object); ok {
				if c.player.LongestDelta(o.Location) <= 20 {
					if !c.player.LocalObjects.Contains(o) {
						localObjects = append(localObjects, o)
					}
				} else {
					if c.player.LocalObjects.Contains(o) {
						removingObjects = append(removingObjects, o)
					}
				}
			}
		}
	}
	// TODO: Clean up appearance list code.
	//	for _, index := range c.player.Appearances {
	//		if v, ok := ClientsIdx[index]; ok {
	//			localAppearances = append(localAppearances, v.player)
	//		}
	//	}
	localAppearances = append(localAppearances, localPlayers...)
	//	c.player.Appearances = c.player.Appearances[:0]
	// POSITIONS BEFORE EVERYTHING ELSE.
	if positions := packets.PlayerPositions(c.player, localPlayers, removingPlayers); positions != nil {
		c.outgoingPackets <- positions
	}
	if appearances := packets.PlayerAppearances(c.player, localAppearances); appearances != nil {
		c.outgoingPackets <- appearances
	}
	if objectUpdates := packets.ObjectLocations(c.player, localObjects, removingObjects); objectUpdates != nil {
		c.outgoingPackets <- objectUpdates
	}
	if boundaryUpdates := packets.BoundaryLocations(c.player, localObjects, removingObjects); boundaryUpdates != nil {
		c.outgoingPackets <- boundaryUpdates
	}
}

//StartNetworking Starts up 3 new goroutines; one for reading incoming data from the socket, one for writing outgoing data to the socket, and one for client state updates and parsing plus handling incoming packets.  When the clients kill signal is sent through the kill channel, the state update and packet handling goroutine will wait for both the reader and writer goroutines to complete their operations before unregistering the client.
func (c *Client) StartNetworking() {
	c.networkingGroup.Add(2)
	go c.StartReader()
	go c.StartWriter()
	go func() {
		defer c.destroy()
		for {
			select {
			case p := <-c.incomingPackets:
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

func (c *Client) sendLoginResponse(i byte) {
	c.outgoingPackets <- packets.LoginResponse(int(i))
	if i != 0 {
		LogInfo.Printf("Denied Client[%v]: {ip:'%v', username:'%v', Response='%v'}\n", c.Index, c.ip, c.player.Username, i)
		c.Destroy()
	} else {
		LogInfo.Printf("Registered: %v\n", c)
		entity.GetRegionFromLocation(c.player.Location).Players.Add(c.player)
		c.player.TransAttrs["plrchanged"] = true
		c.player.Connected = true
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
		c.outgoingPackets <- packets.PlayerInfo(c.player)
		c.outgoingPackets <- packets.PlayerStats(c.player)
		c.outgoingPackets <- packets.EquipmentStats(c.player)
		c.outgoingPackets <- packets.FightMode(c.player)
		c.outgoingPackets <- packets.FriendList(c.player)
		c.outgoingPackets <- packets.IgnoreList(c.player)
		c.outgoingPackets <- packets.ClientSettings(c.player)
		c.outgoingPackets <- packets.Fatigue(c.player)
		c.outgoingPackets <- packets.WelcomeMessage
		c.outgoingPackets <- packets.ServerInfo(len(Clients))
		c.outgoingPackets <- packets.LoginBox(0, c.ip)
		BroadcastLogin(c.player.UserBase37, true)
	}
}

//HandleLogin This method will block until a byte is sent down the reply channel with the login response to send to the client.
func (c *Client) HandleLogin(reply chan byte) {
	defer close(reply)
	select {
	case r := <-reply:
		c.sendLoginResponse(r)
		return
	case <-time.After(time.Second * 10):
		c.sendLoginResponse(8)
		return
	}
}

//IP Parses the players remote IP address and returns it as a go string.  TODO: Should I remove this?
func (c *Client) IP() string {
	return strings.Split(c.socket.RemoteAddr().String(), ":")[0]
}

//NewClient Creates a new instance of a Client, launches goroutines to handle I/O for it, and returns a reference to it.
func NewClient(socket net.Conn) *Client {
	c := &Client{socket: socket, incomingPackets: make(chan *packets.Packet, 20), outgoingPackets: make(chan *packets.Packet, 20), Index: -1, Kill: make(chan struct{}), player: entity.NewPlayer(), buffer: make([]byte, 5000), ip: strings.Split(socket.RemoteAddr().String(), ":")[0]}
	for lastIdx := 0; lastIdx < 2048; lastIdx++ {
		if _, ok := ClientsIdx[lastIdx]; !ok {
			c.Index = lastIdx
			break
		}
	}
	c.StartNetworking()
	return c
}

//String Returns a string populated with some of the more identifying fields from the receiver Client.
func (c *Client) String() string {
	return fmt.Sprintf("Client[%v] {username:'%v', ip:'%v'}", c.Index, c.player.Username, c.ip)
}
