package world

import (
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/spkaeros/rscgo/pkg/game/entity"
	"github.com/spkaeros/rscgo/pkg/isaac"
	"github.com/spkaeros/rscgo/pkg/log"
	rscRand "github.com/spkaeros/rscgo/pkg/rand"
)

type MobState = int

//StateIdle The default MobState, means doing nothing.
const (
	//StateIdle The default MobState, means doing nothing.
	StateIdle MobState = 0
	//StateChatting The mob is chatting with another mob.
	StateChatting = 1 << iota
	//StateFighting The mob is fighting.
	StateFighting
	//StateBanking The mob is banking.
	StateBanking
	//StateMenu The mob is in an option menu.  The option menu handling routines will remove this state as soon
	// as they end, so if this is activated, there is an option menu waiting for a reply.
	StateMenu
	//StateTrading The mob is negotiating a trade.
	StateTrading
	//StateDueling The mob is negotiating a duel.
	StateDueling
	//MSBatching The mob is performing a skill that repeats itself an arbitrary number of times.
	MSBatching
	//StateSleeping The mob is using a bed or sleeping bag, and trying to solve a CAPTCHA
	StateSleeping
	//StateChangingLooks Indicates that the mob in this state is in the player aooearance changing screen
	StateChangingLooks
	//StateShopping Indicates that the mob in this state is using a shop interface
	StateShopping
	//MSItemAction Indicates that the mob in this state is doing an inventory action
	MSItemAction

	StateFightingDuel   = StateDueling | StateFighting
	StateChatChoosing   = StateMenu | StateChatting
	StateItemChoosing   = StateMenu | MSItemAction
	StateObjectChoosing = StateMenu | MSBatching

	StatePanelActive = StateBanking | StateShopping | StateChangingLooks | StateSleeping | StateTrading | StateDueling

	StateBusy      = StatePanelActive | StateChatting | MSItemAction | MSBatching
	StateWaitEvent = StateMenu | StateChatting | MSItemAction | MSBatching
)

const (
	SyncInit = 0
	SyncSprite = 1 << iota
	SyncMoved
	SyncRemoved
	SyncAppearance

	SyncNeedsPosition = SyncRemoved | SyncMoved | SyncSprite
)

// mobSet a collection of entity.MobileEntitys
type mobSet []entity.MobileEntity

//MobList a container type for holding entity.MobileEntitys
type MobList struct {
	mobSet
	sync.RWMutex
}

//NewMobList returns a pointer to a newly pre-allocated MobList, with an initial capacity
// of 255.
func NewMobList() *MobList {
	return &MobList{mobSet: make(mobSet, 0, 255)}
}

//Add Adds a entity.MobileEntity to this MobList
func (l *MobList) Add(m entity.MobileEntity) int {
	l.Lock()
	defer l.Unlock()
	l.mobSet = append(l.mobSet, m)
	return len(l.mobSet)
}

//Range Runs action(entity.MobileEntity) for each entity.MobileEntity in the lists collection, until
// either running out of entries, or action returns true.
func (l *MobList) Range(action func(entity.MobileEntity) bool) int {
	l.RLock()
	defer l.RUnlock()
	for i, v := range l.mobSet {
		if action(v) {
			return i
		}
	}
	return -1
}

//Size returns the number of mobile entitys entered into this list.
func (l *MobList) Size() int {
	l.RLock()
	defer l.RUnlock()
	return len(l.mobSet)
}

//Size returns the number of mobile entitys entered into this list.
func (l *MobList) Contains(mob entity.MobileEntity) bool {
	return l.Range(func(m entity.MobileEntity) bool {
		return m == mob
	}) > -1
}

//Remove removes a entity.MobileEntity from this list and reslices the collection.
func (l *MobList) Remove(m entity.MobileEntity) bool {
	l.Lock()
	defer l.Unlock()
	for i, v := range l.mobSet {
		if v == m {
			if len(l.mobSet) >= i+1 {
				l.mobSet = append(l.mobSet[:i], l.mobSet[i+1:]...)
			} else {
				l.mobSet = l.mobSet[:i]
			}
			return true
		}
	}
	return false
}

func (l *MobList) RangePlayers(action func(*Player) bool) int {
	return l.Range(func(m entity.MobileEntity) bool {
		if p, ok := m.(*Player); ok {
			return action(p)
		}
		return false
	})
}

func (l *MobList) RangeNpcs(action func(*NPC) bool) int {
	return l.Range(func(m entity.MobileEntity) bool {
		if n, ok := m.(*NPC); ok {
			return action(n)
		}
		return false
	})
}

func (l *MobList) Get(idx int) entity.MobileEntity {
	l.RLock()
	defer l.RUnlock()
	for _, v := range l.mobSet {
		if v.ServerIndex() == idx {
			return v
		}
	}
	return nil
}

//Mob Represents a mobile entity within the game world.
type Mob struct {
	SyncMask       int
	skills		   entity.SkillTable
	path		   *Pathway
	Prayers		   [15]bool
	Entity
	*entity.AttributeList
	sync.RWMutex
}

//ExperienceReward returns the total rewarded experience upon killing a mob.
func (m *Mob) ExperienceReward() int {
	meleeTotal := (m.Skills().Current(entity.StatStrength) + m.Skills().Current(entity.StatAttack) + m.Skills().Current(entity.StatDefense)) * 2
	return (int(math.Floor(float64(m.Skills().Current(entity.StatHits)+meleeTotal)/7)) * 2) + 20
}

func (m *Mob) TargetMob() entity.MobileEntity {
	return m.VarMob("targetMob")
}

func (m *Mob) TargetNpc() *NPC {
	if n, ok := m.TargetMob().(*NPC); ok {
		return n
	}
	return nil
}

func (m *Mob) TargetPlayer() *Player {
	if p, ok := m.TargetMob().(*Player); ok {
		return p
	}
	return nil
}

func (p *Player) IsPlayer() bool {
	return true
}

func (p *Player) Type() entity.Type {
	return entity.TypePlayer
}

func (p *Player) ServerIndex() int {
	if p == nil {
		return -1
	}
	return p.Index
}

func (n *NPC) ServerIndex() int {
	return n.Index
}

func (p *Player) IsNpc() bool {
	return false
}

func (n *NPC) Type() entity.Type {
	return entity.TypeNpc
}

func (n *NPC) IsPlayer() bool {
	return false
}

func (n *NPC) IsNpc() bool {
	return true
}

func (m *Mob) Skulls() map[uint64]time.Time {
	return nil
}

//Busy Returns true if this mobs state is anything other than idle. otherwise returns false.
func (m *Mob) Busy() bool {
	return m.State()&StateBusy != 0
}

func (m *Mob) BusyInput() bool {
	return m.State()&StateWaitEvent == StateChatChoosing || m.State()&StateWaitEvent == StateItemChoosing ||
		m.State()&StateWaitEvent == StateObjectChoosing || m.State() == StateIdle
	//	return m.State() != StateIdle && m.State() != MSItemAction
}

func (m *Mob) IsFighting() bool {
	return m.HasState(StateFighting)
}

func (m *Mob) FightTarget() entity.MobileEntity {
	return m.VarMob("fightTarget")
}

func (m *Mob) SetFightTarget(m2 entity.MobileEntity) {
	m.SetVar("fightTarget", m2)
}

func (m *Mob) FightRound() int {
	return m.VarInt("fightRound", 0)
}

func (m *Mob) SetFightRound(i int) {
	m.SetVar("fightRound", i)
}

func (m *Mob) LastRetreat() time.Time {
	return m.VarTime("lastRetreat")
}

func (m *Mob) LastFight() time.Time {
	return m.VarTime("lastFight")
}

func (m *Mob) UpdateLastRetreat() {
	m.SetVar("lastRetreat", time.Now())
}

func (m *Mob) UpdateLastFight() {
	m.SetVar("lastFight", time.Now())
}

//Direction Returns the mobs direction.
func (m *Mob) Direction() int {
	return m.VarInt("direction", North)
}

//SetDirection Sets the mobs direction.
func (m *Mob) SetDirection(direction int) {
	m.SetVar("direction", direction)
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
	m.path = path
	
}

func (m *Mob) WalkTo(end Location) bool {
	path, ok := MakePath(m.Location, end)
	m.SetPath(path)
	return ok
}

//Path returns the path that this mob is trying to traverse.
func (m *Mob) Path() *Pathway {
	return m.path
}

//ResetPath Sets the mobs path to nil, to stop the traversal of the path instantly
func (m *Mob) ResetPath() {
	m.UnsetVar("pathLength")
	m.path = nil
}

//FinishedPath Returns true if the mobs path is nil, the paths current waypoint exceeds the number of waypoints available, or the next tile in the path is not a valid location, implying that we have reached our destination.
func (m *Mob) FinishedPath() bool {
	path := m.Path()
	if path == nil {
		return m.VarInt("pathLength", 0) <= 0
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

func (m *Mob) SessionCache() *entity.AttributeList {
	return m.AttributeList
}

func (m *Mob) State() int {
	return m.VarInt("state", StateIdle)
}

//HasState Returns true if the mob has any of these states
func (m *Mob) HasState(state ...int) bool {
	return m.HasMasks("state", state...)
}

func (m *Mob) AddState(state int) {
	if state == StateIdle {
		m.SetVar("state", StateIdle)
		return
	}
	if m.HasState(state) {
		log.Warn("Attempted to add a Mobstate that we already have:", state)
		return
	}
	m.StoreMask("state", state)
}

func (m *Mob) RemoveState(state int) {
	if state == StateIdle {
		return
	}
	if !m.HasState(state) {
		log.Warn("Attempted to remove a Mobstate that we did not add:", state)
		return
	}
	m.RemoveMask("state", state)
}

//ResetFighting Resets melee fight related variables
func (m *Mob) ResetFighting() {
	target := m.VarMob("fightTarget")
	if m.IsFighting() {
		m.UnsetVar("fightTarget")
		m.UnsetVar("fightRound")
		m.SetDirection(NorthWest)
		m.RemoveState(StateFighting)
		if m.HasState(StateDueling) {
			m.RemoveState(StateDueling)
		}
		m.UpdateLastFight()
	}
	if target != nil && target.IsFighting() {
		target.SessionCache().UnsetVar("fightTarget")
		target.SessionCache().UnsetVar("fightRound")
		target.SetDirection(NorthWest)
		target.RemoveState(StateFighting)
		if target.HasState(StateDueling) {
			target.RemoveState(StateDueling)
		}
		target.UpdateLastFight()
	}
}

//FightMode Returns the players current fight mode.
func (m *Mob) FightMode() int {
	return m.VarInt("fight_mode", 0)
}

//SetFightMode Sets the players fightmode to i.  0=all,1=attack,2=defense,3=strength
func (m *Mob) SetFightMode(i int) {
	m.SetVar("fight_mode", i)
}

//ArmourPoints Returns the players armour points.
func (m *Mob) ArmourPoints() int {
	return m.VarInt("armour_points", 1)
}

//SetArmourPoints Sets the players armour points to i.
func (m *Mob) SetArmourPoints(i int) {
	m.SetVar("armour_points", i)
}

func (m *Mob) IncArmourPoints(i int) {
	points := m.VarInt("armour_points", 1)
	if points == 1 {
		m.SetVar("armour_points", i)
	} else {
		m.Inc("armour_points", i)
	}
}

//PowerPoints Returns the players power points.
func (m *Mob) PowerPoints() int {
	return m.VarInt("power_points", 1)
}

//SetPowerPoints Sets the players power points to i
func (m *Mob) SetPowerPoints(i int) {
	m.SetVar("power_points", i)
}

func (m *Mob) IncPowerPoints(i int) {
	points := m.VarInt("power_points", 1)
	if points == 1 {
		m.SetVar("power_points", i)
	} else {
		m.Inc("power_points", i)
	}
}

//AimPoints Returns the players aim points
func (m *Mob) AimPoints() int {
	return m.VarInt("aim_points", 1)
}

func (m *Mob) IncAimPoints(i int) {
	points := m.VarInt("aim_points", 1)
	if points == 1 {
		m.SetVar("aim_points", i)
	} else {
		m.Inc("aim_points", i)
	}
}

//SetAimPoints Sets the players aim points to i.
func (m *Mob) SetAimPoints(i int) {
	m.SetVar("aim_points", i)
}

//MagicPoints Returns the players magic points
func (m *Mob) MagicPoints() int {
	return m.VarInt("magic_points", 1)
}

func (m *Mob) IncMagicPoints(i int) {
	points := m.VarInt("magic_points", 1)
	if points == 1 {
		m.SetVar("magic_points", i)
	} else {
		m.Inc("magic_points", i)
	}
}

//SetMagicPoints Sets the players magic points to i
func (m *Mob) SetMagicPoints(i int) {
	m.SetVar("magic_points", i)
}

func (m *Mob) IncPrayerPoints(i int) {
	points := m.VarInt("prayer_points", 1)
	if points == 1 {
		m.SetVar("prayer_points", i)
	} else {
		m.Inc("prayer_points", i)
	}
}

//PrayerPoints Returns the players prayer points
func (m *Mob) PrayerPoints() int {
	return m.VarInt("prayer_points", 1)
}

//SetPrayerPoints Sets the players prayer points to i
func (m *Mob) SetPrayerPoints(i int) {
	m.SetVar("prayer_points", i)
}

//RangedPoints Returns the players ranged points.
func (m *Mob) RangedPoints() int {
	return m.VarInt("ranged_points", 1)
}

func (m *Mob) IncRangedPoints(i int) {
	points := m.VarInt("ranged_points", 1)
	if points == 1 {
		m.SetVar("ranged_points", i)
	} else {
		m.Inc("ranged_points", i)
	}
}

//SetRangedPoints Sets the players ranged points tp i.
func (m *Mob) SetRangedPoints(i int) {
	m.SetVar("ranged_points", i)
}

func (m *Mob) Skills() *entity.SkillTable {
	return &m.skills
}

func (m *Mob) PrayerModifiers() [3]float64 {
	var modifiers = [...]float64{1.0, 1.0, 1.0}

	mods := []float64{1.05, 1.10, 1.15}

	defenseMods := []int{0, 3, 9}
	strengthMods := []int{1, 4, 10}
	attackMods := []int{2, 5, 11}
	for idx, prayer := range defenseMods {
		if m.VarBool("prayer"+strconv.Itoa(prayer), false) {
			modifiers[entity.StatDefense] = mods[idx]
		}
	}

	for idx, prayer := range strengthMods {
		if m.VarBool("prayer"+strconv.Itoa(prayer), false) {
			modifiers[entity.StatStrength] = mods[idx]
		}
	}

	for idx, prayer := range attackMods {
		if m.VarBool("prayer"+strconv.Itoa(prayer), false) {
			modifiers[entity.StatAttack] = mods[idx]
		}
	}
	return modifiers
}

func (m *Mob) StyleBonus(stat int) int {
	mode := m.FightMode()
	if mode == 0 {
		return 1
	}
	if mode == (stat+1)%3+1 {
		return 3
	}
	return 0
}

//MaxMeleeDamage Calculates and returns the current max hit for this mob, based on many variables.
func (m *Mob) MaxMeleeDamage() float64 {
	return ((float64(m.Skills().Current(entity.StatStrength))*m.PrayerModifiers()[entity.StatStrength])+float64(m.StyleBonus(entity.StatStrength)))*((float64(m.PowerPoints())*0.00175)+0.1) + 1.05
}

//AttackPoints Calculates and returns the accuracy capability of this mob, based on many variables, as a single variable.
func (m *Mob) AttackPoints() float64 {
	return ((float64(m.Skills().Current(entity.StatAttack))*m.PrayerModifiers()[entity.StatAttack])+float64(m.StyleBonus(entity.StatAttack)))*((float64(m.AimPoints())*0.00175)+0.1) + 1.05
}

//DefensePoints Calculates and returns the defensive capability of this mob, based on many variables, as a single variable.
func (m *Mob) DefensePoints() float64 {
	return ((float64(m.Skills().Current(entity.StatDefense))*m.PrayerModifiers()[entity.StatDefense])+float64(m.StyleBonus(entity.StatDefense)))*((float64(m.ArmourPoints())*0.00175)+0.1) + 1.05
}

func (m *Mob) CombatRng() *rand.Rand {
	rng, ok := m.VarChecked("isaacRng").(*rand.Rand)
	if !ok || rng == nil {
		rng = rand.New(isaac.New(rscRand.Rng.Uint64()))
		m.SetVar("isaacRng", rng)
	}
	return rng
}

func (m *Mob) Isaac() *rand.Rand {
	rng, ok := m.VarChecked("isaac").(*rand.Rand)
	if !ok || rng == nil {
		rng = rand.New(isaac.New(rscRand.Rng.Uint64()))
		m.SetVar("isaac", rng)
	}
	return rng
}

//MagicDamage Calculates and returns a combat spells damage from the
// receiver mob cast unto the target MobileEntity.  This basically wraps a statistically
// random percentage check around a call to GenerateHit.
func (m *Mob) MagicDamage(target entity.MobileEntity, maximum float64) int {
	// TODO: Tweak the defense/armor hit/miss formula to better match RSC--or at least verify this is somewhat close?
	if BoundedChance(float64(m.MagicPoints())/(target.DefensePoints()*4)*100, 0.0, 82.0) {
		return m.GenerateHit(maximum)
	}

	return 0
}

//GenerateHit returns a normally distributed random number from this mobs PRNG instance,
// between 1 and max, inclusive.  This is widely believed to be how Jagex generated damage hits,
// and it feels accurate while playing.
func (m *Mob) GenerateHit(max float64) int {
	var damage float64
	for damage > max || damage < 1.0 {
		damage = math.Floor((m.CombatRng().NormFloat64() * (max / 3)) + (max / 2))
	}

	return int(damage)
}

//MeleeDamage Calculates and returns a melee damage from the receiver mob onto the target mob.
// This basically wraps a statistically random percentage check around a call to GenerateHit.
// Generally believed to be a near-perfect approximation of canonical Jagex RSClassic melee damage formula.
// Kenix mentioned running monte-carlo sims when coming up with it, so presumably this formula matched up
// statistically fairly well to the real game.  I can not say for sure as I didn't do these things myself, though.
func (m *Mob) MeleeDamage(target entity.MobileEntity) int {
	if BoundedChance(m.AttackPoints()/(target.DefensePoints()*4)*100, 0.0, 82.0) {
		return m.GenerateHit(m.MaxMeleeDamage())
	}

	return 0
}

func (m *Mob) Random(low, high int) int {
	return int(m.Isaac().Int63n(int64(high-low))) + low
}

type HitSplat struct {
	Owner entity.MobileEntity
	Damage int
}

func NewHitsplat(target entity.MobileEntity, damage int) HitSplat {
	return HitSplat{target, damage}
}
