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
	"bufio"
	"fmt"
	"io"
	"math"
	stdnet "net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/tasks"
	"github.com/spkaeros/rscgo/pkg/errors"
	"github.com/spkaeros/rscgo/pkg/game/entity"
	"github.com/spkaeros/rscgo/pkg/game/net"
	"github.com/spkaeros/rscgo/pkg/game/social"
	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

type (
	//Player A player in our game world.
	Player struct {
		LocalPlayers     *MobList
		LocalNPCs        *MobList
		LocalObjects     *entityList
		LocalItems       *entityList
		FriendList       *social.FriendsList
		IgnoreList       []uint64
		Appearance       entity.AppearanceTable
		KnownAppearances map[int]int
		AppearanceReq    []*Player
		Socket           stdnet.Conn
		Attributes       *entity.AttributeList
		Inventory        *Inventory
		bank             *Inventory
		TradeOffer       *Inventory
		DuelOffer        *Inventory
		Duel             struct {
			Rules    [4]bool
			Accepted [2]bool
			Target   *Player
		}
		Tickables     tasks.Scripts
		PostTickables tasks.Scripts
		tickAction    tasks.StatusReturnCall
		ActionLock    sync.RWMutex
		ReplyMenuC    chan int8
		killer        sync.Once
		hasReader     bool
		Websocket     bool
		InQueue       chan *net.Packet
		Reader        *bufio.Reader
		Writer		  net.WriteFlusher
		DatabaseIndex int
		Mob
	}
)

func (p *Player) UsernameHash() uint64 {
	if p == nil {
		return 0
	}
	return p.VarLong("username", strutil.Base37.Encode("nil"))
}

func (p *Player) Bank() *Inventory {
	return p.bank
}

func (p *Player) CanAttack(target entity.MobileEntity) bool {
	if target := AsNpc(target); target != nil {
		return target.Attackable()
	}
	if p.State()&StateFightingDuel == StateFightingDuel {
		return p.Duel.Target == target && p.DuelMagic()
	}
	targetp := AsPlayer(target)
	ourWild := p.Wilderness()
	targetWild := targetp.Wilderness()
	if ourWild < 1 || targetWild < 1 {
		p.Message("You cannot attack other players outside of the wilderness!")
		return false
	}
	delta := p.CombatDelta(target)
	if delta < 0 {
		delta = -delta
	}
	if delta > ourWild {
		p.Message("You must move to at least level " + strconv.Itoa(delta) + " wilderness to attack " + targetp.Username() + "!")
		return false
	}
	if delta > targetWild {
		p.Message(targetp.Username() + " is not in high enough wilderness for you to attack!")
		return false
	}
	return true
}

func (p *Player) Username() string {
	return strutil.Base37.Decode(p.VarLong("username", strutil.Base37.Encode("NIL")))
}

//CurrentIP returns the remote IP address this player connected from
func (p *Player) CurrentIP() string {
	return strings.Split(p.RemoteAddress(), ":")[0]
}

//LocalAddress Returns the remote IP:port that this player connected from, or N/A if this player never connected somehow
func (p *Player) RemoteAddress() string {
	if p.Socket == nil {
		return "N/A:N/A"
	}
	return p.Socket.RemoteAddr().String()
}

func (p *Player) IsWebsocket() bool {
	return p.Websocket
}

func (p *Player) Rank() int {
	return p.VarInt("rank", 0)
}

func (p *Player) AppearanceTicket() int {
	return p.VarInt("appearanceTicket", 0)
}

//String returns a string populated with the more identifying features of this player.
func (p *Player) String() string {
	return fmt.Sprintf("'%s'[%d]@'%s'", p.Username(), p.Index, p.CurrentIP())
}

//SetTickAction queues a distanced action to run every game engine tick before path traversal, if action returns true, it will be reset.
func (p *Player) SetTickAction(action func() bool) {
	p.Lock()
	defer p.Unlock()
	p.tickAction = action
}

func (p *Player) TickAction() func() bool {
	p.RLock()
	defer p.RUnlock()
	return p.tickAction
}

//ResetTickAction clears the distanced action, if any is queued.  Should be called any time the player is deliberately performing an action.
func (p *Player) ResetTickAction() {
	p.SetTickAction(nil)
}

//FriendsWith returns true if specified username is in our friend entityList.
func (p *Player) FriendsWith(other uint64) bool {
	return p.FriendList.ContainsHash(other)
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
		if player.FriendList.Contains(p.Username()) {
			if p.FriendList.Contains(player.Username()) || !p.FriendBlocked() {
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
func (p *Player) WalkingArrivalAction(t entity.MobileEntity, dist int, action func()) {
	p.SetTickAction(func() bool {
		if p.Near(t, dist) {
			if !p.ReachableCoords(t.X(), t.Y()) {
				return true
			}
			action()
			return false
		}
		return p.WalkTo(NewLocation(t.X(), t.Y()))
	})
}

//CanReachMob Check if we can reach a mob traversing the most direct tiles toward them, e.g straight lines.
// Used to check ranged combat attacks, or trade requests, basically anything needing local interactions.
func (p *Player) CanReachMob(target entity.MobileEntity) bool {
	if p.FinishedPath() && p.VarInt("triedReach", 0) >= 5 {
		// Tried reaching one mob >=5 times without single success, abort early.
		//		p.ResetPath()
		p.UnsetVar("triedReach")
		return false
	}
	p.Inc("triedReach", 1)

	pathX, pathY := p.X(), p.Y()
	for steps := 0; steps < 256; steps++ {
		if !p.ReachableCoords(pathX, pathY) {
			return false
		}
		// check deltas
		if pathX == target.X() && pathY == target.Y() {
			p.UnsetVar("triedReach")
			return true
		}

		// Update coords toward target in a straight line
		if pathX < target.X() {
			pathX++
		} else if pathX > target.X() {
			pathX--
		}

		if pathY < target.Y() {
			pathY++
		} else if pathY > target.Y() {
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

	// p.UpdateStatus(true)

	Players.Range(func(player *Player) {
		if player.FriendList.Contains(p.Username()) {
			if !p.FriendList.Contains(player.Username()) || p.FriendBlocked() {
				p.SendPacket(FriendUpdate(p.UsernameHash(), !p.FriendBlocked() || p.FriendList.Contains(player.Username())))
			}
		}
	})
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

//NextTo returns true if we can walk a straight line to target without colliding with any walls or objects,
// otherwise returns false.
func (p *Player) NextTo(target Location) bool {
	return p.Reachable(target)
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

	if !p.Reachable(dst) {
		p.ResetPath()
		return
	}

	p.SetLocation(dst, false)
}

func (l Location) Blocked() bool {
	return false
}

//Targetable returns true if you are able to see the other location from the receiever location without hitting
// any obstacles, and you are within range.  Otherwise returns false.
func (l Location) Targetable(other Location) bool {
	return l.WithinRange(other, 5) && l.Reachable(other)
}

//WithinReach returns true if you are able to physically touch the other person you are so close without obstacles
// Otherwise returns false.
func (l Location) WithinReach(other Location) bool {
	return l.WithinRange(other, 2) && l.Reachable(other)
}

func (l Location) Reachable(other Location) bool {
	return l.ReachableCoords(other.X(), other.Y())
}
func (l Location) RD(x, y int) bool {
	if l.ReachableCoords(x, y) {
		log.Debug(x,y,"reachable!")
		return true
	}
	log.Debug("Unreachable:{from:", l.String(), "to:", x,y)
	return false
}

func (l Location) ReachableCoords(x, y int) bool {
	check := func(l, dst Location) bool {
		bitmask := byte(ClipBit(l.DirectionToward(dst)))
		dstmask := byte(ClipBit(dst.DirectionToward(l)))
		// check mask of our tile and dst tile
		if IsTileBlocking(l.X(), l.Y(), bitmask, true) || IsTileBlocking(dst.X(), dst.Y(), dstmask, false) {
			return false
		}

		return true
	}
	cur := l.Clone()
	end := NewLocation(x, y)

	for !cur.Equals(end) {
		next := cur.Step(cur.DirectionToward(end))
		if !check(cur, next) {
			return false
		}
		cur = next
	}
	return true
}

//UpdateRegion if this player is currently in a region, removes it from that region, and adds it to the region at x,y
func (p *Player) UpdateRegion(x, y int) {
	curArea := Region(p.X(), p.Y())
	newArea := Region(x, y)
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
	reqs := definitions.Items[item.ID].Requirements
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
	def := definitions.Equip(item.ID)
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
		otherDef := definitions.Equip(otherItem.ID)
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
		if otherDef.Type&1 == 1 {
			p.Equips()[otherDef.Position] = p.Appearance.Head
		} else if otherDef.Type&2 == 2 {
			p.Equips()[otherDef.Position] = p.Appearance.Body
		} else if otherDef.Type&4 == 4 {
			p.Equips()[otherDef.Position] = p.Appearance.Legs
		} else {
			p.Equips()[otherDef.Position] = 0
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
	p.Equips()[def.Position] = def.Sprite
	// Below will update our appearance ticket (a counter for the number of times this session we updated this player)
	p.UpdateAppearance()
	p.WritePacket(EquipmentStats(p))
	p.WritePacket(InventoryItems(p))
}

func (p *Player) Equips() []int {
	s, ok := p.Var("sprites")
	if !ok || s == nil {
		return []int{}
	}
	return s.([]int)
}

func (p *Player) UpdateAppearance() {
	p.SetAppearanceChanged()
	p.SetVar("appearanceTicket", p.AppearanceTicket()+1)
}

//DequipItem removes an item from this players equips, and sends inventory and equipment bonuses.
func (p *Player) DequipItem(item *Item) {
	def := definitions.Equip(item.ID)
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
	if def.Type&1 == 1 {
		p.Equips()[def.Position] = p.Appearance.Head
	} else if def.Type&2 == 2 {
		p.Equips()[def.Position] = p.Appearance.Body
	} else if def.Type&4 == 4 {
		p.Equips()[def.Position] = p.Appearance.Legs
	} else {
		p.Equips()[def.Position] = 0
	}
	p.UpdateAppearance()
	p.SendEquipBonuses()
	p.SendInventory()
}

//ResetAll in order, calls ResetFighting, ResetTrade, ResetTickAction, ResetFollowing, and CloseOptionMenu.
func (p *Player) ResetAll() {
	p.ResetFighting()
	p.ResetDuel()
	p.ResetTrade()
	p.ResetTickAction()
	p.CloseOptionMenu()
	p.CloseBank()
	p.CloseShop()
}

func (p *Player) ResetAllExceptDueling() {
	p.ResetTrade()
	p.ResetTickAction()
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
	for _, r := range Region(p.X(), p.Y()).neighbors() {
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
	for _, r := range Region(p.X(), p.Y()).neighbors() {
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
	for _, r := range Region(p.X(), p.Y()).neighbors() {
		objects = append(objects, r.Objects.NearbyObjects(p)...)
	}

	return
}

//NewObjects Returns nearby objects that this player is unaware of.
func (p *Player) NewObjects() (objects []*Object) {
	for _, r := range Region(p.X(), p.Y()).neighbors() {
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
	for _, r := range Region(p.X(), p.Y()).neighbors() {
		for _, i := range r.Items.NearbyItems(p) {
			if !p.LocalItems.Contains(i) {
				items = append(items, i)
			}
		}
	}

	return
}

//NewPlayers Returns nearby players that this player is unaware of.
func (p *Player) NewPlayers() (players *MobList) {
	list := &MobList{}
	for _, r := range Region(p.X(), p.Y()).neighbors() {
		r.Players.RangePlayers(func(p1 *Player) bool {
			if !p.LocalPlayers.Contains(p1) && p != p1 && p.WithinRange(p1.Location, 15) {
				list.Add(p1)
			}
			return false
		})
	}

	return list
}

//NewNPCs Returns nearby NPCs that this player is unaware of.
func (p *Player) NewNPCs() (npcs []*NPC) {
	for _, r := range Region(p.X(), p.Y()).neighbors() {
		r.NPCs.RangeNpcs(func(n *NPC) bool {
			if !n.VarBool("removed", false) && !p.LocalNPCs.Contains(n) && p.WithinRange(n.Location, 15) {
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
	return p.Duel.Accepted[screen-1]
}

//DuelAccepted returns the status of the specified duel negotiation screens accepted button for this player.
// Valid screens are 1 and 2.
func (p *Player) SetDuelAccepted(screen int, b bool) {
	if b && screen == 2 && !p.Duel.Accepted[0] {
		log.Suspicious.Println("Attempt to set duelaccept2 before duelaccept1:", p.String())
		return
	}
	p.Duel.Accepted[screen-1] = b
}

//SetDuelRule sets the duel rule associated with the specified index to b.
// Valid rule indices are 0 through 3.
func (p *Player) SetDuelRule(index int, b bool) {
	p.Duel.Rules[index] = !b
}

//DuelRule returns the rule associated with the specified index provided.
// Valid rule indices are 0 through 3.
func (p *Player) duelRule(index int) bool {
	return p.Duel.Rules[index]
}

func (p *Player) DuelRetreating() bool {
	return p.Duel.Rules[0]
}

func (p *Player) DuelMagic() bool {
	return p.Duel.Rules[1]
}

func (p *Player) DuelPrayer() bool {
	return p.Duel.Rules[2]
}

func (p *Player) DuelEquipment() bool {
	return p.Duel.Rules[3]
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
	p.Duel.Target = p1
}

//ResetDuelTarget Removes receivers duel target, if any.
func (p *Player) ResetDuelTarget() {
	p.Duel.Target = nil
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

//SendPacket sends a net to the client.
func (p *Player) SendPacket(packet *net.Packet) {
	if p == nil || (!p.Connected() && packet.Opcode != 0) {
		return
	}
	p.WritePacket(packet)
}

type PlayerService interface {
	PlayerSave(*Player)
}

var DefaultPlayerService PlayerService

//Destroy sends a kill signal to the underlying client to tear down all of the I/O routines and save the player.
func (p *Player) Destroy() {
	p.WritePacket(Logout)
	p.Writer.Flush()
	p.PostTickables.Add(func() bool {
		p.killer.Do(func() {
			if err := p.Socket.Close(); err != nil {
				log.Warn("Couldn't close socket:", err)
			}
			p.Attributes.SetVar("lastIP", p.CurrentIP())
			p.Inventory.Owner = nil
			if Players.Find(p) > -1 {
				log.Debug("Unregistered:{'" + p.Username() + "'@'" + p.CurrentIP() + "'}")
				p.ResetAll()
				p.UpdateStatus(false)
				p.SetConnected(false)
				go func() {
					DefaultPlayerService.PlayerSave(p)
					RemovePlayer(p)
				}()
				return
			}
			log.Debug("Unregistered:{'" + p.CurrentIP() + "'}")
		})
		return true
	})
}

func (p *Player) AtObject(object *Object) bool {
	bounds := object.Boundaries()
	if definitions.ScenaryObjects[object.ID].CollisionType == 2 || definitions.ScenaryObjects[object.ID].CollisionType == 3 {
		// door types
		return (p.Reachable(bounds[0]) || p.Reachable(bounds[1])) && p.WithinArea(bounds)
	}

	// TODO: Maybe replace this with the following:
	//	return (p.Reachable(bounds[0]), bounds[1]) || p.Reachable(bounds[1])) || (p.FinishedPath() && p.CanReachDiag(bounds))
	//	return p.Reachable(bounds[0]) || p.Reachable(bounds[1]) && p.WithinArea(bounds) p.CanReach(bounds) ||  (p.FinishedPath() && p.CanReachDiag(bounds))

	return p.CanReach(bounds) || (p.FinishedPath() && p.CanReachDiag(bounds))
}

func (p *Player) CanReachDiag(bounds [2]Location) bool {
	/*
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

		low, high := p.Masks(bounds[0].X(), bounds[0].Y()), p.Masks(bounds[1].X(), bounds[1].Y())
		mixedNS, mixedEW := p.Masks(bounds[0].X(), bounds[1].Y()), p.Masks(bounds[1].X(), bounds[0].Y())
	*/
	/*
		tile := p.Location.Clone()
		lowX, lowY := p.X() - bounds[0].X(), p.Y() - bounds[0].Y()
		highX, highY := p.X() - bounds[1].X(), p.Y() - bounds[1].Y()
		masks := byte(0)
		if (lowX == 0 || highX == 0) && (lowY == 0 || highY == 0  {
			if lowX < 0 {
				masks |= ClipWest
			}
			if lowY < 0 {
				masks |=
			}
		}
		if lowX >= 1 && highX <= -1 {
			masks |= ClipWest
			tile.x.Dec()
		} else if lowX <= -1  && highX >= 1 {
			masks |= ClipEast
			tile.x.Inc()
		}
		if lowY >= 1 {
			masks |= ClipSouth
			tile.y.Dec()
		} else if lowY <= -1 {
			masks |= ClipNorth
			tile.y.Inc()
		}
		return CollisionData&masks != 0
	*/
	// Northeast target
	if p.X()-1 >= bounds[0].X() && p.X()-1 <= bounds[1].X() && p.Y()-1 >= bounds[0].Y() && p.Y()-1 <= bounds[1].Y() {
		return CollisionData(p.X()-1, p.Y()-1)&(ClipSouth|ClipWest) == 0
	}
	// Northwest target
	if p.X()+1 >= bounds[0].X() && p.X()+1 <= bounds[1].X() && p.Y()-1 >= bounds[0].Y() && p.Y()-1 <= bounds[1].Y() {
		return CollisionData(p.X()-1, p.Y()-1)&(ClipSouth|ClipEast) == 0
	}
	// Southeast target
	if p.X()-1 >= bounds[0].X() && p.X()-1 <= bounds[1].X() && p.Y()+1 >= bounds[0].Y() && p.Y()+1 <= bounds[1].Y() {
		return CollisionData(p.X()-1, p.Y()-1)&(ClipNorth|ClipWest) == 0
	}
	// Southwest target
	if p.X()+1 >= bounds[0].X() && p.X()+1 <= bounds[1].X() && p.Y()+1 >= bounds[0].Y() && p.Y()+1 <= bounds[1].Y() {
		return CollisionData(p.X()-1, p.Y()-1)&(ClipNorth|ClipEast) == 0
	}
	return false
	/*
		// Southeast target
		if lowX >= 1 && lowY <= -1 {
			return CollisionData(p.X()-1, p.Y()-1).CollisionMask&(ClipNorth|ClipWest) == 0
		}
		// Southwest target
		if lowX <= -1 && lowY <= -1 {
			return CollisionData(p.X()-1, p.Y()-1).CollisionMask&(ClipNorth|ClipEast)
		}
		// Northeast target
		if lowX >= 1 && lowY >= 1 {
			return CollisionData(p.X()-1, p.Y()-1).CollisionMask&(ClipSouth|ClipWest) == 0
		}
		// Northwest target
		if lowX <= -1 && lowY >= 1 {
			return CollisionData(p.X()-1, p.Y()-1).CollisionMask&(ClipSouth|ClipEast)
		}
		return CollisionData(tile.X(), tile.Y()).CollisionMask&() == 0
	*/
}

//Initialize informs the client of all of the various attributes of this player, and starts the stat normalization
// routine.
func (p *Player) Initialize() {
	AddPlayer(p)
	// Mark down time of authentication
	p.SetVar("authTime", time.Now())
	// update flags
	p.SetConnected(true)
	p.SetAppearanceChanged()
	p.SetSpriteUpdated()

	// settings panel
	p.WritePacket(ClientSettings(p))
	p.WritePacket(PrivacySettings(p))
	// social panel
	p.WritePacket(FriendList(p))
	p.WritePacket(IgnoreList(p))
	// TODO: Not canonical RSC, but definitely good QoL update...
	//  p.SendPacket(FightMode(p))

	// stat panel
	p.WritePacket(PlayerStats(p))
	p.WritePacket(Fatigue(p))
	p.WritePacket(EquipmentStats(p))
	p.SendCombatPoints()

	// inventory panel
	p.WritePacket(InventoryItems(p))
	// mesh related coordinate info and player index
	p.SendPlane()
	if !p.Attributes.Contains("madeAvatar") {
		p.OpenAppearanceChanger()
	} else {
		if !p.Reconnecting() {
			p.SendPacket(WelcomeMessage)
			p.SendPacket(LoginBox(int(time.Since(p.Attributes.VarTime("lastLogin")).Hours()/24), p.Attributes.VarString("lastIP", "0.0.0.0")))
		}
		p.Attributes.SetVar("lastLogin", time.Now())
	}
	for _, fn := range LoginTriggers {
		fn(p)
	}
}
func (p *Player) WritePacket(packet *net.Packet) {
	defer p.Writer.Flush()
	if packet.Opcode == 0 {
		p.Writer.Write(packet.FrameBuffer)
		return
	}
	header := []byte{0, 0}
	frameLength := len(packet.FrameBuffer)
	if frameLength >= 160 {
		header[0] = byte(frameLength>>8 + 160)
		header[1] = byte(frameLength)
	} else {
		header[0] = byte(frameLength)
		if frameLength > 0 {
			frameLength--
			header[1] = packet.FrameBuffer[frameLength]
		}
	}
	p.Writer.Write(append(header, packet.FrameBuffer[:frameLength]...))
}

//NewPlayer Returns a reference to a new player.
func NewPlayer(socket stdnet.Conn) *Player {
	p := &Player{
		Socket: socket,
		Mob: Mob{
			Entity: Entity{
				Location: Lumbridge.Clone(),
			},
			AttributeList: entity.NewAttributeList(),
		},
		Attributes:       entity.NewAttributeList(),
		LocalPlayers:     NewMobList(),
		LocalNPCs:        NewMobList(),
		LocalObjects:     &entityList{},
		LocalItems:       &entityList{},
		Appearance:       entity.DefaultAppearance(),
		FriendList:       social.New(),
		KnownAppearances: make(map[int]int),
		bank:             &Inventory{Capacity: 48 * 4, stackEverything: true},
		Inventory:        &Inventory{Capacity: 30},
		TradeOffer:       &Inventory{Capacity: 12},
		DuelOffer:        &Inventory{Capacity: 8},
		InQueue:          make(chan *net.Packet, 50),
		Reader:			  bufio.NewReader(socket),
	}
	// TODO: Get rid of this self-referential member; figure out better way to handle client item updating
	p.Inventory.Owner = p
	p.SetVar("sprites", []int{entity.DefaultAppearance().Head, entity.DefaultAppearance().Body, entity.DefaultAppearance().Legs, -1, -1, -1, -1, -1, -1, -1, -1, -1})
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

//Chat sends a player NPC chat message packet to the player and all other players around it.  If multiple msgs are
// provided, will sleep the goroutine for 3-4 ticks between each message, depending on length of message.
func (p *Player) Chat(msgs ...string) {
	for _, msg := range msgs {
		sleep := 3
		if len(msg) >= 83 {
			sleep = 4
		}
		m := p.TargetMob()
		for _, player := range p.NearbyPlayers() {
			player.QueueQuestChat(p, m, msg)
		}
		p.QueueQuestChat(p, m, msg)
		time.Sleep(time.Millisecond*640*time.Duration(sleep))
	}
}

//QueuePublicChat Adds a message to a locked public-chat queue
func (p *Player) QueuePublicChat(owner entity.MobileEntity, message string) {
	if !p.Contains("publicChatQ") {
		p.SetVar("publicChatQ", []ChatMessage{NewChatMessage(owner, message)})
		return
	}
	p.SetVar("publicChatQ", append(p.VarChecked("publicChatQ").([]ChatMessage), NewChatMessage(owner, message)))
}

//QueueQuestChat Adds a message to a locked quest-chat queue
func (p *Player) QueueQuestChat(owner, target entity.MobileEntity, message string) {
	if !p.Contains("questChatQ") {
		p.SetVar("questChatQ", []ChatMessage{NewChatMessage(owner, message)})
		return
	}
	p.SetVar("questChatQ", append(p.VarChecked("questChatQ").([]ChatMessage), NewChatMessage(owner, message)))
}

//QueueNpcChat Adds a message to a locked quest-chat queue
func (p *Player) QueueNpcChat(owner, target entity.MobileEntity, message string) {
	if !p.Contains("npcChatQ") {
		p.SetVar("npcChatQ", []ChatMessage{NewTargetedMessage(owner, target, message)})
		return
	}
	p.SetVar("npcChatQ", append(p.VarChecked("npcChatQ").([]ChatMessage), NewTargetedMessage(owner, target, message)))
}

//QueueNpcSplat Adds a message to a locked quest-chat queue
func (p *Player) QueueNpcSplat(owner *NPC, dmg int) {
	if !p.Contains("npcSplatQ") {
		p.SetVar("npcSplatQ", []HitSplat{NewHitsplat(owner, dmg)})
		return
	}
	p.SetVar("npcSplatQ", append(p.VarChecked("npcSplatQ").([]HitSplat), NewHitsplat(owner, dmg)))
}

//QueueProjectile Adds a missile to a locked projectile queue
func (p *Player) QueueProjectile(owner, target entity.MobileEntity, kind int) {
	if !p.Contains("projectileQ") {
		p.SetVar("projectileQ", []Projectile{NewProjectile(owner, target, kind)})
		return
	}
	p.SetVar("projectileQ", append(p.VarChecked("projectileQ").([]Projectile), NewProjectile(owner, target, kind)))
}

//QueueHitsplat Adds a hit splat to a locked hit-splat queue
func (p *Player) QueueHitsplat(owner entity.MobileEntity, dmg int) {
	if !p.Contains("hitsplatQ") {
		p.SetVar("hitsplatQ", []HitSplat{NewHitsplat(owner, dmg)})
		return
	}
	p.SetVar("hitsplatQ", append(p.VarChecked("hitsplatQ").([]HitSplat), NewHitsplat(owner, dmg)))
}

//QueueItemBubble Adds an action item bubble to a locked item bubble queue
func (p *Player) QueueItemBubble(owner *Player, id int) {
	if !p.Contains("bubbleQ") {
		p.SetVar("bubbleQ", []ItemBubble{{owner, id}})
		return
	}
	p.SetVar("bubbleQ", append(p.VarChecked("bubbleQ").([]ItemBubble), ItemBubble{owner, id}))
}

func (p *Player) OpenOptionMenu(options ...string) int {
	// Can get option menu during most states, even fighting, but not trading, or if we're already in a menu...
	if p.IsPanelOpened() || p.HasState(StateMenu) {
		return -1
	}
	p.AddState(StateMenu)
	p.SendPacket(OptionMenuOpen(options...))
	p.ReplyMenuC = make(chan int8)
	select {
	case r, ok := <-p.ReplyMenuC:
		if !p.HasState(StateMenu) || !ok {
			return -1
		}
		close(p.ReplyMenuC)
		p.RemoveState(StateMenu)
		if r < 0 || r > int8(len(options)-1) {
			log.Warn("Invalid option menu reply:", r)
			return -1
		}

		if p.TargetNpc() != nil && p.HasState(StateChatting) {
			p.Chat(options[r])
		}
		return int(r)
	}
	return -1
}

//CloseOptionMenu closes any open option menus.
func (p *Player) CloseOptionMenu() {
	if p.HasState(StateMenu) {
		p.RemoveState(StateMenu)
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
	if delta > 0 {
		p.PlaySound("advance")
		p.Message("@gre@You just advanced " + strconv.Itoa(delta) + " " + entity.SkillName(idx) + " level!")
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
	p.Skills().SetExp(idx, entity.LevelToExperience(lvl) / 4)
	p.SendStat(idx)
}

//AddItem Adds amount of the item with specified id to the players inventory, if possible, and updates the client about it.
func (p *Player) AddItem(id, amount int) {
	if p.Inventory.CanHold(id, amount) {
		defer p.SendInventory()
	}
	stackSize := 1
	if definitions.Items[id].Stackable {
		stackSize = amount
	}
	for i := 0; i < amount; i += stackSize {
		if p.Inventory.Add(id, amount) < 0 {
			return
		}
	}
}

func (p *Player) TogglePrayer(idx int) bool {
	p.Mob.Prayers[idx] = !p.Mob.Prayers[idx]
	return p.Mob.Prayers[idx]
}

func (p *Player) ActivatePrayer(idx int) {
	p.Mob.Prayers[idx] = true
}

func (p *Player) PrayerOn(idx int) {
	if p.IsDueling() && !p.DuelPrayer() {
		p.Message("You cannot use prayer in this duel!")
		p.SendPrayers()
		return
	}
	boosterPrayers := [3][3]int{
		{0, 3, 9},
		{1, 4, 10},
		{2, 5, 11},
	}
	defer p.ActivatePrayer(idx)
	for stat := 0; stat < 3; stat++ {
		for _, i := range boosterPrayers[stat] {
			if i == idx {
				for _, i1 := range boosterPrayers[stat] {
					p.Mob.Prayers[i1] = false
				}
				return
			}
		}
	}
}

func (p *Player) PrayerOff(idx int) {
	p.DeactivatePrayer(idx)
}

func (p *Player) DeactivatePrayer(idx int) {
	p.Mob.Prayers[idx] = false
}

func (p *Player) PrayerActivated(i int) bool {
	return p.Mob.Prayers[i]
}

func (p *Player) SendPrayers() {
	p.SendPacket(PrayerStatus(p))
}

func (p *Player) Skulled() bool {
	return p.Attributes.VarInt("skullTicks", 0) > 0
}

func (p *Player) SetSkulled(val bool) {
	if val {
		p.Attributes.SetVar("skullTicks", TicksTwentyMin)
	} else {
		p.Attributes.UnsetVar("skullTicks")
	}
	p.UpdateAppearance()
}

func (p *Player) ResetFighting() {
	defer p.Mob.ResetFighting()
	p.ResetDuel()
}

func (p *Player) Skulls() map[uint64]time.Time {
	records, ok := p.Var("attackedList")
	if !ok || records == nil {
		p.SetVar("attackedList", make(map[uint64]time.Time))
		records = p.VarChecked("attackedList").(map[uint64]time.Time)
	}
	return records.(map[uint64]time.Time)
}

func (p *Player) SkullOn(p1 *Player) {
	p.AddSkull(p1.UsernameHash())
}

func (p *Player) SkulledOn(user uint64) bool {
	t, ok := p.Skulls()[user]
	return ok && time.Since(t) <= time.Minute*time.Duration(20)
}

func (p *Player) AddSkull(user uint64) {
	if p.SkulledOn(user) {
		// we skulled on them within 20 mins ago, ignore call
		return
	}
	p.SetSkulled(true)
	p.Skulls()[user] = time.Now()
}

func AsPlayer(m entity.MobileEntity) *Player {
	if p, ok := m.(*Player); ok {
		return p
	}
	return nil
}

func AsNpc(m entity.MobileEntity) *NPC {
	if n, ok := m.(*NPC); ok {
		return n
	}
	return nil
}

func (p *Player) StartCombat(defender entity.MobileEntity) {
	attacker := entity.MobileEntity(p)
	if targetp := AsPlayer(defender); targetp != nil {
		targetp.PlaySound("underattack")
		if !p.IsDueling() && !targetp.SkulledOn(p.UsernameHash()) {
			p.SkullOn(targetp)
		}
	}
	p.SetVar("targetMob", defender)
	p.SetVar("fightTarget", defender)
	defender.SessionCache().SetVar("fightTarget", p)
	defender.SetRegionRemoved()
	p.Teleport(defender.X(), defender.Y())
	p.AddState(StateFighting)
	defender.AddState(StateFighting)
	p.SetDirection(RightFighting)
	defender.SetDirection(LeftFighting)
	tasks.Schedule(2, func() bool {
		if (defender.IsPlayer() && !AsPlayer(defender).Connected()) || !defender.HasState(StateFighting) ||
			!p.HasState(StateFighting) || !p.Connected() || p.LongestDeltaCoords(defender.X(), defender.Y()) > 0 {
			// target is a disconnected player, we are disconnected,
			// one of us is not in a fight, or we are distanced somehow unexpectedly.  Kill tasks.
			// quickfix for possible bugs I imagined will exist
			p.ResetFighting()
			defender.ResetFighting()
			return true
		}
		defer func() {
			attacker, defender = defender, attacker
		}()
		attacker.SessionCache().Inc("fightRound", 1)

		// Paralyze Monster blocker here
		if attacker.IsNpc() && defender.IsPlayer() && AsPlayer(defender).PrayerActivated(12) {
			return false
		}
		
		nextHit := int(math.Min(float64(defender.Skills().Current(entity.StatHits)), float64(attacker.MeleeDamage(defender))))
		defender.Skills().DecreaseCur(entity.StatHits, nextHit)
		if defender.IsNpc() && attacker.IsPlayer() {
			AsNpc(defender).CacheDamage(AsPlayer(attacker).UsernameHash(), nextHit)
		}
		defender.Damage(nextHit)
		if defender.Skills().Current(entity.StatHits) <= 0 {
			if attackerp := AsPlayer(attacker); attackerp != nil {
				attackerp.PlaySound("victory")
			}
			defender.Killed(attacker)
			return true
		}

		sound := "combat"
		// TODO: hit sfx (1/2/3) 1 is standard sound 2 is armor sound 3 is ghostly undead sound
		sound += "1"
		if nextHit > 0 {
			sound += "b"
		} else {
			sound += "a"
		}
		
		if attackerp := AsPlayer(attacker); attackerp != nil {
			attackerp.PlaySound(sound)
		}
		
		if defenderp := AsPlayer(defender); defenderp != nil {
			defenderp.PlaySound(sound)
		}

		return false
	})
}

//Killed kills this player, dropping all of its items where it stands.
func (p *Player) Killed(killer entity.MobileEntity) {
	p.SessionCache().SetVar("deathTime", time.Now())
	p.PlaySound("death")
	p.SendPacket(Death)
	for i := 0; i < 14; i++ {
		p.DeactivatePrayer(i)
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
		if p.Duel.Target != nil {
			p.Duel.Target.ResetDuel()
		}
		p.ResetDuel()
	}

	if killer != nil && killer.IsPlayer() {
		killerp := AsPlayer(killer)
		killerp.DistributeMeleeExp(p.ExperienceReward() / 4)
		killerp.Message("You have defeated " + p.Username() + "!")
	}
	for i, v := range deathItems {
		// becomes universally visible on NPCs, or temporarily private otherwise
		if i == 0 || p.Inventory.RemoveByID(v.ID, v.Amount) > -1 {
			if killer != nil && killer.IsPlayer() {
				v.Owner = AsPlayer(killer).Username()
			}
			AddItem(v)
		} else {
			log.Cheatf("Death item failed during removal: %v,%v owner:%v, killer:%v!\n", v.ID, v.Amount, p, killer)
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
	splat := NewHitsplat(p, amt)
	for _, player := range p.NearbyPlayers() {
		q, ok := player.Var("hitsplatQ")
		if !ok {
			player.SetVar("hitsplatQ", []HitSplat{splat})
			continue
		}
		player.SetVar("hitsplatQ", append(q.([]HitSplat), splat))
	}
	q, ok := p.Var("hitsplatQ")
	if !ok {
		p.SetVar("hitsplatQ", []HitSplat{splat})
		return
	}
	p.SetVar("hitsplatQ", append(q.([]HitSplat), splat))
}

//ItemBubble sends an item action bubble for this player to itself and any nearby players.
func (p *Player) ItemBubble(id int) {
	bubble := ItemBubble{p, id}
	for _, player := range p.NearbyPlayers() {
		q, ok := player.Var("bubbleQ")
		if !ok {
			player.SetVar("bubbleQ", []ItemBubble{bubble})
			continue
		}
		player.SetVar("bubbleQ", append(q.([]ItemBubble), bubble))
	}
	q, ok := p.Var("bubbleQ")
	if !ok {
		p.SetVar("bubbleQ", []ItemBubble{bubble})
		return
	}
	p.SetVar("bubbleQ", append(q.([]ItemBubble)))
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
	if p.IsFighting() || p.IsDueling() || p.State()&(StatePanelActive|StateFighting|StateDueling) != 0 {
		return
	}
	p.AddState(StateShopping)
	shop.Players.Add(p)
	p.SetVar("shop", shop)
	p.SendPacket(ShopOpen(shop))
}

//CloseBank closes the bank screen for this player and sets the appropriate state variables
func (p *Player) CloseShop() {
	if !p.HasState(StateShopping) {
		return
	}
	p.RemoveState(StateShopping)
	p.CurrentShop().Players.Remove(p)
	p.UnsetVar("shop")
	p.SendPacket(ShopClose)
}

//OpenBank opens a bank screen for the player and sets the appropriate state variables.
func (p *Player) OpenBank() {
	if p.IsFighting() || p.IsDueling() || p.State()&(StatePanelActive|StateFighting|StateDueling) != 0 {
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
	p.SendPacket(SystemUpdate(time.Until(UpdateTime).Milliseconds()))
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

//Read implements an io.Reader that detects what type of connection the underlying socket is using,
// and interprets the network byte stream accordingly.  Websockets require a lot of extra book-keeping
// to be used like this, and as such
func (p *Player) Read(data []byte) (n int, err error) {
	written := 0
	for written < len(data) {
		err := p.Socket.SetReadDeadline(time.Now().Add(time.Second * time.Duration(15)))
		if err != nil {
			return -1, errors.NewNetworkError("Deadline reached", true)
		}
		if p.IsWebsocket() && !p.hasReader {
			// reset buffer read index and create the next reader
			header, reader, err := wsutil.NextReader(p.Socket, ws.StateServerSide)
			p.hasReader = true
			p.SetVar("frameFin", header.Fin)
			if err != nil {
				if err == io.EOF && !header.Fin {
					return -1, errors.NewNetworkError("End of file mid-read:", true)
				} else if err == io.ErrUnexpectedEOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
					return -1, errors.NewNetworkError("closed conn", true)
				} else if e, ok := err.(stdnet.Error); ok && e.Timeout() {
					return -1, errors.NewNetworkError("timed out", true)
				}
				log.Warn("Problem creating reader for next websocket frame:", err)
			}
			if p.Reader == nil {
				p.Reader = bufio.NewReader(io.LimitReader(reader, header.Length))
			} else {
				p.Reader.Reset(io.LimitReader(reader, header.Length))
			}
		}
		n, err := p.Reader.Read(data[written:])
		if err != nil {
			if err == io.EOF && p.IsWebsocket() {
				p.hasReader = false
				p.Reader = nil
				if !p.VarBool("frameFin", false) {
					return -1, errors.NewNetworkError("closed conn", true)
				}
				continue
			} else if err == io.ErrUnexpectedEOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
				return -1, errors.NewNetworkError("closed conn", true)
			} else if e, ok := err.(stdnet.Error); ok && e.Timeout() {
				return -1, errors.NewNetworkError("timed out", true)
			}
			// continue
			return -1, errors.NewNetworkError(err.Error(), false)
		}
		written += n
	}
	return written, nil
}
