package world

import (
	"fmt"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/packet"
	"go.uber.org/atomic"
	"strconv"
	"sync"
	"time"
)

//AppearanceTable Represents a players appearance.
type AppearanceTable struct {
	Head      int
	Body      int
	Legs      int
	Male      bool
	HeadColor int
	BodyColor int
	LegsColor int
	SkinColor int
}

//NewAppearanceTable returns a reference to a new appearance table with specified parameters
func NewAppearanceTable(head, body int, male bool, hair, top, bottom, skin int) AppearanceTable {
	return AppearanceTable{head, body, 3, male, hair, top, bottom, skin}
}

//player Represents a single player.
type Player struct {
	Username         string
	UserBase37       uint64
	Password         string
	FriendList       map[uint64]bool
	IgnoreList       []uint64
	LocalPlayers     *entityList
	LocalNPCs        *entityList
	LocalObjects     *entityList
	LocalItems       *entityList
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
	Inventory        *Inventory
	Bank             *Inventory
	TradeOffer       *Inventory
	DistancedAction  func() bool
	ActionLock       sync.RWMutex
	OutgoingPackets  chan *packet.Packet
	OptionMenuC      chan int8
	IP               string
	UID              uint8
	Websocket        bool
	Equips           [12]int
	killer           sync.Once
	Kill             chan struct{}
	*Mob
}

//String returns a string populated with the more identifying features of this player.
func (p *Player) String() string {
	return fmt.Sprintf("Player[%d] {'%v'@'%v'}", p.Index, p.Username, p.IP)
}

//SetDistancedAction queues a distanced action to run every game engine tick before path traversal, if action returns true, it will be reset.
func (p *Player) SetDistancedAction(action func() bool) {
	p.ActionLock.Lock()
	p.DistancedAction = action
	p.ActionLock.Unlock()
}

//ResetDistancedAction clears the distanced action, if any is queued.  Should be called any time the player is deliberately performing an action.
func (p *Player) ResetDistancedAction() {
	p.ActionLock.Lock()
	p.DistancedAction = nil
	p.ActionLock.Unlock()
}

//Friends returns true if specified username is in our friend entityList.
func (p *Player) Friends(other uint64) bool {
	for hash := range p.FriendList {
		if hash == other {
			return true
		}
	}
	return false
}

//Ignoring returns true if specified username is in our ignore entityList.
func (p *Player) Ignoring(hash uint64) bool {
	for _, v := range p.IgnoreList {
		if v == hash {
			return true
		}
	}
	return false
}

//ChatBlocked returns true if public chat is blocked for this player.
func (p *Player) ChatBlocked() bool {
	return p.Attributes.VarBool("chat_block", false)
}

//FriendBlocked returns true if private chat is blocked for this player.
func (p *Player) FriendBlocked() bool {
	return p.Attributes.VarBool("friend_block", false)
}

//TradeBlocked returns true if trade requests are blocked for this player.
func (p *Player) TradeBlocked() bool {
	return p.Attributes.VarBool("trade_block", false)
}

//DuelBlocked returns true if duel requests are blocked for this player.
func (p *Player) DuelBlocked() bool {
	return p.Attributes.VarBool("duel_block", false)
}

//SetPrivacySettings sets privacy settings to specified values.
func (p *Player) SetPrivacySettings(chatBlocked, friendBlocked, tradeBlocked, duelBlocked bool) {
	p.Attributes.SetVar("chat_block", chatBlocked)
	p.Attributes.SetVar("friend_block", friendBlocked)
	p.Attributes.SetVar("trade_block", tradeBlocked)
	p.Attributes.SetVar("duel_block", duelBlocked)
}

//SetClientSetting sets the specified client setting to flag.
func (p *Player) SetClientSetting(id int, flag bool) {
	// TODO: Meaningful names mapped to IDs
	p.Attributes.SetVar("client_setting_"+strconv.Itoa(id), flag)
}

//GetClientSetting looks up the client setting with the specified ID, and returns it.  If it can't be found, returns false.
func (p *Player) GetClientSetting(id int) bool {
	// TODO: Meaningful names mapped to IDs
	return p.Attributes.VarBool("client_setting_"+strconv.Itoa(id), false)
}

//IsFollowing returns true if the player is following another mob, otherwise false.
func (p *Player) IsFollowing() bool {
	return p.FollowRadius() >= 0
}

//ServerSeed returns the seed for the ISAAC cipher provided by the server for this player, if set, otherwise returns 0
func (p *Player) ServerSeed() uint64 {
	return p.TransAttrs.VarLong("server_seed", 0)
}

//SetServerSeed sets the player's stored server seed to seed for later comparison to ensure we decrypted the login block properly and the player received the proper seed.
func (p *Player) SetServerSeed(seed uint64) {
	p.TransAttrs.SetVar("server_seed", seed)
}

//Reconnecting returns true if the player is reconnecting, false otherwise.
func (p *Player) Reconnecting() bool {
	return p.TransAttrs.VarBool("reconnecting", false)
}

//SetReconnecting sets the player's reconnection status to flag.
func (p *Player) SetReconnecting(flag bool) {
	p.TransAttrs.SetVar("reconnecting", flag)
}

//Connected returns true if the player is connected, false otherwise.
func (p *Player) Connected() bool {
	return p.TransAttrs.VarBool("connected", false)
}

//SetConnected sets the player's connected status to flag.
func (p *Player) SetConnected(flag bool) {
	p.TransAttrs.SetVar("connected", flag)
}

//FirstLogin returns true if this player has never logged in before, otherwise false.
func (p *Player) FirstLogin() bool {
	return p.Attributes.VarBool("first_login", true)
}

//SetFirstLogin sets the player's persistent logged in before status to flag.
func (p *Player) SetFirstLogin(flag bool) {
	p.Attributes.SetVar("first_login", flag)
}

//StartFollowing sets the transient attribute for storing the radius with which we want to stay near our target
func (p *Player) StartFollowing(radius int) {
	p.TransAttrs.SetVar("followrad", radius)
}

//FollowRadius returns the radius within which we should follow whatever mob we are following, or -1 if we aren't following anyone.
func (p *Player) FollowRadius() int {
	return p.TransAttrs.VarInt("followrad", -1)
}

//ResetFollowing resets the transient attribute for storing the radius within which we want to stay to our target mob
// and resets our path.
func (p *Player) ResetFollowing() {
	p.TransAttrs.UnsetVar("followrad")
	p.ResetPath()
}

//NextTo returns true if we can walk a straight line to target without colliding with any walls or objects,
// otherwise returns false.
func (p *Player) NextTo(target Location) bool {
	curLoc := NewLocation(p.X(), p.Y())
	for !curLoc.Equals(target) {
		nextTile := curLoc.NextTileToward(target)
		dir := curLoc.DirectionTo(nextTile.X(), nextTile.Y())
		switch dir {
		case North:
			if IsTileBlocking(nextTile.X(), nextTile.Y(), ClipSouth, true) {
				return false
			}
		case South:
			if IsTileBlocking(nextTile.X(), nextTile.Y(), ClipNorth, true) {
				return false
			}
		case East:
			if IsTileBlocking(nextTile.X(), nextTile.Y(), ClipWest, true) {
				return false
			}
		case West:
			if IsTileBlocking(nextTile.X(), nextTile.Y(), ClipEast, true) {
				return false
			}
		case NorthWest:
			if IsTileBlocking(nextTile.X()+1, nextTile.Y(), ClipSouth, true) {
				return false
			}
			if IsTileBlocking(nextTile.X(), nextTile.Y()+1, ClipEast, true) {
				return false
			}
			if IsTileBlocking(nextTile.X(), nextTile.Y(), ClipSouth|ClipEast, true) {
				return false
			}
		case NorthEast:
			if IsTileBlocking(nextTile.X()-1, nextTile.Y(), ClipSouth, true) {
				return false
			}
			if IsTileBlocking(nextTile.X(), nextTile.Y()+1, ClipWest, true) {
				return false
			}
			if IsTileBlocking(nextTile.X(), nextTile.Y(), ClipSouth|ClipWest, true) {
				return false
			}
		case SouthWest:
			if IsTileBlocking(nextTile.X()+1, nextTile.Y(), ClipNorth, true) {
				return false
			}
			if IsTileBlocking(nextTile.X(), nextTile.Y()-1, ClipEast, true) {
				return false
			}
			if IsTileBlocking(nextTile.X(), nextTile.Y(), ClipNorth|ClipEast, true) {
				return false
			}
		case SouthEast:
			if IsTileBlocking(nextTile.X()-1, nextTile.Y(), ClipNorth, true) {
				return false
			}
			if IsTileBlocking(nextTile.X(), nextTile.Y()-1, ClipWest, true) {
				return false
			}
			if IsTileBlocking(nextTile.X(), nextTile.Y(), ClipNorth|ClipWest, true) {
				return false
			}
		}
		curLoc = nextTile
	}

	return true
}

//TraversePath if the mob has a path, calling this method will change the mobs location to the next location described by said Path data structure.  This should be called no more than once per game tick.
func (p *Player) TraversePath() {
	path := p.Path()
	if path == nil {
		return
	}
	if p.AtLocation(path.nextTile()) {
		path.CurrentWaypoint++
	}
	if p.FinishedPath() {
		p.ResetPath()
		return
	}
	dst := path.nextTile()
	x, y := p.X(), p.Y()
	next := NewLocation(x, y)
	xBlocked, yBlocked := false, false
	newXBlocked, newYBlocked := false, false
	if y > dst.Y() {
		yBlocked = IsTileBlocking(x, y, ClipNorth, true)
		newYBlocked = IsTileBlocking(x, y-1, ClipSouth, false)
		if !newYBlocked {
			next.y.Dec()
		}
	} else if y < dst.Y() {
		yBlocked = IsTileBlocking(x, y, ClipSouth, true)
		newYBlocked = IsTileBlocking(x, y+1, ClipNorth, false)
		if !newYBlocked {
			next.y.Inc()
		}
	}
	if x > dst.X() {
		xBlocked = IsTileBlocking(x, next.Y(), ClipEast, true)
		newXBlocked = IsTileBlocking(x-1, next.Y(), ClipWest, false)
		if !newXBlocked {
			next.x.Dec()
		}
	} else if x < dst.X() {
		xBlocked = IsTileBlocking(x, next.Y(), ClipWest, true)
		newXBlocked = IsTileBlocking(x+1, next.Y(), ClipEast, false)
		if !newXBlocked {
			next.x.Inc()
		}
	}

	if (xBlocked && yBlocked) || (xBlocked && y == dst.Y()) || (yBlocked && x == dst.X()) {
		p.ResetPath()
		return
	}
	if (newXBlocked && newYBlocked) || (newXBlocked && x != next.X() && y == next.Y()) || (newYBlocked && y != next.Y() && x == next.X()) {
		p.ResetPath()
		return
	}

	if next.X() > x {
		newXBlocked = IsTileBlocking(next.X(), next.Y(), ClipEast, false)
	} else if next.X() < x {
		newXBlocked = IsTileBlocking(next.X(), next.Y(), ClipWest, false)
	}
	if next.Y() > y {
		newYBlocked = IsTileBlocking(next.X(), next.Y(), ClipNorth, false)
	} else if next.Y() < y {
		newYBlocked = IsTileBlocking(next.X(), next.Y(), ClipSouth, false)
	}

	if (newXBlocked && newYBlocked) || (newXBlocked && y == next.Y()) || (newYBlocked && x == next.X()) {
		p.ResetPath()
		return
	}

	p.SetLocation(next, false)
}

//UpdateRegion if this player is currently in a region, removes it from that region, and adds it to the region at x,y
func (p *Player) UpdateRegion(x, y int) {
	curArea := getRegion(p.X(), p.Y())
	newArea := getRegion(x, y)
	if newArea != curArea {
		if curArea.Players.Contains(p) {
			curArea.Players.Remove(p)
		}
		newArea.Players.Add(p)
	}
}

//EquipItem equips an item to this player, and sends inventory and equipment bonuses.
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
	p.NeedsSelf()
	p.Inventory.Range(func(otherItem *Item) bool {
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

//DequipItem removes an item from this players equips, and sends inventory and equipment bonuses.
func (p *Player) DequipItem(item *Item) {
	def := GetEquipmentDefinition(item.ID)
	if def == nil {
		return
	}
	if !item.Worn {
		return
	}
	p.NeedsSelf()
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

//ResetAll in order, calls ResetFighting, ResetTrade, ResetDistancedAction, ResetFollowing, and CloseOptionMenu.
func (p *Player) ResetAll() {
	p.ResetFighting()
	p.ResetTrade()
	p.ResetDistancedAction()
	p.ResetFollowing()
	p.CloseOptionMenu()
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
	for _, r := range surroundingRegions(p.X(), p.Y()) {
		players = append(players, r.Players.NearbyPlayers(p)...)
	}

	return
}

//NearbyPlayers Returns nearby players.
func (p *Player) NearbyNpcs() (npcs []*NPC) {
	for _, r := range surroundingRegions(p.X(), p.Y()) {
		npcs = append(npcs, r.Players.NearbyNpcs(p)...)
	}

	return
}

//NearbyObjects Returns nearby objects.
func (p *Player) NearbyObjects() (objects []*Object) {
	for _, r := range surroundingRegions(p.X(), p.Y()) {
		objects = append(objects, r.Objects.NearbyObjects(p)...)
	}

	return
}

//NewObjects Returns nearby objects that this player is unaware of.
func (p *Player) NewObjects() (objects []*Object) {
	for _, r := range surroundingRegions(p.X(), p.Y()) {
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
	for _, r := range surroundingRegions(p.X(), p.Y()) {
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
	for _, r := range surroundingRegions(p.X(), p.Y()) {
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
	for _, r := range surroundingRegions(p.X(), p.Y()) {
		for _, n := range r.NPCs.NearbyNpcs(p) {
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

//IsTrading returns true if this player is in a trade, otherwise returns false.
func (p *Player) IsTrading() bool {
	return p.HasState(MSTrading)
}

//ResetTrade resets trade-related variables.
func (p *Player) ResetTrade() {
	if p.IsTrading() {
		p.TransAttrs.UnsetVar("tradetarget")
		p.TransAttrs.UnsetVar("trade1accept")
		p.TransAttrs.UnsetVar("trade2accept")
		p.TradeOffer.Clear()
		p.RemoveState(MSTrading)
	}
}

//TradeTarget returns the server index of the player we are trying to trade with, or -1 if we have not made a trade request.
func (p *Player) TradeTarget() int {
	return p.TransAttrs.VarInt("tradetarget", -1)
}

//SendPacket sends a packet to the client.
func (p *Player) SendPacket(packet *packet.Packet) {
	p.OutgoingPackets <- packet
}

//Destroy sends a kill signal to the underlying client to tear down all of the I/O routines and save the player.
func (p *Player) Destroy() {
	p.killer.Do(func() {
		p.Attributes.SetVar("lastIP", p.IP)
		close(p.Kill)
	})
}

//Initialize informs the client of all of the various attributes of this player, and starts the stat normalization
// routine.
func (p *Player) Initialize() {
	p.SetConnected(true)
	AddPlayer(p)
	p.SendPacket(FriendList(p))
	p.SendPacket(IgnoreList(p))
	p.SendPlane()
	p.SendEquipBonuses()
	p.SendInventory()
	p.SendPacket(Fatigue(p))
	// TODO: Not canonical RSC, but definitely good QoL update...
	//  p.SendPacket(FightMode(p))
	p.SendPacket(ClientSettings(p))
	p.SendPacket(PrivacySettings(p))
	if !p.Reconnecting() {
		p.SendPacket(WelcomeMessage)
		if !p.FirstLogin() {
			if tString := p.Attributes.VarString("lastLogin", ""); tString != "" {
				if t, err := time.Parse(time.ANSIC, tString); err == nil {
					p.SendPacket(LoginBox(int(time.Since(t).Hours()/24), p.Attributes.VarString("lastIP", "0.0.0.0")))
				} else {
					log.Info.Println(err)
				}
			}
		} else {
			p.SetFirstLogin(false)
			for i := 0; i < 18; i++ {
				exp := 0
				if i == 3 {
					exp = 1154
				}
				p.Skills().SetCur(i, ExperienceToLevel(exp))
				p.Skills().SetMax(i, ExperienceToLevel(exp))
				p.Skills().SetExp(i, exp)
			}
			p.OpenAppearanceChanger()
		}

	}
	p.SendStats()
	p.Attributes.SetVar("lastLogin", time.Now().Format(time.ANSIC))
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if !p.Connected() {
				return
			}

			for idx := 0; idx < 18; idx++ {
				cur := p.Skills().Current(idx)
				max := p.Skills().Maximum(idx)
				delta := max - cur
				if idx == StatPrayer {
					continue
				}

				if delta > 0 {
					p.SetCurStat(idx, cur+1)
				} else if delta < 0 {
					p.SetCurStat(idx, cur-1)
				}
				if idx != 3 && delta == 1 || delta == -1 {
					// TODO: Look this real message up
					p.Message("Your " + SkillName(idx) + " level has returned to normal.")
				}
			}
		}
	}()
}

//NewPlayer Returns a reference to a new player.
func NewPlayer(index int, ip string) *Player {
	p := &Player{Mob: &Mob{Entity: &Entity{Index: index, Location: Location{atomic.NewUint32(0), atomic.NewUint32(0)}},
		TransAttrs: &AttributeList{set: make(map[string]interface{})}}, Attributes: &AttributeList{set: make(map[string]interface{})},
		LocalPlayers: &entityList{}, LocalNPCs: &entityList{}, LocalObjects: &entityList{}, Appearance: NewAppearanceTable(1, 2, true, 2, 8, 14, 0),
		FriendList: make(map[uint64]bool), KnownAppearances: make(map[int]int), Inventory: &Inventory{Capacity: 30},
		TradeOffer: &Inventory{Capacity: 12}, LocalItems: &entityList{}, IP: ip, OutgoingPackets: make(chan *packet.Packet, 20),
		Kill: make(chan struct{}), Bank: &Inventory{Capacity: 48 * 4, stackEverything: true}}
	p.Transients().SetVar("skills", &SkillTable{})
	p.Equips[0] = p.Appearance.Head
	p.Equips[1] = p.Appearance.Body
	p.Equips[2] = p.Appearance.Legs
	return p
}

//Message sends a message to the player.
func (p *Player) Message(msg string) {
	p.SendPacket(ServerMessage(msg))
}

//OpenAppearanceChanger If the player is not fighting or trading, opens the appearance window.
func (p *Player) OpenAppearanceChanger() {
	if p.IsFighting() || p.IsTrading() {
		return
	}
	p.AddState(MSChangingAppearance)
	p.SendPacket(OpenChangeAppearance)
}

//Chat sends a player NPC chat message packet to the player and all other players around it.  If multiple msgs are
// provided, will sleep the goroutine for 1800ms between each message.
func (p *Player) Chat(msgs ...string) {
	for _, msg := range msgs {
		for _, player := range p.NearbyPlayers() {
			player.SendPacket(PlayerMessage(p, msg))
		}
		p.SendPacket(PlayerMessage(p, msg))

		//		if i < len(msgs)-1 {
		time.Sleep(time.Millisecond * 1800)
		// TODO: is 3 ticks right?
		//		}
	}
}

//OpenOptionMenu opens an option menu with the provided options, and returns the reply index, or -1 upon timeout..
func (p *Player) OpenOptionMenu(options ...string) int {
	// Can get option menu during most states, even fighting, but not trading, or if we're already in a menu...
	if p.IsTrading() || p.HasState(MSOptionMenu) {
		return -1
	}
	p.OptionMenuC = make(chan int8)
	p.AddState(MSOptionMenu)
	defer func() {
		if p.HasState(MSOptionMenu) {
			p.RemoveState(MSOptionMenu)
			close(p.OptionMenuC)
		}
	}()
	p.SendPacket(OptionMenuOpen(options...))

	select {
	case reply := <-p.OptionMenuC:
		if reply < 0 || int(reply) > len(options)-1 || !p.HasState(MSOptionMenu) {
			return -1
		}

		if p.HasState(MSChatting) {
			p.Chat(options[reply])
		}
		return int(reply)
	case <-time.After(time.Second * 10):
		p.SendPacket(OptionMenuClose)
		return -1
	}
}

//CloseOptionMenu closes any open option menus.
func (p *Player) CloseOptionMenu() {
	if p.HasState(MSOptionMenu) {
		p.RemoveState(MSOptionMenu)
		close(p.OptionMenuC)
		p.SendPacket(OptionMenuClose)
	}
}

//CanWalk returns true if this player is in a state that allows walking.
func (p *Player) CanWalk() bool {
	if p.HasState(MSOptionMenu) && p.HasState(MSChatting) {
		// If player tries to walk but is in an option menu, they clearly have closed the menu, so we will kill the
		// routine waiting for a reply when ResetAll is called before the new path is set.
		return true
	}
	return !p.HasState(MSBatching, MSFighting, MSTrading, MSDueling, MSChangingAppearance, MSSleeping, MSChatting, MSBusy)
}

//PlaySound sends a command to the client to play a sound by its file name.
func (p *Player) PlaySound(soundName string) {
	p.SendPacket(Sound(soundName))
}

//SendStat sends the information for the stat at idx to the player.
func (p *Player) SendStat(idx int) {
	p.SendPacket(PlayerStat(p, idx))
}

//SendStats sends all stat information to this player.
func (p *Player) SendStats() {
	p.SendPacket(PlayerStats(p))
}

//SendInventory sends inventory information to this player.
func (p *Player) SendInventory() {
	p.SendPacket(InventoryItems(p))
}

//SetCurStat sets this players current stat at idx to lvl and updates the client about it.
func (p *Player) SetCurStat(idx int, lvl int) {
	p.Skills().SetCur(idx, lvl)
	p.SendStat(idx)
}

//SetMaxStat sets this players maximum stat at idx to lvl and updates the client about it.
func (p *Player) SetMaxStat(idx int, lvl int) {
	p.Skills().SetMax(idx, lvl)
	p.Skills().SetExp(idx, LevelToExperience(lvl))
	p.SendStat(idx)
}

//AddItem Adds amount of the item with specified id to the players inventory, if possible, and updates the client about it.
func (p *Player) AddItem(id, amount int) {
	if !ItemDefs[id].Stackable {
		for i := 0; i < amount; i++ {
			if p.Inventory.Size() >= p.Inventory.Capacity {
				break
			}
			p.Inventory.Add(id, 1)
		}
	} else {
		p.Inventory.Add(id, amount)
	}
	p.SendInventory()
}

//Killed kills this player, dropping all of its items where it stands.
func (p *Player) Killed() {
	p.Transients().SetVar("deathTime", time.Now())
	p.PlaySound("death")
	p.SendPacket(Death)
	for i := 0; i < 18; i++ {
		p.Skills().SetCur(i, p.Skills().Maximum(i))
	}
	p.SendStats()
	// TODO: Keep 3 most valuable items

	if killer, ok := p.FightTarget().(*Player); killer != nil && ok {
		AddItem(NewGroundItemFor(killer.UserBase37, 20, 1, p.X(), p.Y()))
		p.Inventory.Range(func(item *Item) bool {
			p.DequipItem(item)
			AddItem(NewGroundItemFor(killer.UserBase37, item.ID, item.Amount, p.X(), p.Y()))
			return true
		})
	} else {
		p.Inventory.Range(func(item *Item) bool {
			p.DequipItem(item)
			AddItem(NewGroundItem(item.ID, item.Amount, p.X(), p.Y()))
			return true
		})
	}
	p.Inventory.Clear()
	p.SendInventory()
	p.SendEquipBonuses()
	p.ResetFighting()
	plane := p.Plane()
	p.SetLocation(SpawnPoint, true)
	if p.Plane() != plane {
		p.SendPlane()
	}
}

//SendPlane sends the current plane of this player.
func (p *Player) SendPlane() {
	p.SendPacket(PlaneInfo(p))
}

//SendEquipBonuses sends the current equipment bonuses of this player.
func (p *Player) SendEquipBonuses() {
	p.SendPacket(EquipmentStats(p))
}

//RemoveItem removes amount of the item at index in this players inventory, then updates the client about it.
func (p *Player) RemoveItem(index, amount int) {
	p.Inventory.Remove(index, amount)
	p.SendInventory()
}

//RemoveItemByID Removes amount of the item with specified id in this players inventory, then updates the client about
// it.
func (p *Player) RemoveItemByID(id, amount int) {
	p.Inventory.RemoveByID(id, amount)
	p.SendInventory()
}

//Damage sends a player damage bubble for this player to itself and any nearby players.
func (p *Player) Damage(amt int) {
	for _, player := range p.NearbyPlayers() {
		player.SendPacket(PlayerDamage(p, amt))
	}
	p.SendPacket(PlayerDamage(p, amt))
}

//SetStat sets the current, maximum, and experience levels of the skill at idx to lvl, and updates the client about it.
func (p *Player) SetStat(idx, lvl int) {
	p.Skills().SetCur(idx, lvl)
	p.Skills().SetMax(idx, lvl)
	p.Skills().SetExp(idx, LevelToExperience(lvl))
	p.SendStat(idx)
}

//OpenBank opens a bank screen for the player and sets the appropriate state variables.
func (p *Player) OpenBank() {
	if p.IsFighting() || p.IsTrading() || p.HasState(MSBanking) {
		return
	}
	p.AddState(MSBanking)
	p.SendPacket(BankOpen(p))
}

//CloseBank closes the bank screen for this player and sets the appropriate state variables
func (p *Player) CloseBank() {
	if !p.HasState(MSBanking) {
		return
	}
	p.RemoveState(MSBanking)
	p.SendPacket(BankClose)
}

//SendUpdateTimer sends a system update countdown timer to the client.
func (p *Player) SendUpdateTimer() {
	p.SendPacket(SystemUpdate(int(time.Until(UpdateTime).Seconds())))
}
