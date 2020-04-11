/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package world

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/spkaeros/rscgo/pkg/game/entity"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/strutil"
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

func DefaultAppearance() AppearanceTable {
	return NewAppearanceTable(1, 2, true, 2, 8, 14, 0)
}

type friendSet map[uint64]bool

type FriendsList struct {
	sync.RWMutex
	friendSet
	Owner string
}

func (f *FriendsList) contains(name string) bool {
	f.RLock()
	defer f.RUnlock()
	_, ok := f.friendSet[strutil.Base37.Encode(name)]
	return ok
}

func (f *FriendsList) containsHash(hash uint64) bool {
	f.RLock()
	defer f.RUnlock()
	_, ok := f.friendSet[hash]
	return ok
}

func (f *FriendsList) Add(name string) {
	f.Lock()
	defer f.Unlock()
	hash := strutil.Base37.Encode(name)
	p, ok := Players.FromUserHash(hash)
	f.friendSet[hash] = p != nil && ok && (p.FriendList.contains(f.Owner) || !p.FriendBlocked())
}

func (f *FriendsList) Remove(name string) {
	f.Lock()
	defer f.Unlock()
	hash := strutil.Base37.Encode(name)
	delete(f.friendSet, hash)
	if p, ok := Players.FromUserHash(hash); ok && p.FriendList.contains(f.Owner) {
		p.SendPacket(FriendUpdate(strutil.Base37.Encode(f.Owner), false))
	}
}

func (f *FriendsList) ForEach(fn func(string, bool) bool) {
	f.Lock()
	defer f.Unlock()
	for name, status := range f.friendSet {
		if fn(strutil.Base37.Decode(name), status) {
			break
		}
	}
}

func (f *FriendsList) size() int {
	f.RLock()
	defer f.RUnlock()
	return len(f.friendSet)
}

//player Represents a single player.
type Player struct {
	LocalPlayers     *MobList
	LocalNPCs        *MobList
	LocalObjects     *entityList
	LocalItems       *entityList
	FriendList       *FriendsList
	IgnoreList       []uint64
	Appearance       AppearanceTable
	KnownAppearances map[int]int
	AppearanceReq    []*Player
	AppearanceLock   sync.RWMutex
	Attributes       *entity.AttributeList
	Inventory        *Inventory
	TradeOffer       *Inventory
	DuelOffer        *Inventory
	Duel             struct {
		rules []string
		
	}
	DistancedAction  func() bool
	ActionLock       sync.RWMutex
	OutgoingPackets  chan *net.Packet
	ReplyMenuC       chan int8
	Equips           [12]int
	killer           sync.Once
	KillC            chan struct{}
	UpdateWG         sync.RWMutex
	Tickables        []interface{}
	*Mob
}

func (p *Player) UsernameHash() uint64 {
	return p.VarLong("username", strutil.Base37.Encode("nil"))
}

func (p *Player) Bank() *Inventory {
	i, ok := p.Var("bank")
	if ok {
		return i.(*Inventory)
	}
	return nil
}

func (p *Player) CanAttack(target entity.MobileEntity) bool {
	if target.IsNpc() {
		return target.(*NPC).Attackable()
	}
	if p.State()&StateFightingDuel==StateFightingDuel {
		return p.DuelTarget() == target && p.DuelMagic()
	}
	p1 := target.(*Player)
	ourWild := p.Wilderness()
	targetWild := p1.Wilderness()
	if ourWild < 1 || targetWild < 1 {
		p.Message("You cannot attack other players outside of the wilderness!")
		return false
	}
	delta := p.CombatDelta(target)
	if delta < 0 {
		delta = -delta
	}
	if delta > ourWild {
		p.Message("You must move to at least level " + strconv.Itoa(delta) + " wilderness to attack " + p1.Username() + "!")
		return false
	}
	if delta > targetWild {
		p.Message(p1.Username() + " is not in high enough wilderness for you to attack!")
		return false
	}
	return true
}

func (p *Player) Username() string {
	return strutil.Base37.Decode(p.VarLong("username", strutil.Base37.Encode("NIL")))
}

func (p *Player) CurrentIP() string {
	return p.VarString("currentIP", "0.0.0.0")
}

func (p *Player) Rank() int {
	return p.VarInt("rank", 0)
}

func (p *Player) DatabaseID() int {
	return p.VarInt("dbID", -1)
}

func (p *Player) AppearanceTicket() int {
	return p.VarInt("appearanceTicket", 0)
}

//String returns a string populated with the more identifying features of this player.
func (p *Player) String() string {
	return fmt.Sprintf("'%s'[%d]@'%s'", p.Username(), p.Index, p.CurrentIP())
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

//FriendsWith returns true if specified username is in our friend entityList.
func (p *Player) FriendsWith(other uint64) bool {
	return p.FriendList.containsHash(other)
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

func (p *Player) UpdateStatus(status bool) {
	Players.Range(func(player *Player) {
		if player.FriendList.contains(p.Username()) {
			if p.FriendList.contains(player.Username()) || !p.FriendBlocked() {
				p.SendPacket(FriendUpdate(p.UsernameHash(), status))
			}
		}
	})
}

//WalkingRangedAction Runs `fn` once arriving anywhere within 5 tiles in any direction of `t`,
// with a straight line of sight, e.g no intersecting boundaries, large objects, walls, etc.
// Runs everything on game engine ticks, retries until catastrophic failure or success.
func (p *Player) WalkingRangedAction(t entity.MobileEntity, fn func()) {
	//if t == nil || p.State()&StateFightingDuel == 0 || (p.State()&StateFightingDuel == StateFightingDuel &&
	//		(!p.VarBool("duelCanMagic", true) || t != p.DuelTarget())) {
	//	p.ResetPath()
	//	return
	//}
	p.WalkingArrivalAction(t, 5, fn)
}

//WalkingArrivalAction Runs `action` once arriving within dist (min 1 max 2 tiles)
// of `target` mob, with a straight line of sight, e.g no intersecting boundaries, large
// objects, walls, etc.
// Runs everything on game engine ticks, retries until catastrophic failure or success.
func (p *Player) WalkingArrivalAction(target entity.MobileEntity, dist int, action func()) {
	p.SetDistancedAction(func() bool {
		if target == nil {
			p.ResetPath()
			return true
		}
		if p.WithinRange(NewLocation(target.X(), target.Y()), dist) {
			// make sure we can shoot them without obstacles
			if !p.CanReachMob(target) {
				return false
			}
			// shoot them
			action()
			return true
		}
		p.WalkTo(NewLocation(target.X(), target.Y()))
		return false
	})
}

//CanReachMob Check if we can reach a mob traversing the most direct tiles toward them, e.g straight lines.
// Used to check ranged combat attacks, or trade requests, basically anything needing local interactions.
func (p *Player) CanReachMob(target entity.MobileEntity) bool {
	if p.FinishedPath() && p.VarInt("triedReach", 0) >= 5 {
		// Tried reaching one mob >=5 times without single success, abort early.
		p.ResetPath()
		p.UnsetVar("triedReach")
		return false
	}
	p.Inc("triedReach", 1)
		
	
	pathX := p.X()
	pathY := p.Y()
	for steps := 0; steps < p.VarInt("viewRadius", 16)+5; steps++ {
		// check deltas
		if pathX == target.X() && pathY == target.Y() {
			p.UnsetVar("triedReach")
			return true
		}
		if !p.Reachable(pathX, pathY) {
			return false
		}
		
		// Update coords toward target in a straight line
		if pathX<target.X() {
			pathX++
		} else if pathX>target.X() {
			pathX--
		}
		
		if pathY<target.Y() {
			pathY++
		} else if pathY>target.Y() {
			pathY--
		}
	}
	return false
}

//SetPrivacySettings sets privacy settings to specified values.
func (p *Player) SetPrivacySettings(chatBlocked, friendBlocked, tradeBlocked, duelBlocked bool) {
	p.Attributes.SetVar("chat_block", chatBlocked)
	p.Attributes.SetVar("friend_block", friendBlocked)
	p.Attributes.SetVar("trade_block", tradeBlocked)
	p.Attributes.SetVar("duel_block", duelBlocked)

	//Players.Range(func(player *Player) {
	//	if player.FriendList.contains(p.Username()) {
	//		if !p.FriendList.contains(player.Username()) || p.FriendBlocked() {
	//			p.SendPacket(FriendUpdate(p.UsernameHash(), !p.FriendBlocked() || p.FriendList.contains(player.Username())))
	//		}
	//	}
	//})
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

//ServerSeed returns the seed for the ISAAC cipher provided by the game for this player, if set, otherwise returns 0
func (p *Player) ServerSeed() uint64 {
	return p.VarLong("server_seed", 0)
}

//SetServerSeed sets the player's stored game seed to seed for later comparison to ensure we decrypted the login block properly and the player received the proper seed.
func (p *Player) SetServerSeed(seed uint64) {
	p.SetVar("server_seed", seed)
}

//Reconnecting returns true if the player is reconnecting, false otherwise.
func (p *Player) Reconnecting() bool {
	return p.VarBool("reconnecting", false)
}

//SetReconnecting sets the player's reconnection status to flag.
func (p *Player) SetReconnecting(flag bool) {
	p.SetVar("reconnecting", flag)
}

//Connected returns true if the player is connected, false otherwise.
func (p *Player) Connected() bool {
	return p.VarBool("connected", false)
}

//SetConnected sets the player's connected status to flag.
func (p *Player) SetConnected(flag bool) {
	p.SetVar("connected", flag)
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
	p.SetVar("followrad", radius)
}

//FollowRadius returns the radius within which we should follow whatever mob we are following, or -1 if we aren't following anyone.
func (p *Player) FollowRadius() int {
	return p.VarInt("followrad", -1)
}

//ResetFollowing resets the transient attribute for storing the radius within which we want to stay to our target mob
// and resets our path.
func (p *Player) ResetFollowing() {
	p.UnsetVar("followrad")
	p.ResetPath()
}

//NextTo returns true if we can walk a straight line to target without colliding with any walls or objects,
// otherwise returns false.
func (p *Player) NextTo(target Location) bool {
	return p.Reachable(target.X(), target.Y())
}

func (p *Player) NextToCoords(x, y int) bool {
	return p.NextTo(NewLocation(x, y))
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
	dst := p.NextTileToward(path.nextTile())
	
	if p.FinishedPath() {
		p.ResetPath()
		return
	}

	if !p.Reachable(dst.X(), dst.Y()) {
		p.ResetPath()
		return
	}

	p.SetLocation(dst, false)
}

func (l Location) Blocked() bool {
	return false
}

func (l Location) Reachable(x, y int) bool {
	dst := NewLocation(x, y)
	if l.LongestDelta(dst) > 1 {
		dst = l.NextTileToward(dst)
	}
	bitmask := byte(ClipBit(l.DirectionToward(dst)))
	dstmask := byte(ClipBit(dst.DirectionToward(l)))
	// check mask of our tile and dst tile
	if IsTileBlocking(l.X(), l.Y(), bitmask, true) || IsTileBlocking(dst.X(), dst.Y(), dstmask, false) {
		return false
	}

	// does the next step toward our goal affect both X and Y coords?
	if dst.X() != l.X() && dst.Y() != l.Y() {
		// if so, we must scan for adjacent bitmasks of certain diags('_|', '‾|', '|‾' or '|_')
		// Since / and \ masks block a whole tile, those masks are auto-checked in the guts of the API
		// However, this leaves possible holes in x+1,y and x,y+1 at certain angles where we should block
		masks := l.Masks(dst.X(), dst.Y())
		if IsTileBlocking(l.X(), dst.Y(), masks[0], false) && IsTileBlocking(dst.X(), l.Y(), masks[1], false) {
			return false
		}
	}
	return true
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

//DistributeMeleeExp This is a helper method to distribute experience amongst the players melee stats according to
// its current fight stance.
//
// If the player is in controlled stance, each melee skill gets (experience).
// Otherwise, whatever fight stance the player was in will get (experience)*3, and hits will get (experience).
func (p *Player) DistributeMeleeExp(experience int) {
	switch p.FightMode() {
	case 0:
		for i := 0; i < 3; i++ {
			p.IncExp(i, experience)
		}
	case 1:
		p.IncExp(entity.StatStrength, experience*3)
	case 2:
		p.IncExp(entity.StatAttack, experience*3)
	case 3:
		p.IncExp(entity.StatDefense, experience*3)
	}
	p.IncExp(entity.StatHits, experience)
}

//EquipItem equips an item to this player, and sends inventory and equipment bonuses.
func (p *Player) EquipItem(item *Item) {
	reqs := ItemDefs[item.ID].Requirements
	if reqs != nil {
		var needed string
		for skill, lvl := range reqs {
			if p.Skills().Current(skill) < lvl {
				needed += strconv.Itoa(lvl) + " " + entity.SkillName(skill) + ", "
			}
		}
		if len(needed) > 0 {
			p.Message("You must have at least " + needed[:len(needed)-2] + " to wield a " + item.Name())
			return
		}
	}
	def := GetEquipmentDefinition(item.ID)
	if def == nil {
		return
	}
	if def.Female && p.Appearance.Male {
		// TODO: Look up canonical message
		p.Message("You must be a female to wear that")
		return
	}
	p.PlaySound("click")
	p.Inventory.Range(func(otherItem *Item) bool {
		otherDef := GetEquipmentDefinition(otherItem.ID)
		if otherItem == item || !otherItem.Worn || otherDef == nil || def.Type&otherDef.Type == 0 {
			return true
		}
		p.SetAimPoints(p.AimPoints() - otherDef.Aim)
		p.SetPowerPoints(p.PowerPoints() - otherDef.Power)
		p.SetArmourPoints(p.ArmourPoints() - otherDef.Armour)
		p.SetMagicPoints(p.MagicPoints() - otherDef.Magic)
		p.SetPrayerPoints(p.PrayerPoints() - otherDef.Prayer)
		p.SetRangedPoints(p.RangedPoints() - otherDef.Ranged)
		otherItem.Worn = false
		p.AppearanceLock.Lock()
		if otherDef.Type&1 == 1 {
			p.Equips[otherDef.Position] = p.Appearance.Head
		} else if otherDef.Type&2 == 2 {
			p.Equips[otherDef.Position] = p.Appearance.Body
		} else if otherDef.Type&4 == 4 {
			p.Equips[otherDef.Position] = p.Appearance.Legs
		} else {
			p.Equips[otherDef.Position] = 0
		}
		p.AppearanceLock.Unlock()
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
	p.AppearanceLock.Unlock()
	p.UpdateAppearance()
	p.SendEquipBonuses()
	p.SendInventory()
}

func (p *Player) UpdateAppearance() {
	p.SetAppearanceChanged()
	p.SetVar("appearanceTicket", p.AppearanceTicket()+1)
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
	item.Worn = false
	p.SetAimPoints(p.AimPoints() - def.Aim)
	p.SetPowerPoints(p.PowerPoints() - def.Power)
	p.SetArmourPoints(p.ArmourPoints() - def.Armour)
	p.SetMagicPoints(p.MagicPoints() - def.Magic)
	p.SetPrayerPoints(p.PrayerPoints() - def.Prayer)
	p.SetRangedPoints(p.RangedPoints() - def.Ranged)
	p.AppearanceLock.Lock()
	if def.Type&1 == 1 {
		p.Equips[def.Position] = p.Appearance.Head
	} else if def.Type&2 == 2 {
		p.Equips[def.Position] = p.Appearance.Body
	} else if def.Type&4 == 4 {
		p.Equips[def.Position] = p.Appearance.Legs
	} else {
		p.Equips[def.Position] = 0
	}
	p.AppearanceLock.Unlock()
	p.UpdateAppearance()
	p.SendEquipBonuses()
	p.SendInventory()
}

//ResetAll in order, calls ResetFighting, ResetTrade, ResetDistancedAction, ResetFollowing, and CloseOptionMenu.
func (p *Player) ResetAll() {
	p.ResetFighting()
	p.ResetDuel()
	p.ResetTrade()
	p.ResetDistancedAction()
	p.ResetFollowing()
	p.CloseOptionMenu()
	p.CloseBank()
	p.CloseShop()
}

func (p *Player) ResetAllExceptDueling() {
	p.ResetTrade()
	p.ResetDistancedAction()
	p.ResetFollowing()
	p.CloseOptionMenu()
	p.CloseBank()
	p.CloseShop()
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
		r.Players.RangePlayers(func(p1 *Player) bool {
			if p.WithinRange(p1.Location, 16) && p != p1 {
				players = append(players, p1)
			}
			return false
		})
	}

	return
}

//NearbyNpcs Returns nearby NPCs.
func (p *Player) NearbyNpcs() (npcs []*NPC) {
	for _, r := range surroundingRegions(p.X(), p.Y()) {
		r.NPCs.RangeNpcs(func(n *NPC) bool {
			if p.WithinRange(n.Location, 16) && !n.Location.Equals(DeathPoint) {
				npcs = append(npcs, n)
			}
			return false
		})
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
		r.Players.RangePlayers(func(p1 *Player) bool {
			if !p.LocalPlayers.Contains(p1) && p != p1 && p.WithinRange(p1.Location, 15) {
				players = append(players, p1)
			}
			return false
		})
	}

	return
}

//NewNPCs Returns nearby NPCs that this player is unaware of.
func (p *Player) NewNPCs() (npcs []*NPC) {
	for _, r := range surroundingRegions(p.X(), p.Y()) {
		r.NPCs.RangeNpcs(func(n *NPC) bool {
			if !p.LocalNPCs.Contains(n) && p.WithinRange(n.Location, 15) {
				npcs = append(npcs, n)
			}
			return false
		})
	}

	return
}

//SetTradeTarget Sets the variable for the index of the player we are trying to trade
func (p *Player) SetTradeTarget(index int) {
	p.SetVar("tradetarget", index)
}

//IsTrading returns true if this player is in a trade, otherwise returns false.
func (p *Player) IsTrading() bool {
	return p.HasState(StateTrading)
}

func (p *Player) IsPanelOpened() bool {
	return p.HasState(StatePanelActive)
}

//ResetTrade resets trade-related variables.
func (p *Player) ResetTrade() {
	if p.IsTrading() {
		p.UnsetVar("tradetarget")
		p.UnsetVar("trade1accept")
		p.UnsetVar("trade2accept")
		p.TradeOffer.Clear()
		p.RemoveState(StateTrading)
	}
}

//TradeTarget returns the game index of the player we are trying to trade with, or -1 if we have not made a trade request.
func (p *Player) TradeTarget() int {
	return p.VarInt("tradetarget", -1)
}

//CombatDelta returns the difference between our combat level and the other mobs combat level
func (p *Player) CombatDelta(other entity.MobileEntity) int {
	delta := p.Skills().CombatLevel() - other.Skills().CombatLevel()
	if delta < 0 {
		return -delta
	}
	return delta
}

//DuelAccepted returns the status of the specified duel negotiation screens accepted button for this player.
// Valid screens are 1 and 2.
func (p *Player) DuelAccepted(screen int) bool {
	return p.VarBool("duel" + strconv.Itoa(screen) + "accept", false)
}

//DuelAccepted returns the status of the specified duel negotiation screens accepted button for this player.
// Valid screens are 1 and 2.
func (p *Player) SetDuelAccepted(screen int, b bool) {
	duelAttr := "duel" + strconv.Itoa(screen) + "accept"
	if b && screen == 2 && !p.DuelAccepted(1) {
		log.Suspicious.Println("Attempt to set duel2accept before duel1accept:", p.String())
		return
	}
	p.SetVar(duelAttr, b)
}

//SetDuelRule sets the duel rule associated with the specified index to b.
// Valid rule indices are 0 through 3.
func (p *Player) SetDuelRule(index int, b bool) {
	rules := [4]string{"duelCanRetreat", "duelCanMagic", "duelCanPrayer", "duelCanEquip"}
	p.SetVar(rules[index], !b)
}

//DuelRule returns the rule associated with the specified index provided.
// Valid rule indices are 0 through 3.
func (p *Player) duelRule(index int) bool {
	return p.VarBool([4]string{"duelCanRetreat", "duelCanMagic", "duelCanPrayer", "duelCanEquip"}[index], true)
}

func (p *Player) DuelRetreating() bool {
	return p.duelRule(0)
}

func (p *Player) DuelMagic() bool {
	return p.duelRule(1)
}

func (p *Player) DuelPrayer() bool {
	return p.duelRule(2)
}

func (p *Player) DuelEquipment() bool {
	return p.duelRule(3)
}

//ResetDuel resets duel-related variables.
func (p *Player) ResetDuel() {
	//if target := p.DuelTarget(); target != nil && target.IsDueling()
	if p.IsDueling() {
		p.ResetDuelAccepted()
		p.ResetDuelRules()
		p.DuelOffer.Clear()
		p.ResetDuelTarget()
		p.RemoveState(StateDueling)
	}
}

//IsDueling returns true if this player is negotiating a duel, otherwise returns false.
func (p *Player) IsDueling() bool {
	return p.HasState(StateDueling)
}

//SetDuelTarget Sets p1 as the receivers dueling target.
func (p *Player) SetDuelTarget(p1 *Player) {
	p.SetVar("duelTarget", p1)
}

//ResetDuelTarget Removes receivers duel target, if any.
func (p *Player) ResetDuelTarget() {
	p.UnsetVar("duelTarget")
}

//ResetDuelAccepted Resets receivers duel negotiation settings to indicate that neither screens are accepted.
func (p *Player) ResetDuelAccepted() {
	p.SetDuelAccepted(1, false)
	p.SetDuelAccepted(2, false)
}

func (p *Player) ResetDuelRules() {
	for i := 0; i < 4; i++ {
		p.SetDuelRule(i, false)
	}
}

//DuelTarget Returns the player that the receiver is targeting to duel with, or if none, returns nil
func (p *Player) DuelTarget() *Player {
	if p1, ok := p.VarPlayer("duelTarget").(*Player); ok {
		return p1
	}
	return nil
}

//SendPacket sends a net to the client.
func (p *Player) SendPacket(packet *net.Packet) {
	if p == nil || !p.Connected() {
		return
	}
	p.OutgoingPackets <- packet
}

//Destroy sends a kill signal to the underlying client to tear down all of the I/O routines and save the player.
func (p *Player) Destroy() {
	p.killer.Do(func() {
		if p.Connected() {
			p.SendPacket(Logout)
			p.ResetAll()
		}
		p.Inventory.Owner = nil
		close(p.KillC)
	})
}

func (p *Player) AtObject(object *Object) bool {
	x, y := p.X(), p.Y()
	bounds := object.Boundaries()
	if ObjectDefs[object.ID].CollisionType == 2 || ObjectDefs[object.ID].CollisionType == 3 {
		return (p.NextTo(bounds[0]) || p.NextTo(bounds[1])) && (x >= bounds[0].X() && x <= bounds[1].X() && y >= bounds[0].Y() && y <= bounds[1].Y())
	}

	return p.CanReach(bounds) || (p.FinishedPath() && p.CanReachDiag(bounds))
}

func (p *Player) CanReach(bounds [2]Location) bool {
	x, y := p.X(), p.Y()

	if x >= bounds[0].X() && x <= bounds[1].X() && y >= bounds[0].Y() && y <= bounds[1].Y() {
		return true
	}
	if x-1 >= bounds[0].X() && x-1 <= bounds[1].X() && y >= bounds[0].Y() && y <= bounds[1].Y() &&
		(CollisionData(x-1, y).CollisionMask&ClipWest) == 0 {
		return true
	}
	if x+1 >= bounds[0].X() && x+1 <= bounds[1].X() && y >= bounds[0].Y() && y <= bounds[1].Y() &&
		(CollisionData(x+1, y).CollisionMask&ClipEast) == 0 {
		return true
	}
	if x >= bounds[0].X() && x <= bounds[1].X() && bounds[0].Y() <= y-1 && bounds[1].Y() >= y-1 &&
		(CollisionData(x, y-1).CollisionMask&ClipSouth) == 0 {
		return true
	}
	if x >= bounds[0].X() && x <= bounds[1].X() && bounds[0].Y() <= y+1 && bounds[1].Y() >= y+1 &&
		(CollisionData(x, y+1).CollisionMask&ClipNorth) == 0 {
		return true
	}
	return false
}

func (p *Player) CanReachDiag(bounds [2]Location) bool {
	x, y := p.X(), p.Y()
	if x-1 >= bounds[0].X() && x-1 <= bounds[1].X() && y-1 >= bounds[0].Y() && y-1 <= bounds[1].Y() &&
		(CollisionData(x-1, y-1).CollisionMask&ClipSouth|ClipWest) == 0 {
		return true
	}
	if x-1 >= bounds[0].X() && x-1 <= bounds[1].X() && y+1 >= bounds[0].Y() && y+1 <= bounds[1].Y() &&
		(CollisionData(x-1, y+1).CollisionMask&ClipNorth|ClipWest) == 0 {
		return true
	}
	if x+1 >= bounds[0].X() && x+1 <= bounds[1].X() && y-1 >= bounds[0].Y() && y-1 <= bounds[1].Y() &&
		(CollisionData(x+1, y-1).CollisionMask&ClipSouth|ClipEast) == 0 {
		return true
	}
	if x+1 >= bounds[0].X() && x+1 <= bounds[1].X() && y+1 >= bounds[0].Y() && y+1 <= bounds[1].Y() &&
		(CollisionData(x+1, y+1).CollisionMask&ClipNorth|ClipEast) == 0 {
		return true
	}

	return false
}

func (p *Player) SendFatigue() {
	p.SendPacket(Fatigue(p))
}

//Initialize informs the client of all of the various attributes of this player, and starts the stat normalization
// routine.
func (p *Player) Initialize() {
	p.SetAppearanceChanged()
	p.SetSpriteUpdated()
	AddPlayer(p)
	p.UpdateStatus(true)
	p.SendPacket(FriendList(p))
	p.SendPacket(IgnoreList(p))
	// TODO: Not canonical RSC, but definitely good QoL update...
	//  p.SendPacket(FightMode(p))
	p.SendPacket(ClientSettings(p))
	p.SendPacket(PrivacySettings(p))
	timestamp := p.Attributes.VarTime("lastLogin")
	if timestamp.IsZero() {
		for i := 0; i < 18; i++ {
			if i != 3 {
				p.Skills().SetCur(i, 1)
				p.Skills().SetMax(i, 1)
				p.Skills().SetExp(i, 0)
			}
		}
		p.Skills().SetCur(entity.StatHits, 10)
		p.Skills().SetMax(entity.StatHits, 10)
		p.Skills().SetExp(entity.StatHits, entity.LevelToExperience(10))

		p.Bank().Add(546, 96000)
		p.Bank().Add(373, 96000)

		p.Inventory.Add(77, 1)
		p.Inventory.Add(316, 1)

		p.OpenAppearanceChanger()
	}
	if !p.Reconnecting() {
		p.SendPacket(WelcomeMessage)
		p.SendPacket(LoginBox(int(time.Since(timestamp).Hours()/24), p.Attributes.VarString("lastIP", "0.0.0.0")))
	}
	p.SendEquipBonuses()
	p.SendInventory()
	p.SendFatigue()
	p.SendCombatPoints()
	p.SendStats()
	p.SendPlane()
	p.Attributes.SetVar("lastLogin", time.Now())
	for _, fn := range LoginTriggers {
		fn(p)
	}
}

//NewPlayer Returns a reference to a new player.
func NewPlayer(index int, ip string) *Player {
	p := &Player{Mob: &Mob{Entity: &Entity{Index: index, Location: Lumbridge.Clone()}, AttributeList: entity.NewAttributeList()},
		Attributes: entity.NewAttributeList(), LocalPlayers: NewMobList(), LocalNPCs: NewMobList(), LocalObjects: &entityList{},
		Appearance: DefaultAppearance(), FriendList: &FriendsList{friendSet: make(friendSet)}, KnownAppearances: make(map[int]int),
		Inventory: &Inventory{Capacity: 30}, TradeOffer: &Inventory{Capacity: 12}, DuelOffer: &Inventory{Capacity: 8},
		LocalItems: &entityList{}, OutgoingPackets: make(chan *net.Packet, 20), KillC: make(chan struct{})}
	p.SetVar("currentIP", ip)
	p.SetVar("viewRadius", 16)
	p.SetVar("skills", &entity.SkillTable{})
	p.SetVar("bank", &Inventory{Capacity: 48 * 4, stackEverything: true})

	p.Equips[0] = p.Appearance.Head
	p.Equips[1] = p.Appearance.Body
	p.Equips[2] = p.Appearance.Legs
	p.Inventory.Owner = p
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
	p.AddState(StateChangingLooks)
	p.SendPacket(OpenChangeAppearance)
}

//Chat sends a player NPC chat message net to the player and all other players around it.  If multiple msgs are
// provided, will sleep the goroutine for 1800ms between each message.
func (p *Player) Chat(msgs ...string) {
	for _, msg := range msgs {
		for _, player := range p.NearbyPlayers() {
			player.SendPacket(PlayerMessage(p, msg))
		}
		p.SendPacket(PlayerMessage(p, msg))

		//		if i < len(msgs)-1 {
		time.Sleep(time.Millisecond * 1920)
		// TODO: is 3 ticks right?
		//		}
	}
}

//OpenOptionMenu opens an option menu with the provided options, and returns the reply index, or -1 upon timeout..
func (p *Player) OpenOptionMenu(options ...string) int {
	// Can get option menu during most states, even fighting, but not trading, or if we're already in a menu...
	if p.IsPanelOpened() || p.HasState(StateMenu) {
		return -1
	}
	p.ReplyMenuC = make(chan int8)
	p.AddState(StateMenu)
	p.SendPacket(OptionMenuOpen(options...))

	select {
	case reply := <-p.ReplyMenuC:
		if !p.HasState(StateMenu) {
			return -1
		}
		p.RemoveState(StateMenu)
		close(p.ReplyMenuC)
		if reply < 0 || int(reply) > len(options)-1 {
			return -1
		}

		if p.HasState(StateChatting) {
			p.Chat(options[reply])
		}
		return int(reply)
	case <-time.After(time.Second * 60):
		if p.HasState(StateMenu) {
			p.RemoveState(StateMenu)
			close(p.ReplyMenuC)
			p.SendPacket(OptionMenuClose)
		}
		return -1
	}
}

//CloseOptionMenu closes any open option menus.
func (p *Player) CloseOptionMenu() {
	if p.HasState(StateMenu) {
		p.RemoveState(StateMenu)
		close(p.ReplyMenuC)
		p.SendPacket(OptionMenuClose)
	}
}

//CanWalk returns true if this player is in a state that allows walking.
func (p *Player) CanWalk() bool {
	if p.BusyInput() {
		return true
	}
	return !p.HasState(MSBatching, StateFighting, StateTrading, StateDueling, StateChangingLooks, StateSleeping, StateChatting, StateBusy, StateShopping)
}

//PlaySound sends a command to the client to play a sound by its file name.
func (p *Player) PlaySound(soundName string) {
	p.SendPacket(Sound(soundName))
}

//SendStat sends the information for the stat at idx to the player.
func (p *Player) SendStat(idx int) {
	p.SendPacket(PlayerStat(p, idx))
}

//SendStatExp sends the experience information for the stat at idx to the player.
func (p *Player) SendStatExp(idx int) {
	p.SendPacket(PlayerExperience(p, idx))
}

func (p *Player) SendCombatPoints() {
	p.SendPacket(PlayerCombatPoints(p))
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

//IncCurStat sets this players current stat at idx to Current(idx)+lvl and updates the client about it.
func (p *Player) IncCurStat(idx int, lvl int) {
	p.Skills().IncreaseCur(idx, lvl)
	p.SendStat(idx)
}

//SetCurStat sets this players current stat at idx to lvl and updates the client about it.
func (p *Player) IncExp(idx int, amt int) {
	//if idx <= 3 {
	//	p.Attributes.Inc("combatPoints", amt)
	//	p.SendCombatPoints()
	//	return
	//}
	amt *= 20
	p.Skills().IncExp(idx, amt)
	// TODO: Fatigue
	delta := entity.ExperienceToLevel(p.Skills().Experience(idx)) - p.Skills().Maximum(idx)
	if delta != 0 {
		p.PlaySound("advance")
		p.Message(fmt.Sprintf("@gre@You just advanced %d %v level!", delta, entity.SkillName(idx)))
		oldCombat := p.Skills().CombatLevel()
		p.Skills().IncreaseCur(idx, delta)
		p.Skills().IncreaseMax(idx, delta)
		p.SendStat(idx)
		if oldCombat != p.Skills().CombatLevel() {
			p.UpdateAppearance()
		}
	} else {
		p.SendStatExp(idx)
	}
}

//SetMaxStat sets this players maximum stat at idx to lvl and updates the client about it.
func (p *Player) SetMaxStat(idx int, lvl int) {
	p.Skills().SetMax(idx, lvl)
	p.Skills().SetExp(idx, entity.LevelToExperience(lvl))
	p.SendStat(idx)
}

//AddItem Adds amount of the item with specified id to the players inventory, if possible, and updates the client about it.
func (p *Player) AddItem(id, amount int) {
	if p.Inventory.CanHold(id, amount) {
		defer p.SendInventory()
	}
	stackSize := 1
	if ItemDefs[id].Stackable {
		stackSize = amount
	}
	for i := 0; i < amount; i += stackSize {
		if p.Inventory.Add(id, amount) < 0 {
			return
		}
	}
}

func (p *Player) PrayerActivated(idx int) bool {
	return p.VarBool("prayer"+strconv.Itoa(idx), false)
}

func (p *Player) PrayerOn(idx int) {
	if p.IsDueling() && !p.DuelPrayer() {
		p.Message("You cannot use prayer in this duel!")
		p.SendPrayers()
		return
	}
	if idx == 0 || idx == 3 || idx == 9 {
		p.PrayerOff(0)
		p.PrayerOff(3)
		p.PrayerOff(9)
	}
	if idx == 1 || idx == 4 || idx == 10 {
		p.PrayerOff(1)
		p.PrayerOff(4)
		p.PrayerOff(10)
	}
	if idx == 2 || idx == 5 || idx == 11 {
		p.PrayerOff(2)
		p.PrayerOff(5)
		p.PrayerOff(11)
	}
	p.SetVar("prayer"+strconv.Itoa(idx), true)
}

func (p *Player) PrayerOff(idx int) {
	p.SetVar("prayer"+strconv.Itoa(idx), false)
}

func (p *Player) SendPrayers() {
	p.SendPacket(PrayerStatus(p))
}

func (p *Player) Skulled() bool {
	return p.Attributes.VarInt("skullTime", 0) > 0
}

func (p *Player) SetSkulled(val bool) {
	if val {
		p.Attributes.SetVar("skullTime", TicksTwentyMin)
	} else {
		p.Attributes.UnsetVar("skullTime")
	}
	p.UpdateAppearance()
}
 
func (p *Player) ResetFighting() {
       defer p.Mob.ResetFighting()
       p.ResetDuel()
}

func (p *Player) StartCombat(target entity.MobileEntity) {
	if target.IsPlayer() {
		target.(*Player).PlaySound("underattack")
		if !p.IsDueling() {
			p.SetSkulled(true)
		}
	}
	target.SetRegionRemoved()
	p.Teleport(target.X(), target.Y())
	p.AddState(StateFighting)
	target.AddState(StateFighting)
	p.SetDirection(RightFighting)
	target.SetDirection(LeftFighting)
	p.Transients().SetVar("fightTarget", target)
	target.Transients().SetVar("fightTarget", p)
	curTick := 0
	attacker := entity.MobileEntity(p)
	defender := target
	//var defender entity.MobileEntity = target
	p.Tickables = append(p.Tickables, func() bool {
		if ptarget, ok := target.(*Player); (ok && !ptarget.Connected()) || !target.HasState(StateFighting) ||
			!p.HasState(StateFighting) || !p.Connected() || p.LongestDeltaCoords(target.X(), target.Y()) > 0 {
			// target is a disconnected player, we are disconnected,
			// one of us is not in a fight, or we are distanced somehow unexpectedly.  Kill tasks.
			// quickfix for possible bugs I imagined will exist
			p.ResetFighting()
			target.ResetFighting()
			return true
		}

		// One round per 2 ticks
		curTick++
		if curTick%2 == 0 {
			// TODO: tickables return tick delay count, e.g return 2 will wait 2 ticks and rerun, maybe??
			// would get ridda this per-tickable counter var paradigm
			return false
		}

		defer func() {
			attacker, defender = defender, attacker
		}()

		attacker.Transients().Inc("fightRound", 1)
		if p.PrayerActivated(12) && attacker.IsNpc() {
			return false
		}
		nextHit := int(math.Min(float64(defender.Skills().Current(entity.StatHits)), float64(attacker.MeleeDamage(defender))))
		defender.Skills().DecreaseCur(entity.StatHits, nextHit)
		if defender.Skills().Current(entity.StatHits) <= 0 {
			if attacker, ok := attacker.(*Player); ok {
				attacker.PlaySound("victory")
			}
			defender.Killed(attacker)
			return true
		}
		defender.Damage(nextHit)

		sound := "combat"
		// TODO: hit sfx (1/2/3) 1 is standard sound 2 is armor sound 3 is ghostly undead sound
		sound += "1"
		if nextHit > 0 {
			sound += "b"
		} else {
			sound += "a"
		}
		if attacker.IsPlayer() {
			attacker.(*Player).PlaySound(sound)
		}
		if defender.IsPlayer() {
			defender.(*Player).PlaySound(sound)
		}

		return false
	})
}

//Killed kills this player, dropping all of its items where it stands.
func (p *Player) Killed(killer entity.MobileEntity) {
	p.Transients().SetVar("deathTime", time.Now())
	p.PlaySound("death")
	p.SendPacket(Death)
	for i := 0; i < 14; i++ {
		p.PrayerOff(i)
	}
	for i := 0; i < 18; i++ {
		p.Skills().SetCur(i, p.Skills().Maximum(i))
	}
	p.SendPrayers()
	p.SendStats()
	p.SetDirection(North)

	deathItems := []*GroundItem{NewGroundItem(DefaultDrop, 1, p.X(), p.Y())}
	if !p.IsDueling() {
		keepCount := 0
		if p.PrayerActivated(8) {
			// protect item prayer
			keepCount++
		}
		if !p.Skulled() {
			keepCount += 3
		}
		deathItems = append(deathItems, p.Inventory.DeathDrops(keepCount)...)
	} else {
		p.DuelOffer.Lock.RLock()
		for _, i := range p.DuelOffer.List {
			deathItems = append(deathItems, NewGroundItem(i.ID, i.Amount, p.X(), p.Y()))
		}
		p.DuelOffer.Lock.RUnlock()
		if p.DuelTarget() != nil {
			p.DuelTarget().ResetDuel()
		}
		p.ResetDuel()
	}

	if killer != nil && killer.IsPlayer() {
		killer := killer.(*Player)
		killer.DistributeMeleeExp(int(math.Ceil(MeleeExperience(p) / 4.0)))
		killer.Message("You have defeated " + p.Username() + "!")
	}
	for i, v := range deathItems {
		// becomes universally visible on NPCs, or temporarily private otherwise
		if i == 0 || p.Inventory.RemoveByID(v.ID, v.Amount) > -1 {
			if killer != nil && killer.IsPlayer() {
				v.SetVar("belongsTo", killer.Transients().VarLong("username", 0))
			}
			AddItem(v)
		} else {
			log.Suspicious.Printf("Death item failed during removal: %v,%v owner:%v, killer:%v!\n", v.ID, v.Amount, p, killer)
		}
	}

	p.SendEquipBonuses()
	p.ResetFighting()
	p.SetSkulled(false)
	
	plane := p.Plane()
	p.SetLocation(SpawnPoint, true)
	if p.Plane() != plane {
		p.SendPlane()
	}
}

func (p *Player) NpcWithin(id int, rad int) *NPC {
	return NpcNearest(id, p.X(), p.Y())
}

//SendPlane sends the current plane of this player.
func (p *Player) SendPlane() {
	p.SendPacket(PlaneInfo(p))
}

//SendEquipBonuses sends the current equipment bonuses of this player.
func (p *Player) SendEquipBonuses() {
	p.SendPacket(EquipmentStats(p))
}

//Damage sends a player damage bubble for this player to itself and any nearby players.
func (p *Player) Damage(amt int) {
	for _, player := range p.NearbyPlayers() {
		player.SendPacket(PlayerDamage(p, amt))
	}
	p.SendPacket(PlayerDamage(p, amt))
}

//ItemBubble sends an item action bubble for this player to itself and any nearby players.
func (p *Player) ItemBubble(id int) {
	for _, player := range p.NearbyPlayers() {
		player.SendPacket(PlayerItemBubble(p, id))
	}
	p.SendPacket(PlayerItemBubble(p, id))
}

//SetStat sets the current, maximum, and experience levels of the skill at idx to lvl, and updates the client about it.
func (p *Player) SetStat(idx, lvl int) {
	p.Skills().SetCur(idx, lvl)
	p.Skills().SetMax(idx, lvl)
	p.Skills().SetExp(idx, entity.LevelToExperience(lvl))
	p.SendStat(idx)
}

func (p *Player) CurrentShop() *Shop {
	if !p.Contains("shop") {
		return nil
	}
	return p.VarChecked("shop").(*Shop)
}

//OpenBank opens a shop screen for the player and sets the appropriate state variables.
func (p *Player) OpenShop(shop *Shop) {
	if p.IsFighting() || p.IsTrading() || p.IsDueling() || p.HasState(StateShopping) || p.HasState(StateBanking) {
		return
	}
	p.AddState(StateShopping)
	shop.Players.Add(p)
	p.Transients().SetVar("shop", shop)
	p.SendPacket(ShopOpen(shop))
}

//CloseBank closes the bank screen for this player and sets the appropriate state variables
func (p *Player) CloseShop() {
	if !p.HasState(StateShopping) {
		return
	}
	p.RemoveState(StateShopping)
	p.CurrentShop().Players.Remove(p)
	p.Transients().UnsetVar("shop")
	p.SendPacket(ShopClose)
}

//OpenBank opens a bank screen for the player and sets the appropriate state variables.
func (p *Player) OpenBank() {
	if p.IsFighting() || p.IsTrading() || p.IsDueling() || p.HasState(StateShopping) || p.HasState(StateBanking) {
		return
	}
	p.AddState(StateBanking)
	p.SendPacket(BankOpen(p))
}

//CloseBank closes the bank screen for this player and sets the appropriate state variables
func (p *Player) CloseBank() {
	if !p.HasState(StateBanking) {
		return
	}
	p.RemoveState(StateBanking)
	p.SendPacket(BankClose)
}

//SendUpdateTimer sends a system update countdown timer to the client.
func (p *Player) SendUpdateTimer() {
	p.SendPacket(SystemUpdate(int(time.Until(UpdateTime).Seconds())))
}

func (p *Player) SendMessageBox(msg string, big bool) {
	if big {
		p.SendPacket(BigInformationBox(msg))
	} else {
		p.SendPacket(InformationBox(msg))
	}
}

func (p *Player) SetCache(name string, val interface{}) {
	p.Attributes.SetVar(name, val)
}

func (p *Player) RemoveCache(name string) {
	p.Attributes.UnsetVar(name)
}

func (p *Player) Cache(name string) interface{} {
	if !p.Attributes.Contains(name) {
		return nil
	}
	return p.Attributes.VarChecked(name)
}

func (p *Player) OpenSleepScreen() {
	p.AddState(StateSleeping)
	p.SendPacket(SleepWord(p))
}
