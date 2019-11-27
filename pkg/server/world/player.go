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
	OutgoingPackets  chan *packet.Packet
	OptionMenuC      chan int8
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

//Connected Returns true if the player is connected, false otherwise.
func (p *Player) Connected() bool {
	return p.TransAttrs.VarBool("connected", false)
}

//SetConnected Sets the player's connected status to flag.
func (p *Player) SetConnected(flag bool) {
	p.TransAttrs.SetVar("connected", flag)
}

//FirstLogin Returns true if this player has never logged in before, otherwise false.
func (p *Player) FirstLogin() bool {
	return p.Attributes.VarBool("first_login", true)
}

//SetFirstLogin Sets the player's persistent logged in before status to flag.
func (p *Player) SetFirstLogin(flag bool) {
	p.Attributes.SetVar("first_login", flag)
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
		dir := curLoc.DirectionTo(nextTile.X(), nextTile.Y())
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

//EquipItem Equips an item to this player, and sends inventory and equipment bonuses.
func (p *Player) EquipItem(item *Item) {
	var itemAffectedTypes = map[int][]int{32: {32, 33}, 33: {32, 33}, 64: {64, 322}, 512: {512, 640, 644},
		8: {8, 24, 8216}, 1024: {1024}, 128: {128, 640, 644}, 644: {128, 512, 640, 644},
		640: {128, 512, 640, 644}, 2048: {2048}, 16: {16, 24, 8216}, 256: {256, 322},
		322: {64, 256, 322}, 24: {8, 16, 24, 8216}, 8216: {8, 16, 24, 8216},
	}
	def := GetEquipmentDefinition(item.ID)
	if def == nil {
		return
	}
	p.TransAttrs.SetVar("self", false)
	p.Items.Range(func(otherItem *Item) bool {
		if otherDef := GetEquipmentDefinition(otherItem.ID); otherDef != nil {
			if otherItem == item || !otherItem.Worn {
				return true
			}
			for _, i := range itemAffectedTypes[def.Type] {
				if i == otherDef.Type {
					p.SetAimPoints(p.AimPoints() - otherDef.Aim)
					p.SetPowerPoints(p.PowerPoints() - otherDef.Power)
					p.SetArmourPoints(p.ArmourPoints() - otherDef.Armour)
					p.SetMagicPoints(p.MagicPoints() - otherDef.Magic)
					p.SetPrayerPoints(p.PrayerPoints() - otherDef.Prayer)
					p.SetRangedPoints(p.RangedPoints() - otherDef.Ranged)
					otherItem.Worn = false
					var value int
					switch otherDef.Position {
					case 0:
						value = p.Appearance.Head
					case 1:
						value = p.Appearance.Body
					case 2:
						value = p.Appearance.Legs
					default:
						value = 0
					}
					p.Equips[otherDef.Position] = value
				}
			}
		}
		return true
	})
	item.Worn = true
	p.SetAimPoints(p.AimPoints() + def.Aim)
	p.SetPowerPoints(p.PowerPoints() + def.Power)
	p.SetArmourPoints(p.ArmourPoints() + def.Armour)
	p.SetMagicPoints(p.MagicPoints() + def.Magic)
	p.SetPrayerPoints(p.PrayerPoints() + def.Prayer)
	p.SetRangedPoints(p.RangedPoints() + def.Ranged)
	p.AppearanceLock.Lock()
	p.Equips[def.Position] = def.Sprite
	p.AppearanceTicket++
	p.AppearanceLock.Unlock()
}

//DequipItem Removes an item from this clients player equips, and sends inventory and equipment bonuses.
func (p *Player) DequipItem(item *Item) {
	def := GetEquipmentDefinition(item.ID)
	if def == nil {
		return
	}
	if !item.Worn {
		return
	}
	p.TransAttrs.SetVar("self", false)
	item.Worn = false
	p.SetAimPoints(p.AimPoints() - def.Aim)
	p.SetPowerPoints(p.PowerPoints() - def.Power)
	p.SetArmourPoints(p.ArmourPoints() - def.Armour)
	p.SetMagicPoints(p.MagicPoints() - def.Magic)
	p.SetPrayerPoints(p.PrayerPoints() - def.Prayer)
	p.SetRangedPoints(p.RangedPoints() - def.Ranged)
	var value int
	switch def.Position {
	case 0:
		value = p.Appearance.Head
	case 1:
		value = p.Appearance.Body
	case 2:
		value = p.Appearance.Legs
	default:
		value = 0
	}
	p.AppearanceLock.Lock()
	p.Equips[def.Position] = value
	p.AppearanceTicket++
	p.AppearanceLock.Unlock()
}

//ResetFollowing Resets the transient attributes holding: Path, Follow radius, and Distanced action triggers...
func (p *Player) ResetAll() {
	p.ResetFighting()
	p.ResetTrade()
	p.ResetDistancedAction()
	p.ResetFollowing()
}

//Fatigue Returns the players current fatigue.
func (p *Player) Fatigue() int {
	return p.Attributes.VarInt("fatigue", 0)
}

//SetFatigue Sets the players current fatigue to i.
func (p *Player) SetFatigue(i int) {
	p.Attributes.SetVar("fatigue", i)
}

//NearbyPlayers Returns nearby players.
func (p *Player) NearbyPlayers() (players []*Player) {
	for _, r := range SurroundingRegions(p.X(), p.Y()) {
		players = append(players, r.Players.NearbyPlayers(p)...)
	}

	return
}

//NearbyPlayers Returns nearby players.
func (p *Player) NearbyNpcs() (npcs []*NPC) {
	for _, r := range SurroundingRegions(p.X(), p.Y()) {
		npcs = append(npcs, r.Players.NearbyNPCs(p)...)
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

func (p *Player) IsTrading() bool {
	return p.HasState(MSTrading)
}

//ResetTrade Resets trade-related variables.
func (p *Player) ResetTrade() {
	p.TransAttrs.UnsetVar("tradetarget")
	p.TransAttrs.UnsetVar("trade1accept")
	p.TransAttrs.UnsetVar("trade2accept")
	p.TradeOffer.Clear()
	p.RemoveState(MSTrading)
}

//TradeTarget Returns the server index of the player we are trying to trade with, or -1 if we have not made a trade request.
func (p *Player) TradeTarget() int {
	return p.TransAttrs.VarInt("tradetarget", -1)
}

func (p *Player) SendPacket(packet *packet.Packet) {
	p.OutgoingPackets <- packet
}

//NewPlayer Returns a reference to a new player.
func NewPlayer(index int, ip string) *Player {
	p := &Player{Mob: &Mob{Entity: &Entity{Index: index, Location: Location{atomic.NewUint32(0), atomic.NewUint32(0)}}, TransAttrs: &AttributeList{Set: make(map[string]interface{})}}, Attributes: &AttributeList{Set: make(map[string]interface{})}, LocalPlayers: &List{}, LocalNPCs: &List{}, LocalObjects: &List{}, Appearance: NewAppearanceTable(1, 2, true, 2, 8, 14, 0), FriendList: make(map[uint64]bool), KnownAppearances: make(map[int]int), Items: &Inventory{Capacity: 30}, TradeOffer: &Inventory{Capacity: 12}, LocalItems: &List{}, IP: ip, OutgoingPackets: make(chan *packet.Packet, 20), OptionMenuC: make(chan int8)}
	p.Transients().SetVar("skills", &SkillTable{})
	p.Equips[0] = p.Appearance.Head
	p.Equips[1] = p.Appearance.Body
	p.Equips[2] = p.Appearance.Legs
	return p
}
