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

	"github.com/spkaeros/rscgo/pkg/config"
	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/engine/tasks"
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
		AppearanceLock   sync.RWMutex
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
		Tickables     scripts
		PostTickables scripts
		tickAction    statusReturnCall
		ActionLock    sync.RWMutex
		SendHandler   chan net.Packet
		ReplyMenuC    chan int8
		Equips        [12]int
		killer        sync.Once
		hasReader     bool
		SigKill       chan struct{}
		InQueue       chan net.Packet
		OutQueue      chan net.Packet
		Reader        *bufio.Reader
		DatabaseIndex int
		Mob
	}
)

func (p *Player) AsyncTick() bool {
	p.Tickables.async(interface{}(p))
	return true
}

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
	if target.IsNpc() {
		return target.(*NPC).Attackable()
	}
	if p.State()&StateFightingDuel == StateFightingDuel {
		return p.Duel.Target == target && p.DuelMagic()
	}
	p1 := AsPlayer(target)
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

//LocalAddress returns the local IP:port that this player connected to, or N/A if this player never connected somehow
func (p *Player) LocalAddress() string {
	if p.Socket == nil {
		return "N/A:N/A"
	}
	return p.Socket.LocalAddr().String()
}

func (p *Player) ConnectionPort() int {
	if p.Socket == nil {
		return -1
	}
	i, err := strconv.Atoi(strings.Split(p.Socket.LocalAddr().String(), ":")[1])
	if err != nil {
		return -1
	}
	return i
}

func (p *Player) IsWebsocket() bool {
	return p.ConnectionPort() == config.WSPort()
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
	p.ActionLock.Lock()
	p.tickAction = action
	p.ActionLock.Unlock()
}

func (p *Player) TickAction() func() bool {
	return p.tickAction
}

//ResetTickAction clears the distanced action, if any is queued.  Should be called any time the player is deliberately performing an action.
func (p *Player) ResetTickAction() {
	p.ActionLock.Lock()
	p.tickAction = nil
	p.ActionLock.Unlock()
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
func (p *Player) WalkingArrivalAction(target entity.MobileEntity, dist int, action func()) {
	p.tickAction = func() bool {
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
		//		if !p.WalkTo(NewLocation(target.X(), target.Y())) {
		///			return true
		//		}
		return !p.WalkTo(NewLocation(target.X(), target.Y()))
	}
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

	pathX := p.X()
	pathY := p.Y()
	for steps := 0; steps < 256; steps++ {
		if !p.ReachableCoords(pathX, pathY) {
			return false
		}
		// check deltas
		if pathX == target.X() && pathY == target.Y() {
			p.UnsetVar("triedReach")
			break
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

func (l Location) ReachableCoords(x, y int) bool {
	check := func(l, dst Location) bool {
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
			if IsTileBlocking(dst.X(), dst.Y(), masks[0]|masks[1], false) {
				return false
			}
		}
		return true
	}
	dst := l.Clone()
	start := dst.Clone()

	//	steps := 0
	for dst.X() != x || dst.Y() != y {
		if dst.X() > x {
			dst.x.Dec()
		} else if dst.X() < x {
			dst.x.Inc()
		}
		if dst.Y() > y {
			dst.y.Dec()
		} else if dst.Y() < y {
			dst.y.Inc()
		}
		//		if steps >= 8 {
		//			return false
		//		}
		//		steps++
		if !check(start.Clone(), dst.Clone()) {
			return false
		}
		start = dst.Clone()
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

//ResetAll in order, calls ResetFighting, ResetTrade, ResetTickAction, ResetFollowing, and CloseOptionMenu.
func (p *Player) ResetAll() {
	p.ResetFighting()
	p.ResetDuel()
	p.ResetTrade()
	p.ResetTickAction()
	p.ResetFollowing()
	p.CloseOptionMenu()
	p.CloseBank()
	p.CloseShop()
}

func (p *Player) ResetAllExceptDueling() {
	p.ResetTrade()
	p.ResetTickAction()
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
func (p *Player) NewPlayers() (players *MobList) {
	list := &MobList{}
	for _, r := range surroundingRegions(p.X(), p.Y()) {
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
	if p == nil || (packet.Opcode != 0 && !p.Connected()) {
		return
	}
	p.OutQueue <- *packet
}

type PlayerService interface {
	PlayerSave(*Player)
}

var DefaultPlayerService PlayerService

//Destroy sends a kill signal to the underlying client to tear down all of the I/O routines and save the player.
func (p *Player) Destroy() {
	p.killer.Do(func() {
		if p.Connected() {
			p.UpdateStatus(false)
			p.ResetAll()
		}
		p.OutQueue <- *Logout
		//go time.AfterFunc(time.Millisecond*time.Duration(640)*2, func() {
		tasks.TickList.Add(func() bool {
			if err := p.Socket.Close(); err != nil {
				log.Warn("Couldn't close socket:", err)
			}
			p.Inventory.Owner = nil
			p.Attributes.SetVar("lastIP", p.CurrentIP())
			close(p.SigKill)
			if player, ok := Players.FromIndex(p.Index); ok && player.UsernameHash() != p.UsernameHash() || !ok || !p.Connected() {
				if ok {
					log.Cheatf("Unauthenticated player being destroyed had index %d and there is a player that is assigned that index already! (%v)\n", p.Index, player)
				}
				return true
			}
			if p.Connected() {
				go DefaultPlayerService.PlayerSave(p)
			}
			RemovePlayer(p)
			p.SetConnected(false)
			log.Debug("Unregistered:{'" + p.Username() + "'@'" + p.CurrentIP() + "'}")
			return true
		})
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

func (l Location) WithinArea(area [2]Location) bool {
	return l.X() >= area[0].X() && l.X() <= area[1].X() && l.Y() >= area[0].Y() && l.Y() <= area[1].Y()
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
		return CollisionData(p.X()-1, p.Y()-1).CollisionMask&(ClipSouth|ClipWest) == 0
	}
	// Northwest target
	if p.X()+1 >= bounds[0].X() && p.X()+1 <= bounds[1].X() && p.Y()-1 >= bounds[0].Y() && p.Y()-1 <= bounds[1].Y() {
		return CollisionData(p.X()-1, p.Y()-1).CollisionMask&(ClipSouth|ClipEast) == 0
	}
	// Southeast target
	if p.X()-1 >= bounds[0].X() && p.X()-1 <= bounds[1].X() && p.Y()+1 >= bounds[0].Y() && p.Y()+1 <= bounds[1].Y() {
		return CollisionData(p.X()-1, p.Y()-1).CollisionMask&(ClipNorth|ClipWest) == 0
	}
	// Southwest target
	if p.X()+1 >= bounds[0].X() && p.X()+1 <= bounds[1].X() && p.Y()+1 >= bounds[0].Y() && p.Y()+1 <= bounds[1].Y() {
		return CollisionData(p.X()-1, p.Y()-1).CollisionMask&(ClipNorth|ClipEast) == 0
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

func (p *Player) SendFatigue() {
	p.SendPacket(Fatigue(p))
}

//Initialize informs the client of all of the various attributes of this player, and starts the stat normalization
// routine.
func (p *Player) Initialize() {
	AddPlayer(p)
	p.SetVar("initTime", time.Now())
	// update flags
	p.SetAppearanceChanged()
	p.SetSpriteUpdated()

	// settings panel
	p.SendPacket(ClientSettings(p))
	p.SendPacket(PrivacySettings(p))
	// social panel
	p.SendPacket(FriendList(p))
	p.SendPacket(IgnoreList(p))
	// TODO: Not canonical RSC, but definitely good QoL update...
	//  p.SendPacket(FightMode(p))

	// stat panel
	p.SendStats()
	p.SendFatigue()
	p.SendEquipBonuses()
	p.SendCombatPoints()

	// inventory panel
	p.SendInventory()
	// mesh related coordinate info and player index
	p.SendPlane()
	if !p.Attributes.Contains("madeAvatar") {
		p.OpenAppearanceChanger()
	} else {
		if !p.Reconnecting() {
			p.SendPacket(LoginBox(int(time.Since(p.Attributes.VarTime("lastLogin")).Hours()/24), p.Attributes.VarString("lastIP", "0.0.0.0")))
			p.SendPacket(WelcomeMessage)
		}
		p.Attributes.SetVar("lastLogin", time.Now())
	}
	for _, fn := range LoginTriggers {
		fn(p)
	}
}

func (p *Player) FlushOutgoing() {
	//	go func() {
	var writer net.WriteFlusher
	if p.IsWebsocket() {
		writer = wsutil.NewWriter(p.Socket, ws.StateServerSide, ws.OpBinary)
	} else {
		writer = bufio.NewWriter(p.Socket)
	}
	for {
		select {
		case packet, ok := <-p.OutQueue:
			if !ok {
				return
			}
			if packet.Opcode == 0 {
				writer.Write(packet.FrameBuffer)
				writer.Flush()
				continue
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
			writer.Write(append(header, packet.FrameBuffer[:frameLength]...))
			writer.Flush()
			continue
		default:
			return
		}
	}
	return
	//	}()
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
		Equips:           [12]int{entity.DefaultAppearance().Head, entity.DefaultAppearance().Body, entity.DefaultAppearance().Legs},
		SigKill:          make(chan struct{}),
		InQueue:          make(chan net.Packet, 50),
		OutQueue:         make(chan net.Packet, 50),
	}
	p.Reader = bufio.NewReader(socket)
	// TODO: Get rid of this self-referential member; figure out better way to handle client item updating
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

//Chat sends a player NPC chat message packet to the player and all other players around it.  If multiple msgs are
// provided, will sleep the goroutine for 3-4 ticks between each message, depending on length of message.
func (p *Player) Chat(msgs ...string) {
	for _, msg := range msgs {
		for _, player := range p.NearbyPlayers() {
			player.QueueQuestChat(p, nil, msg)
		}
		p.QueueQuestChat(p, nil, msg)

		sleepTicks := 3
		if len(msg) > 82 {
			sleepTicks++
		}
		time.Sleep(time.Millisecond * time.Duration(640*sleepTicks))
		//tasks.Schedule(sleepTicks, func() {

		//})
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
	reply := make(chan int8)
	p.ReplyMenuC = make(chan int8)
	p.Tickables.Add(func(p *Player) bool {
		select {
		case r, ok := <-p.ReplyMenuC:
			if !p.HasState(StateMenu) || !ok {
				return true
			}
			close(p.ReplyMenuC)
			p.RemoveState(StateMenu)
			if r < 0 || r > int8(len(options)-1) {
				log.Warn("Invalid option menu reply:", r)
				return true
			}

			if p.TargetNpc() != nil && p.HasState(StateChatting) {
				go func() {
					p.Chat(options[r])
					reply <- r
				}()
			}
			return true
		default:
			return false
		}
	})
	return int(<-reply)
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
	p.Skills().SetExp(idx, entity.LevelToExperience(lvl))
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

func (p *Player) PrayerActivated(idx int) bool {
	return p.VarBool("prayer"+strconv.Itoa(idx), false)
}

func (p *Player) PrayerOn(idx int) {
	if p.IsDueling() && !p.DuelPrayer() {
		p.Message("You cannot use prayer in this duel!")
		p.SendPrayers()
		return
	}
	boosterPrayers := [3][]int{
		{0, 3, 9},
		{1, 4, 10},
		{2, 5, 11},
	}
	defer p.SetVar("prayer"+strconv.Itoa(idx), true)
	for stat := 0; stat < 3; stat++ {
		for _, index := range boosterPrayers[stat] {
			if index == idx {
				for _, other := range boosterPrayers[stat] {
					if other != idx {
						p.PrayerOff(other)
					}
				}
				return
			}
		}
	}
	/*
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
	*/
}

func (p *Player) PrayerOff(idx int) {
	p.SetVar("prayer"+strconv.Itoa(idx), false)
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
		p.Attributes.UnsetVar("skullTime")
	}
	p.UpdateAppearance()
}

func (p *Player) ResetFighting() {
	defer p.Mob.ResetFighting()
	p.ResetDuel()
}

func (p *Player) Skulls() map[uint64]time.Time {
	records, ok := p.Var("skullRecord")
	if !ok || records == nil {
		p.SetVar("skullRecord", make(map[uint64]time.Time))
		records = p.VarChecked("skullRecord").(map[uint64]time.Time)
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
	if !p.SkulledOn(user) {
		// we skulled on them within 20 mins ago, ignore call
		return
	}
	p.SetSkulled(true)
	p.Skulls()[user] = time.Now()
}

func AsPlayer(m entity.MobileEntity) *Player {
	return m.(*Player)
}

func AsNpc(m entity.MobileEntity) *NPC {
	return m.(*NPC)
}

func (p *Player) StartCombat(target entity.MobileEntity) {
	if target.IsPlayer() {
		targetp := AsPlayer(target)
		targetp.PlaySound("underattack")
		if !p.IsDueling() && !targetp.SkulledOn(p.UsernameHash()) {
			p.SkullOn(targetp)
		}
	}
	p.SetVar("fightTarget", target)
	target.SessionCache().SetVar("fightTarget", p)
	target.SetRegionRemoved()
	p.Teleport(target.X(), target.Y())
	p.AddState(StateFighting)
	target.AddState(StateFighting)
	p.SetDirection(RightFighting)
	target.SetDirection(LeftFighting)
	curTick := 0
	attacker := entity.MobileEntity(p)
	defender := target
	tasks.TickList.Add(func() bool {
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

		attacker.SessionCache().Inc("fightRound", 1)
		if p.PrayerActivated(12) && attacker.IsNpc() {
			return false
		}
		nextHit := int(math.Min(float64(defender.Skills().Current(entity.StatHits)), float64(attacker.MeleeDamage(defender))))
		defender.Skills().DecreaseCur(entity.StatHits, nextHit)
		if defender.Skills().Current(entity.StatHits) <= 0 {
			if attacker.IsPlayer() {
				AsPlayer(attacker).PlaySound("victory")
			}
			defender.Killed(attacker)
			return true
		}
		if defender.IsNpc() && attacker.IsPlayer() {
			AsNpc(defender).CacheDamage(AsPlayer(attacker).UsernameHash(), nextHit)
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
	//tasks.TickList.Add(fightClosure)
}

//Killed kills this player, dropping all of its items where it stands.
func (p *Player) Killed(killer entity.MobileEntity) {
	p.SessionCache().SetVar("deathTime", time.Now())
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
		if p.Duel.Target != nil {
			p.Duel.Target.ResetDuel()
		}
		p.ResetDuel()
	}

	if killer != nil && killer.IsPlayer() {
		killer := killer.(*Player)
		killer.DistributeMeleeExp(p.ExperienceReward() / 4)
		killer.Message("You have defeated " + p.Username() + "!")
	}
	for i, v := range deathItems {
		// becomes universally visible on NPCs, or temporarily private otherwise
		if i == 0 || p.Inventory.RemoveByID(v.ID, v.Amount) > -1 {
			if killer != nil && killer.IsPlayer() {
				v.SetVar("belongsTo", killer.SessionCache().VarLong("username", 0))
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
	splat := NewHitsplat(p, amt)
	for _, player := range p.NearbyPlayers() {
		player.SetVar("hitsplatQ", append(player.VarChecked("hitsplatQ").([]HitSplat), splat))
	}
	p.SetVar("hitsplatQ", append(p.VarChecked("hitsplatQ").([]HitSplat), splat))
}

//ItemBubble sends an item action bubble for this player to itself and any nearby players.
func (p *Player) ItemBubble(id int) {
	//	for _, player := range p.NearbyPlayers() {
	//		player.SendPacket(PlayerItemBubble(p, id))
	//	}
	//	p.SendPacket(PlayerItemBubble(p, id))
	bubble := ItemBubble{p, id}
	for _, player := range p.NearbyPlayers() {
		player.SetVar("bubbleQ", append(player.VarChecked("bubbleQ").([]ItemBubble), bubble))
	}
	p.SetVar("bubbleQ", append(p.VarChecked("bubbleQ").([]ItemBubble)))
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

/*
//Read implements an io.Reader that detects what type of connection the underlying socket is using,
// and interprets the network byte stream accordingly.  Websockets require a lot of extra book-keeping
// to be used like this, and as such
func (p *Player) Read(data []byte) (n int, err error) {
	written := 0
	for written < len(data) {
		err := p.Socket.SetDeadline(time.Now().Add(time.Second * time.Duration(15)))
		if err != nil {
			p.Destroy()
			return -1, errors.NewNetworkError("Deadline reached", true)
		}
		if p.IsWebsocket() && !p.hasReader {
			// reset buffer read index and create the next reader
			header, reader, err := wsutil.NextReader(p.Socket, ws.StateServerSide)
			p.hasReader = true
			p.SetVar("frameFinished", header.Fin)
			if err != nil {
				if err == io.EOF {
					p.Destroy()
					return -1, errors.NewNetworkError("EOF encountered during wsutil.NextReader:" + err.Error(), true)
				} else if err == io.ErrUnexpectedEOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
					p.Destroy()
					return -1, errors.NewNetworkError("closed conn", true)
				} else if e, ok := err.(stdnet.Error); ok && e.Timeout() {
					p.Destroy()
					return -1, errors.NewNetworkError("timed out", true)
				}
				log.Warn("Problem creating reader for next websocket frame:", err)
			}
//			if p.Reader != nil {
				p.Reader.Reset(reader, header.)
//			} else {
//				p.Reader = bufio.NewReader(io.LimitReader(reader, header.Length))
//			}
		}
		n, err := p.Reader.Read(data[written:])
//		p.hasReader = err == io.EOF
		if err != nil {
			if err == io.EOF {
				if !p.VarBool("frameFinished", false) {
					p.Destroy()
					log.Debug("EOF on unfinished frame")
//					return -1, errors.NewNetworkError("EOF on unfinished frame", true)
				}
				p.hasReader=false
				log.Debug("End of frame detected")
				continue
			}
			if err == io.ErrUnexpectedEOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
				p.Destroy()
				return -1, errors.NewNetworkError("closed conn", true)
			}
			if e, ok := err.(stdnet.Error); ok && e.Timeout() {
				p.Destroy()
				return -1, errors.NewNetworkError("timed out", true)
			}
//			if p.IsWebsocket() && err == io.EOF {
//				p.hasReader=false
//				log.Warn("Problem reading from socket:", err.Error())
//				p.Reader = nil
//				continue
//			}
			return -1, errors.NewNetworkError("Unknown:" + err.Error(), false)
			//continue
		}
		written += n
	}
	return written, nil
}
*/

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
			if err == io.EOF {
				p.hasReader = false
				if !p.VarBool("frameFin", false) {
					return -1, errors.NewNetworkError("closed conn", true)
				}
				continue
			} else if err == io.ErrUnexpectedEOF || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "use of closed") {
				return -1, errors.NewNetworkError("closed conn", true)
			} else if e, ok := err.(stdnet.Error); ok && e.Timeout() {
				return -1, errors.NewNetworkError("timed out", true)
			}
			//	continue
			return -1, errors.NewNetworkError("idklol", false)
		}
		written += n
	}
	return written, nil
}
