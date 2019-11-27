package world

import (
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"go.uber.org/atomic"
	"strconv"
	"sync"
)

//player Represents a single player.
type Player struct {
	Username         string
	UserBase37       uint64
	Password         string
	FriendList       map[uint64]bool
	IgnoreList       []uint64
	LocalPlayers     *List
	LocalNPCs        *List
	LocalObjects     *List
	LocalItems       *List
	Updating         bool
	Appearances      []int
	DatabaseIndex    int
	Rank             int
	Appearance       AppearanceTable
	AppearanceTicket int
	KnownAppearances map[int]int
	AppearanceReq    []*Player
	AppearanceLock   sync.RWMutex
	Attributes       *AttributeList
	Items            *Inventory
	TradeOffer       *Inventory
	DistancedAction  func() bool
	ActionLock       sync.RWMutex
	OutgoingPackets chan *packet.Packet
	IP               string
	UID              uint8
	Websocket        bool
	Equips           [12]int
	*Mob
}

//String Returns a string populated with the more identifying features of this player.
func (p *Player) String() string {
	return "[" + p.Username + ", " + p.IP + "]"
}

func (p *Player) Inventory() *Inventory {
	return p.Items
}

//SetDistancedAction Queues a distanced action to run every game engine tick before path traversal, if action returns true, it will be reset.
func (p *Player) SetDistancedAction(action func() bool) {
	p.ActionLock.Lock()
	p.DistancedAction = action
	p.ActionLock.Unlock()
}

//ResetDistancedAction Clears the distanced action, if any is queued.  Should be called any time the player is deliberately performing an action.
func (p *Player) ResetDistancedAction() {
	p.ActionLock.Lock()
	p.DistancedAction = nil
	p.ActionLock.Unlock()
}

//Friends Returns true if specified username is in our friend list.
func (p *Player) Friends(other uint64) bool {
	for hash := range p.FriendList {
		if hash == other {
			return true
		}
	}
	return false
}

//Ignoring Returns true if specified username is in our ignore list.
func (p *Player) Ignoring(hash uint64) bool {
	for _, v := range p.IgnoreList {
		if v == hash {
			return true
		}
	}
	return false
}

//ChatBlocked Returns true if public chat is blocked for this player.
func (p *Player) ChatBlocked() bool {
	return p.Attributes.VarBool("chat_block", false)
}

//FriendBlocked Returns true if private chat is blocked for this player.
func (p *Player) FriendBlocked() bool {
	return p.Attributes.VarBool("friend_block", false)
}

//TradeBlocked Returns true if trade requests are blocked for this player.
func (p *Player) TradeBlocked() bool {
	return p.Attributes.VarBool("trade_block", false)
}

//DuelBlocked Returns true if duel requests are blocked for this player.
func (p *Player) DuelBlocked() bool {
	return p.Attributes.VarBool("duel_block", false)
}

//SetPrivacySettings Sets privacy settings to specified values.
func (p *Player) SetPrivacySettings(chatBlocked, friendBlocked, tradeBlocked, duelBlocked bool) {
	p.Attributes.SetVar("chat_block", chatBlocked)
	p.Attributes.SetVar("friend_block", friendBlocked)
	p.Attributes.SetVar("trade_block", tradeBlocked)
	p.Attributes.SetVar("duel_block", duelBlocked)
}

//SetClientSetting Sets the specified client setting to flag.
func (p *Player) SetClientSetting(id int, flag bool) {
	// TODO: Meaningful names mapped to IDs
	p.Attributes.SetVar("client_setting_"+strconv.Itoa(id), flag)
}

//GetClientSetting Looks up the client setting with the specified ID, and returns it.  If it can't be found, returns false.
func (p *Player) GetClientSetting(id int) bool {
	// TODO: Meaningful names mapped to IDs
	return p.Attributes.VarBool("client_setting_"+strconv.Itoa(id), false)
}

//IsFollowing Returns true if the player is following another mob, otherwise false.
func (p *Player) IsFollowing() bool {
	return p.FollowRadius() >= 0
}

//ServerSeed Returns the seed for the ISAAC cipher provided by the server for this player, if set, otherwise returns 0
func (p *Player) ServerSeed() uint64 {
	return p.TransAttrs.VarLong("server_seed", 0)
}

//SetServerSeed Sets the player's stored server seed to seed for later comparison to ensure we decrypted the login block properly and the player received the proper seed.
func (p *Player) SetServerSeed(seed uint64) {
	p.TransAttrs.SetVar("server_seed", seed)
}

//Reconnecting Returns true if the player is reconnecting, false otherwise.
func (p *Player) Reconnecting() bool {
	return p.TransAttrs.VarBool("reconnecting", false)
}

//SetReconnecting Sets the player's reconnection status to flag.
func (p *Player) SetReconnecting(flag bool) {
	p.TransAttrs.SetVar("reconnecting", flag)
}

//StartFollowing Sets the transient attribute for storing the server index of the player we want to follow to index.
func (p *Player) StartFollowing(radius int) {
	p.TransAttrs.SetVar("followrad", radius)
}

//FollowRadius Returns the radius within which we should follow whatever mob we are following, or -1 if we aren't following anyone.
func (p *Player) FollowRadius() int {
	return p.TransAttrs.VarInt("followrad", -1)
}

//ResetFollowing Resets the transient attribute for storing the server index of the player we want to follow.
func (p *Player) ResetFollowing() {
	p.TransAttrs.UnsetVar("followrad")
	p.ResetPath()
}

//NextTo Returns true if we can walk a straight line to target without colliding with any walls or objects, otherwise returns false.
func (p *Player) NextTo(target Location) bool {
	curLoc := NewLocation(p.X(), p.Y())
	for !curLoc.Equals(target) {
		nextTile := curLoc.NextTileToward(target)
		dir := curLoc.directionTo(nextTile.X(), nextTile.Y())
		switch dir {
		case North:
			if IsTileBlocking(nextTile.X(), nextTile.Y(), WallSouth, true) {
				return false
			}
		case South:
			if IsTileBlocking(nextTile.X(), nextTile.Y(), WallNorth, true) {
				return false
			}
		case East:
			if IsTileBlocking(nextTile.X(), nextTile.Y(), WallWest, true) {
				return false
			}
		case West:
			if IsTileBlocking(nextTile.X(), nextTile.Y(), WallEast, true) {
				return false
			}
		case NorthWest:
			if IsTileBlocking(nextTile.X()+1, nextTile.Y(), WallSouth, true) {
				return false
			}
			if IsTileBlocking(nextTile.X(), nextTile.Y()+1, WallEast, true) {
				return false
			}
			if IsTileBlocking(nextTile.X(), nextTile.Y(), WallSouth|WallEast, true) {
				return false
			}
		case NorthEast:
			if IsTileBlocking(nextTile.X()-1, nextTile.Y(), WallSouth, true) {
				return false
			}
			if IsTileBlocking(nextTile.X(), nextTile.Y()+1, WallWest, true) {
				return false
			}
			if IsTileBlocking(nextTile.X(), nextTile.Y(), WallSouth|WallWest, true) {
				return false
			}
		case SouthWest:
			if IsTileBlocking(nextTile.X()+1, nextTile.Y(), WallNorth, true) {
				return false
			}
			if IsTileBlocking(nextTile.X(), nextTile.Y()-1, WallEast, true) {
				return false
			}
			if IsTileBlocking(nextTile.X(), nextTile.Y(), WallNorth|WallEast, true) {
				return false
			}
		case SouthEast:
			if IsTileBlocking(nextTile.X()-1, nextTile.Y(), WallNorth, true) {
				return false
			}
			if IsTileBlocking(nextTile.X(), nextTile.Y()-1, WallWest, true) {
				return false
			}
			if IsTileBlocking(nextTile.X(), nextTile.Y(), WallNorth|WallWest, true) {
				return false
			}
		}
		curLoc = nextTile
	}

	return true
}

//ResetFollowing Resets the transient attributes holding: Path, Follow radius, and Distanced action triggers...
func (p *Player) ResetAll() {
	p.TransAttrs.UnsetVar("followrad")
	p.ResetFighting()
	p.ResetDistancedAction()
	p.ResetPath()
}

//ArmourPoints Returns the players armour points.
func (p *Player) ArmourPoints() int {
	return p.TransAttrs.VarInt("armour_points", 1)
}

//SetArmourPoints Sets the players armour points to i.
func (p *Player) SetArmourPoints(i int) {
	p.TransAttrs.SetVar("armour_points", i)
}

//PowerPoints Returns the players power points.
func (p *Player) PowerPoints() int {
	return p.TransAttrs.VarInt("power_points", 1)
}

//SetPowerPoints Sets the players power points to i
func (p *Player) SetPowerPoints(i int) {
	p.TransAttrs.SetVar("power_points", i)
}

//AimPoints Returns the players aim points
func (p *Player) AimPoints() int {
	return p.TransAttrs.VarInt("aim_points", 1)
}

//SetAimPoints Sets the players aim points to i.
func (p *Player) SetAimPoints(i int) {
	p.TransAttrs.SetVar("aim_points", i)
}

//MagicPoints Returns the players magic points
func (p *Player) MagicPoints() int {
	return p.TransAttrs.VarInt("magic_points", 1)
}

//SetMagicPoints Sets the players magic points to i
func (p *Player) SetMagicPoints(i int) {
	p.TransAttrs.SetVar("magic_points", i)
}

//PrayerPoints Returns the players prayer points
func (p *Player) PrayerPoints() int {
	return p.TransAttrs.VarInt("prayer_points", 1)
}

//SetPrayerPoints Sets the players prayer points to i
func (p *Player) SetPrayerPoints(i int) {
	p.TransAttrs.SetVar("prayer_points", i)
}

//RangedPoints Returns the players ranged points.
func (p *Player) RangedPoints() int {
	return p.TransAttrs.VarInt("ranged_points", 1)
}

//SetRangedPoints Sets the players ranged points tp i.
func (p *Player) SetRangedPoints(i int) {
	p.TransAttrs.SetVar("ranged_points", i)
}

//Fatigue Returns the players current fatigue.
func (p *Player) Fatigue() int {
	return p.Attributes.VarInt("fatigue", 0)
}

//SetFatigue Sets the players current fatigue to i.
func (p *Player) SetFatigue(i int) {
	p.Attributes.SetVar("fatigue", i)
}

//FightMode Returns the players current fight mode.
func (p *Player) FightMode() int {
	return p.Attributes.VarInt("fight_mode", 0)
}

//SetFightMode Sets the players fightmode to i.  0=all,1=attack,2=defense,3=strength
func (p *Player) SetFightMode(i int) {
	p.Attributes.SetVar("fight_mode", i)
}

//NearbyPlayers Returns nearby players.
func (p *Player) NearbyPlayers() (players []*Player) {
	for _, r := range SurroundingRegions(p.X(), p.Y()) {
		players = append(players, r.Players.NearbyPlayers(p)...)
	}

	return
}

//NearbyObjects Returns nearby objects.
func (p *Player) NearbyObjects() (objects []*Object) {
	for _, r := range SurroundingRegions(p.X(), p.Y()) {
		objects = append(objects, r.Objects.NearbyObjects(p)...)
	}

	return
}

//NewObjects Returns nearby objects that this player is unaware of.
func (p *Player) NewObjects() (objects []*Object) {
	for _, r := range SurroundingRegions(p.X(), p.Y()) {
		for _, o := range r.Objects.NearbyObjects(p) {
			if !p.LocalObjects.Contains(o) {
				objects = append(objects, o)
			}
		}
	}

	return
}

//NewItems Returns nearby ground items that this player is unaware of.
func (p *Player) NewItems() (items []*GroundItem) {
	for _, r := range SurroundingRegions(p.X(), p.Y()) {
		for _, i := range r.Items.NearbyItems(p) {
			if !p.LocalItems.Contains(i) {
				items = append(items, i)
			}
		}
	}

	return
}

//NewPlayers Returns nearby players that this player is unaware of.
func (p *Player) NewPlayers() (players []*Player) {
	for _, r := range SurroundingRegions(p.X(), p.Y()) {
		for _, p1 := range r.Players.NearbyPlayers(p) {
			if !p.LocalPlayers.Contains(p1) {
				players = append(players, p1)
			}
		}
	}

	return
}

//NewNPCs Returns nearby NPCs that this player is unaware of.
func (p *Player) NewNPCs() (npcs []*NPC) {
	for _, r := range SurroundingRegions(p.X(), p.Y()) {
		for _, n := range r.NPCs.NearbyNPCs(p) {
			if !p.LocalNPCs.Contains(n) {
				npcs = append(npcs, n)
			}
		}
	}

	return
}

//SetTradeTarget Sets the variable for the index of the player we are trying to trade
func (p *Player) SetTradeTarget(index int) {
	p.TransAttrs.SetVar("tradetarget", index)
}

//ResetTrade Resets trade-related variables.
func (p *Player) ResetTrade() {
	p.TransAttrs.UnsetVar("tradetarget")
	p.TransAttrs.UnsetVar("trade1accept")
	p.TransAttrs.UnsetVar("trade2accept")
	p.TradeOffer.Clear()
}

//TradeTarget Returns the server index of the player we are trying to trade with, or -1 if we have not made a trade request.
func (p *Player) TradeTarget() int {
	return p.TransAttrs.VarInt("tradetarget", -1)
}

//IsFighting Returns true if this player is currently in a fighting stance, otherwise returns false.
func (p *Player) IsFighting() bool {
	sprite := p.Direction() // Prevent locking too frequently
	return sprite == LeftFighting || sprite == RightFighting
}

func (p *Player) SendPacket(packet *packet.Packet) {
	p.OutgoingPackets <- packet
}

//NewPlayer Returns a reference to a new player.
func NewPlayer(index int, ip string) *Player {
	p := &Player{Mob: &Mob{Entity: &Entity{Index: index, Location: Location{atomic.NewUint32(0), atomic.NewUint32(0)}}, Skillset: &SkillTable{}, State: MSIdle, TransAttrs: &AttributeList{Set: make(map[string]interface{})}}, Attributes: &AttributeList{Set: make(map[string]interface{})}, LocalPlayers: &List{}, LocalNPCs: &List{}, LocalObjects: &List{}, Appearance: NewAppearanceTable(1, 2, true, 2, 8, 14, 0), FriendList: make(map[uint64]bool), KnownAppearances: make(map[int]int), Items: &Inventory{Capacity: 30}, TradeOffer: &Inventory{Capacity: 12}, LocalItems: &List{}, IP: ip, OutgoingPackets: make(chan *packet.Packet, 20)}
	p.Equips[0] = p.Appearance.Head
	p.Equips[1] = p.Appearance.Body
	p.Equips[2] = p.Appearance.Legs
	return p
}
