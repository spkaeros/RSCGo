package world

import (
	"github.com/spkaeros/rscgo/pkg/rand"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"sync"
	"time"
)

//MSIdle The default MobState, means doing nothing.
const (
	//MSIdle The default MobState, means doing nothing.
	MSIdle = iota
	//MSBanking The mob is banking.
	MSBanking
	//MSChatting The mob is chatting with a NPC
	MSChatting
	//MSOptionMenu The mob is in an option menu.  The option menu handling routines will remove this state as soon
	// as they end, so if this is activated, there is an option menu waiting for a reply.
	MSOptionMenu
	//MSTrading The mob is negotiating a trade.
	MSTrading
	//MSDueling The mob is negotiating a duel.
	MSDueling
	//MSFighting The mob is fighting.
	MSFighting
	//MSBatching The mob is performing a skill that repeats itself an arbitrary number of times.
	MSBatching
	//MSSleeping The mob is using a bed or sleeping bag, and trying to solve a CAPTCHA
	MSSleeping
	//MSBusy Generic busy state
	MSBusy
	//MSChangingAppearance Indicates that the mob in this state is in the player aooearance changing screen
	MSChangingAppearance
)

const (
	SyncSprite     = 1
	SyncMoved      = 1 << 1
	SyncRemoved    = 1 << 2
	SyncAppearance = 1 << 3

	SyncNeedsPosition = SyncRemoved | SyncMoved | SyncSprite
)

//Mob Represents a mobile entity within the game world.
type Mob struct {
	*Entity
	TransAttrs *AttributeList
	SyncMask   int
	sync.RWMutex
}

type MobileEntity interface {
	X() int
	Y() int
	Skills() *SkillTable
	MeleeDamage(target MobileEntity) int
	Defense(float32, float32) float32
	Transients() *AttributeList
	IsFighting() bool
	FightTarget() MobileEntity
	SetFightTarget(MobileEntity)
	FightRound() int
	SetFightRound(int)
	ResetFighting()
	HasState(...int) bool
	AddState(int)
	RemoveState(int)
	State() int
	Busy() bool
	SetX(int)
	SetY(int)
	SetCoords(int, int, bool)
	Teleport(int, int)
	Direction() int
	SetDirection(int)
	Path() *Pathway
	ResetPath()
	SetPath(*Pathway)
	TraversePath()
	FinishedPath() bool
	SetLocation(Location, bool)
	UpdateLastFight()
	LastFight() time.Time
	UpdateLastRetreat()
	LastRetreat() time.Time
	SetRegionMoved()
	SetRegionRemoved()
	SetSpriteUpdated()
	SetAppearanceChanged()
	ResetRegionMoved()
	ResetRegionRemoved()
	ResetSpriteUpdated()
	ResetAppearanceChanged()
}

func (m *Mob) Transients() *AttributeList {
	return m.TransAttrs
}

//Busy Returns true if this mobs state is anything other than idle. otherwise returns false.
func (m *Mob) Busy() bool {
	return m.State() != MSIdle
}

func (m *Mob) IsFighting() bool {
	return m.HasState(MSFighting)
}

func (m *Mob) FightTarget() MobileEntity {
	return m.TransAttrs.VarMob("fightTarget")
}

func (m *Mob) SetFightTarget(m2 MobileEntity) {
	m.TransAttrs.SetVar("fightTarget", m2)
}

func (m *Mob) FightRound() int {
	return m.TransAttrs.VarInt("fightRound", 0)
}

func (m *Mob) SetFightRound(i int) {
	m.TransAttrs.SetVar("fightRound", i)
}

func (m *Mob) LastRetreat() time.Time {
	return m.TransAttrs.VarTime("lastRetreat")
}

func (m *Mob) LastFight() time.Time {
	return m.TransAttrs.VarTime("lastFight")
}

func (m *Mob) UpdateLastRetreat() {
	m.TransAttrs.SetVar("lastRetreat", time.Now())
}

func (m *Mob) UpdateLastFight() {
	m.TransAttrs.SetVar("lastFight", time.Now())
}

//Direction Returns the mobs direction.
func (m *Mob) Direction() int {
	return m.TransAttrs.VarInt("direction", North)
}

//SetDirection Sets the mobs direction.
func (m *Mob) SetDirection(direction int) {
	m.TransAttrs.SetVar("direction", direction)
	m.SetSpriteUpdated()
}

//SetSpriteUpdated Sets the synchronization flag for whether this mob changed directions to true.
func (m *Mob) SetSpriteUpdated() {
	m.Lock()
	defer m.Unlock()
	m.SyncMask |= SyncSprite
}

//SetRegionRemoved Sets the synchronization flag for whether this mob needs to be removed to true.
func (m *Mob) SetRegionRemoved() {
	m.Lock()
	defer m.Unlock()
	m.SyncMask |= SyncRemoved
}

//UpdateSelf Sets the synchronization flag for whether this mob has moved to true.
func (m *Mob) SetRegionMoved() {
	m.Lock()
	defer m.Unlock()
	m.SyncMask |= SyncMoved
}

func (m *Mob) SetAppearanceChanged() {
	m.Lock()
	defer m.Unlock()
	m.SyncMask |= SyncAppearance
}

func (m *Mob) ResetRegionMoved() {
	m.Lock()
	defer m.Unlock()
	m.SyncMask &= ^SyncMoved
}

func (m *Mob) ResetRegionRemoved() {
	m.Lock()
	defer m.Unlock()
	m.SyncMask &= ^SyncRemoved
}

func (m *Mob) ResetAppearanceChanged() {
	m.Lock()
	defer m.Unlock()
	m.SyncMask &= ^SyncAppearance
}

func (m *Mob) ResetSpriteUpdated() {
	m.Lock()
	defer m.Unlock()
	m.SyncMask &= ^SyncSprite
}

//SetPath Sets the mob's current pathway to path.  If path is nil, effectively resets the mobs path.
func (m *Mob) SetPath(path *Pathway) {
	m.TransAttrs.SetVar("path", path)
}

func (m *Mob) WalkTo(end Location) {
	path := MakePath(m.Location, end)
	m.SetPath(path)
}

//Path returns the path that this mob is trying to traverse.
func (m *Mob) Path() *Pathway {
	return m.TransAttrs.VarPath("path")
}

//ResetPath Sets the mobs path to nil, to stop the traversal of the path instantly
func (m *Mob) ResetPath() {
	m.TransAttrs.UnsetVar("path")
}

//FinishedPath Returns true if the mobs path is nil, the paths current waypoint exceeds the number of waypoints available, or the next tile in the path is not a valid location, implying that we have reached our destination.
func (m *Mob) FinishedPath() bool {
	path := m.Path()
	if path == nil {
		return true
	}
	return path.CurrentWaypoint >= path.countWaypoints() || !path.nextTileFrom(m.Location).IsValid()
}

//SetLocation Sets the mobs location.
func (m *Mob) SetLocation(location Location, teleport bool) {
	m.SetCoords(location.X(), location.Y(), teleport)
}

func (p *Player) SetLocation(l Location, teleport bool) {
	p.UpdateRegion(l.X(), l.Y())
	p.Mob.SetLocation(l, teleport)
}

func (n *NPC) SetLocation(l Location, teleport bool) {
	n.UpdateRegion(l.X(), l.Y())
	n.Mob.SetLocation(l, teleport)
}

//SetCoords Sets the mobs locations coordinates.
func (m *Mob) SetCoords(x, y int, teleport bool) {
	if !teleport {
		m.SetDirection(m.DirectionTo(x, y))
		m.SetRegionMoved()
	} else {
		m.SetRegionRemoved()
	}
	m.SetX(x)
	m.SetY(y)
}

func (p *Player) SetCoords(x, y int, teleport bool) {
	p.UpdateRegion(x, y)
	p.Mob.SetCoords(x, y, teleport)
}

func (n *NPC) SetCoords(x, y int, teleport bool) {
	n.UpdateRegion(x, y)
	n.Mob.SetCoords(x, y, teleport)
}

func (p *Player) Teleport(x, y int) {
	p.SetCoords(x, y, true)
}

func (n *NPC) Teleport(x, y int) {
	n.SetCoords(x, y, true)
}

func (m *Mob) State() int {
	return m.TransAttrs.VarInt("state", MSIdle)
}

//HasState Returns true if the mob has any of these states
func (m *Mob) HasState(state ...int) bool {
	return m.Transients().HasMasks("state", state...)
}

func (m *Mob) AddState(state int) {
	if state == MSIdle {
		m.Transients().SetVar("state", MSIdle)
		return
	}
	if m.HasState(state) {
		log.Warning.Println("Attempted to add a Mobstate that we already have:", state)
		return
	}
	m.Transients().StoreMask("state", state)
}

func (m *Mob) RemoveState(state int) {
	if state == MSIdle {
		return
	}
	if !m.HasState(state) {
		log.Warning.Println("Attempted to remove a Mobstate that we did not add:", state)
		return
	}
	m.Transients().RemoveMask("state", state)
}

//ResetFighting Resets melee fight related variables
func (m *Mob) ResetFighting() {
	target := m.TransAttrs.VarMob("fightTarget")
	if target != nil && target.IsFighting() {
		target.UpdateLastFight()
		target.Transients().UnsetVar("fightTarget")
		target.Transients().UnsetVar("fightRound")
		target.SetDirection(North)
		target.RemoveState(MSFighting)
		target.UpdateLastFight()
	}
	if m.IsFighting() {
		target.UpdateLastFight()
		m.TransAttrs.UnsetVar("fightTarget")
		m.TransAttrs.UnsetVar("fightRound")
		m.SetDirection(North)
		m.RemoveState(MSFighting)
		m.UpdateLastFight()
	}
}

//FightMode Returns the players current fight mode.
func (m *Mob) FightMode() int {
	return m.TransAttrs.VarInt("fight_mode", 0)
}

//SetFightMode Sets the players fightmode to i.  0=all,1=attack,2=defense,3=strength
func (m *Mob) SetFightMode(i int) {
	m.TransAttrs.SetVar("fight_mode", i)
}

//ArmourPoints Returns the players armour points.
func (m *Mob) ArmourPoints() int {
	return m.TransAttrs.VarInt("armour_points", 1)
}

//SetArmourPoints Sets the players armour points to i.
func (m *Mob) SetArmourPoints(i int) {
	m.TransAttrs.SetVar("armour_points", i)
}

//PowerPoints Returns the players power points.
func (m *Mob) PowerPoints() int {
	return m.TransAttrs.VarInt("power_points", 1)
}

//SetPowerPoints Sets the players power points to i
func (m *Mob) SetPowerPoints(i int) {
	m.TransAttrs.SetVar("power_points", i)
}

//AimPoints Returns the players aim points
func (m *Mob) AimPoints() int {
	return m.TransAttrs.VarInt("aim_points", 1)
}

//SetAimPoints Sets the players aim points to i.
func (m *Mob) SetAimPoints(i int) {
	m.TransAttrs.SetVar("aim_points", i)
}

//MagicPoints Returns the players magic points
func (m *Mob) MagicPoints() int {
	return m.TransAttrs.VarInt("magic_points", 1)
}

//SetMagicPoints Sets the players magic points to i
func (m *Mob) SetMagicPoints(i int) {
	m.TransAttrs.SetVar("magic_points", i)
}

//PrayerPoints Returns the players prayer points
func (m *Mob) PrayerPoints() int {
	return m.TransAttrs.VarInt("prayer_points", 1)
}

//SetPrayerPoints Sets the players prayer points to i
func (m *Mob) SetPrayerPoints(i int) {
	m.TransAttrs.SetVar("prayer_points", i)
}

//RangedPoints Returns the players ranged points.
func (m *Mob) RangedPoints() int {
	return m.TransAttrs.VarInt("ranged_points", 1)
}

//SetRangedPoints Sets the players ranged points tp i.
func (m *Mob) SetRangedPoints(i int) {
	m.TransAttrs.SetVar("ranged_points", i)
}

func (m *Mob) Skills() *SkillTable {
	return m.TransAttrs.VarSkills("skills")
}

func PrayerModifiers(m MobileEntity) [3]float32 {
	modifiers := [3]float32{1.0, 1.0, 1.0}

	if m, ok := m.(*Player); ok {
		if m.PrayerActivated(0) {
			modifiers[1] += .05
		}
		if m.PrayerActivated(1) {
			modifiers[2] += .05
		}
		if m.PrayerActivated(2) {
			modifiers[0] += .05
		}
		if m.PrayerActivated(3) {
			modifiers[1] += .1
		}
		if m.PrayerActivated(4) {
			modifiers[2] += .1
		}
		if m.PrayerActivated(5) {
			modifiers[0] += .1
		}
		if m.PrayerActivated(9) {
			modifiers[1] += .15
		}
		if m.PrayerActivated(10) {
			modifiers[2] += .15
		}
		if m.PrayerActivated(11) {
			modifiers[0] += .15
		}
	}
	return modifiers
}

func (m *Mob) StyleBonus(stat int) int {
	mode := m.FightMode()
	if mode == 0 {
		return 1
	} else if (mode == 2 && stat == 0) || (mode == 1 && stat == 2) || (mode == 3 && stat == 1) {
		return 3
	}
	return 0
}

//MaxHit Calculates and returns the current max hit for this mob.
func (m *Mob) MaxHit(prayer float32) int {
	newStr := (float32(m.Skills().Current(StatStrength)) * prayer) + float32(m.StyleBonus(StatStrength))
	return int((newStr*((float32(m.PowerPoints())*0.00175)+0.1) + 1.05) * 0.95)
}

func (m *Mob) Accuracy(npcMul float32, prayer float32) float32 {
	styleBonus := float32(m.StyleBonus(StatAttack))
	attackLvl := (float32(m.Skills().Current(StatAttack)) * prayer) + styleBonus + 8
	multiplier := float32(m.AimPoints() + 64)
	multiplier *= npcMul
	return attackLvl * multiplier
}

func (m *Mob) Defense(npcMul float32, prayer float32) float32 {
	styleBonus := float32(m.StyleBonus(StatDefense))
	defenseLvl := (float32(m.Skills().Current(StatDefense)) * prayer) + styleBonus + 8
	multiplier := float32(m.ArmourPoints() + 64)
	multiplier *= npcMul
	return defenseLvl * multiplier
}

func (n *NPC) MeleeDamage(target MobileEntity) int {
	att := n.Accuracy(0.9, PrayerModifiers(n)[0])
	mul := float32(1.0)
	if _, ok := target.(*NPC); ok {
		mul = 0.9
	}
	def := target.Defense(mul, PrayerModifiers(target)[1])
	max := n.MaxHit(PrayerModifiers(n)[2])
	if att*10 < def {
		return 0
	}

	finalAtt := int((att / (2.0 * (def + 1.0))) * 10000.0)

	if att > def {
		finalAtt = int((1.0 - ((def + 2.0) / (2.0 * (att + 1.0)))) * 10000.0)
	}

	roll := rand.Int31N(0, 10000)
	//	log.Info.Println(finalAtt, roll, att, def, max)
	if finalAtt > roll {
		return rand.Int31N(0, max)
	}
	return 0
}

func (p *Player) MeleeDamage(target MobileEntity) int {
	att := p.Accuracy(1.0, PrayerModifiers(p)[0])
	mul := float32(1.0)
	if _, ok := target.(*NPC); ok {
		mul = 0.9
	}
	def := target.Defense(mul, PrayerModifiers(target)[1])
	max := p.MaxHit(PrayerModifiers(p)[2])
	if att*10 < def {
		return 0
	}

	finalAtt := int((att / (2.0 * (def + 1.0))) * 10000.0)

	if att > def {
		finalAtt = int((1.0 - ((def + 2.0) / (2.0 * (att + 1.0)))) * 10000.0)
	}

	roll := rand.Int31N(0, 10000)
	//	log.Info.Println(finalAtt, roll, att, def, max)
	if finalAtt > roll {
		return rand.Int31N(0, max)
	}
	return 0
}
