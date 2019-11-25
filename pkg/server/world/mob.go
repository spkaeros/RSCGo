package world

import (
	"github.com/spkaeros/rscgo/pkg/rand"
	"go.uber.org/atomic"
	"sync"
	"time"
)

//MobState Mob state.
type MobState uint8

const (
	//MSIdle The default MobState, means doing nothing.
	MSIdle MobState = iota
	//MSWalking The mob is walking.
	MSWalking
	//MSBanking The mob is banking.
	MSBanking
	//MSChatting The mob is chatting with a NPC
	MSChatting
	//MSMenuChoosing The mob is in a query menu
	MSMenuChoosing
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

//Mob Represents a mobile entity within the game world.
type Mob struct {
	*Entity
	State      MobState
	Skillset   *SkillTable
	TransAttrs *AttributeList
}

type Locatable interface {
	CurX() int
	CurY() int
}

//Busy Returns true if this mobs state is anything other than idle. otherwise returns false.
func (m *Mob) Busy() bool {
	return m.State != MSIdle
}

//Direction Returns the mobs direction.
func (m *Mob) Direction() int {
	return m.TransAttrs.VarInt("direction", North)
}

//SetDirection Sets the mobs direction.
func (m *Mob) SetDirection(direction int) {
	m.Change()
	m.TransAttrs.SetVar("direction", direction)
}

//Change Sets the synchronization flag for whether this mob changed directions to true.
func (m *Mob) Change() {
	m.TransAttrs.SetVar("changed", true)
}

//Remove Sets the synchronization flag for whether this mob needs to be removed to true.
func (m *Mob) Remove() {
	m.TransAttrs.SetVar("remove", true)
}

//UpdateSelf Sets the synchronization flag for whether this mob needs to update itself to true.
func (m *Mob) UpdateSelf() {
	m.TransAttrs.SetVar("self", true)
}

//UpdateSelf Sets the synchronization flag for whether this mob has moved to true.
func (m *Mob) Move() {
	m.TransAttrs.SetVar("moved", true)
}

func (m *Mob) ResetMoved() {
	m.TransAttrs.UnsetVar("moved")
}

func (m *Mob) ResetRemoved() {
	m.TransAttrs.UnsetVar("remove")
}

func (m *Mob) ResetNeedsSelf() {
	m.TransAttrs.UnsetVar("self")
}

func (m *Mob) ResetChanged() {
	m.TransAttrs.UnsetVar("changed")
}

//SetPath Sets the mob's current pathway to path.  If path is nil, effectively resets the mobs path.
func (m *Mob) SetPath(path *Pathway) {
	m.TransAttrs.SetVar("path", path)
}

//Path returns the path that this mob is trying to traverse.
func (m *Mob) Path() *Pathway {
	return m.TransAttrs.VarPath("path")
}

//ResetPath Sets the mobs path to nil, to stop the traversal of the path instantly
func (m *Mob) ResetPath() {
	m.ResetMoved()
	m.TransAttrs.UnsetVar("path")
}

//TraversePath If the mob has a path, calling this method will change the mobs location to the next location described by said Path data structure.  This should be called no more than once per game tick.
func (m *Mob) TraversePath() Location {
	path := m.Path()
	if path == nil {
		return m.Location
	}
	if m.AtLocation(path.NextWaypointTile()) {
		path.CurrentWaypoint++
	}
	if m.FinishedPath() {
		m.ResetPath()
		return m.Location
	}
	dst := path.NextWaypointTile()
	x, y := m.CurX(), m.CurY()
	next := NewLocation(x, y)
	xBlocked, yBlocked := false, false
	newXBlocked, newYBlocked := false, false
	if y > dst.CurY() {
		yBlocked = IsTileBlocking(x, y, 1, true)
		newYBlocked = IsTileBlocking(x, y-1, 4, false)
		if !newYBlocked {
			next.Y.Dec()
		}
	} else if y < dst.CurY() {
		yBlocked = IsTileBlocking(x, y, 4, true)
		newYBlocked = IsTileBlocking(x, y+1, 1, false)
		if !newYBlocked {
			next.Y.Inc()
		}
	}
	if x > dst.CurX() {
		xBlocked = IsTileBlocking(x, next.CurY(), 2, true)
		newXBlocked = IsTileBlocking(x-1, next.CurY(), 8, false)
		if !newXBlocked {
			next.X.Dec()
		}
	} else if x < dst.CurX() {
		xBlocked = IsTileBlocking(x, next.CurY(), 8, true)
		newXBlocked = IsTileBlocking(x+1, next.CurY(), 2, false)
		if !newXBlocked {
			next.X.Inc()
		}
	}

	if (xBlocked && yBlocked) || (xBlocked && y == dst.CurY()) || (yBlocked && x == dst.CurX()) {
		m.ResetPath()
		return m.Location
	}
	if (newXBlocked && newYBlocked) || (newXBlocked && x != next.CurX() && y == next.CurY()) || (newYBlocked && y != next.CurY() && x == next.CurX()) {
		m.ResetPath()
		return m.Location
	}

	if next.CurX() > x {
		newXBlocked = IsTileBlocking(next.CurX(), next.CurY(), 2, false)
	} else if next.CurX() < x {
		newXBlocked = IsTileBlocking(next.CurX(), next.CurY(), 8, false)
	}
	if next.CurY() > y {
		newYBlocked = IsTileBlocking(next.CurX(), next.CurY(), 1, false)
	} else if next.CurY() < y {
		newYBlocked = IsTileBlocking(next.CurX(), next.CurY(), 4, false)
	}

	if (newXBlocked && newYBlocked) || (newXBlocked && y == next.CurY()) || (newYBlocked && x == next.CurX()) {
		m.ResetPath()
		return m.Location
	}

	return next
}

func UpdateRegionMob(m Locatable, x, y int) {
	curArea := GetRegion(m.CurX(), m.CurY())
	newArea := GetRegion(x, y)
	if newArea != curArea {
		if _, ok := m.(*Player); ok {
			if curArea.Players.Contains(m) {
				curArea.Players.Remove(m)
			}
			newArea.Players.Add(m)
		} else if _, ok := m.(*NPC); ok {
			if curArea.NPCs.Contains(m) {
				curArea.NPCs.Remove(m)
			}
			newArea.NPCs.Add(m)
		}
	}
}

//FinishedPath Returns true if the mobs path is nil, the paths current waypoint exceeds the number of waypoints available, or the next tile in the path is not a valid location, implying that we have reached our destination.
func (m *Mob) FinishedPath() bool {
	path := m.Path()
	if path == nil {
		return true
	}
	return path.CurrentWaypoint >= path.CountWaypoints() || !path.NextTileFrom(m.Location).IsValid()
}

//SetLocation Sets the mobs location.
func (m *Mob) SetLocation(location Location) {
	x := location.CurX()
	y := location.CurY()
	m.SetDirection(m.directionTo(x, y))
	m.SetCoords(x, y)
}

//SetCoords Sets the mobs locations coordinates.
func (m *Mob) SetCoords(x, y int) {
//	m.Location = NewLocation(x, y)
	m.SetX(x)
	m.SetY(y)
}

//Teleport Moves the mob to x,y and sets a flag to remove said mob from the local players list of every nearby player.
func (m *Mob) Teleport(x, y int) {
	m.Remove()
	m.SetCoords(x, y)
}

//AttrList A type alias for a map of strings to empty interfaces, to hold generic mob information for easy serialization and to provide dynamic insertion/deletion of new mob properties easily
type AttrList map[string]interface{}

//AttributeList A concurrency-safe collection data type for storing misc. variables by a descriptive name.
type AttributeList struct {
	Set  map[string]interface{}
	Lock sync.RWMutex
}

//Range Runs fn(key, value) for every entry in this attribute list.
func (attributes *AttributeList) Range(fn func(string, interface{})) {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	for k, v := range attributes.Set {
		fn(k, v)
	}
}

//SetVar Sets the attribute mapped at name to value in the attribute map.
func (attributes *AttributeList) SetVar(name string, value interface{}) {
	attributes.Lock.Lock()
	attributes.Set[name] = value
	attributes.Lock.Unlock()
}

//UnsetVar Removes the attribute with the key `name` from this attribute set.
func (attributes *AttributeList) UnsetVar(name string) {
	attributes.Lock.Lock()
	delete(attributes.Set, name)
	attributes.Lock.Unlock()
}

//VarInt If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes *AttributeList) VarInt(name string, zero int) int {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	if _, ok := attributes.Set[name].(int); !ok {
		return zero
	}

	return attributes.Set[name].(int)
}

//VarPlayer If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes *AttributeList) VarPlayer(name string) *Player {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	if _, ok := attributes.Set[name].(*Player); !ok {
		return nil
	}

	return attributes.Set[name].(*Player)
}

//VarPlayer If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes *AttributeList) VarNpc(name string) *NPC {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	if _, ok := attributes.Set[name].(*NPC); !ok {
		return nil
	}

	return attributes.Set[name].(*NPC)
}

//VarLong If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes *AttributeList) VarLong(name string, zero uint64) uint64 {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	if _, ok := attributes.Set[name].(uint64); !ok {
		return zero
	}

	return attributes.Set[name].(uint64)
}

//VarBool If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes *AttributeList) VarBool(name string, zero bool) bool {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	if _, ok := attributes.Set[name].(bool); !ok {
		return zero
	}

	return attributes.Set[name].(bool)
}

//VarTime If there is a time.Duration attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes *AttributeList) VarTime(name string) time.Time {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	if _, ok := attributes.Set[name].(time.Time); !ok {
		return time.Time{}
	}

	return attributes.Set[name].(time.Time)
}

//VarTime If there is a time.Duration attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes *AttributeList) VarPath(name string) *Pathway {
	attributes.Lock.RLock()
	defer attributes.Lock.RUnlock()
	if _, ok := attributes.Set[name].(*Pathway); !ok {
		return nil
	}

	return attributes.Set[name].(*Pathway)
}

//AppearanceTable Represents a mobs appearance.
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

//NewAppearanceTable Returns a reference to a new appearance table with specified parameters
func NewAppearanceTable(head, body int, male bool, hair, top, bottom, skin int) AppearanceTable {
	return AppearanceTable{head, body, 3, male, hair, top, bottom, skin}
}

//SkillTable Represents a skill table for a mob.
type SkillTable struct {
	current    [18]int
	maximum    [18]int
	experience [18]int
	Lock       sync.RWMutex
}

//Current Returns the current level of the skill indicated by idx.
func (s *SkillTable) Current(idx int) int {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	return s.current[idx]
}

func (s *SkillTable) DecreaseCur(idx, delta int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.current[idx] -= delta
}

func (s *SkillTable) IncreaseCur(idx, delta int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.current[idx] += delta
}

func (s *SkillTable) SetCur(idx, val int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.current[idx] = val
}

func (s *SkillTable) DecreaseMax(idx, delta int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.maximum[idx] -= delta
}

func (s *SkillTable) IncreaseMax(idx, delta int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.maximum[idx] += delta
}

func (s *SkillTable) SetMax(idx, val int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.maximum[idx] = val
}

func (s *SkillTable) SetExp(idx, val int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.experience[idx] = val
}

func (s *SkillTable) IncExp(idx, val int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.experience[idx] += val
}

//Maximum Returns the maximum level of the skill indicated by idx.
func (s *SkillTable) Maximum(idx int) int {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	return s.maximum[idx]
}

//Experience Returns the current level of the skill indicated by idx.
func (s *SkillTable) Experience(idx int) int {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	return s.experience[idx]
}

//CombatLevel Calculates and returns the combat level for this skill table.
func (s *SkillTable) CombatLevel() int {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	aggressiveTotal := float32(s.maximum[0] + s.maximum[2])
	defensiveTotal := float32(s.maximum[1] + s.maximum[3])
	spiritualTotal := float32((s.maximum[5] + s.maximum[6]) / 8)
	ranged := float32(s.maximum[4])
	if aggressiveTotal < ranged*1.5 {
		return int((defensiveTotal / 4) + (ranged * 0.375) + spiritualTotal)
	}
	return int((aggressiveTotal / 4) + (defensiveTotal / 4) + spiritualTotal)
}

//NpcDefinition This represents a single definition for a single NPC in the game.
type NpcDefinition struct {
	ID          int
	Name        string
	Description string
	Command     string
	Hits        int
	Attack      int
	Strength    int
	Defense     int
	Attackable  bool
}

//NpcDefs This holds the defining characteristics for all of the game's NPCs, ordered by ID.
var NpcDefs []NpcDefinition

//NpcCounter Counts the number of total NPCs within the world.
var NpcCounter = atomic.NewUint32(0)

//Npcs A collection of every NPC in the game, sorted by index
var Npcs []*NPC
var npcsLock sync.RWMutex

//NPC Represents a single non-playable character within the game world.
type NPC struct {
	Mob
	ID         int
	Boundaries [2]Location
	StartPoint Location
}

//NewNpc Creates a new NPC and returns a reference to it
func NewNpc(id int, startX int, startY int, minX, maxX, minY, maxY int) *NPC {
	n := &NPC{ID: id, Mob: Mob{Entity: &Entity{Index: int(NpcCounter.Swap(NpcCounter.Load() + 1)), Location: NewLocation(startX,startY)}, Skillset: &SkillTable{}, State: MSIdle, TransAttrs: &AttributeList{Set: make(map[string]interface{})}}}
	n.Boundaries[0] = NewLocation(minX, minY)
	n.Boundaries[1] = NewLocation(maxX, maxY)
	n.StartPoint = NewLocation(startX, startY)
	if id < 794 {
		n.Skillset.current[0] = NpcDefs[id].Attack
		n.Skillset.current[1] = NpcDefs[id].Defense
		n.Skillset.current[2] = NpcDefs[id].Strength
		n.Skillset.current[3] = NpcDefs[id].Hits
		n.Skillset.maximum[0] = NpcDefs[id].Attack
		n.Skillset.maximum[1] = NpcDefs[id].Defense
		n.Skillset.maximum[2] = NpcDefs[id].Strength
		n.Skillset.maximum[3] = NpcDefs[id].Hits
	}
	npcsLock.Lock()
	Npcs = append(Npcs, n)
	npcsLock.Unlock()
	return n
}

func UpdateNPCPaths() {
	npcsLock.RLock()
	for _, n := range Npcs {
		if n.LongestDelta(DeathSpot) == 0 {
			continue
		}
		if n.TransAttrs.VarBool("fighting", false) {
			continue
		}
		if n.TransAttrs.VarTime("nextMove").Before(time.Now()) {
			for _, r := range SurroundingRegions(n.CurX(), n.CurY()) {
				r.Players.lock.RLock()
				if len(r.Players.List) > 0 {
					r.Players.lock.RUnlock()
					n.TransAttrs.SetVar("nextMove", time.Now().Add(time.Second*time.Duration(rand.Int31N(5, 15))))
					n.SetPath(MakePath(n.Location, NewRandomLocation(n.Boundaries)))
					//n.SetPath(NewPathwayToLocation(NewRandomLocation(n.Boundaries)))
					break
				}
				r.Players.lock.RUnlock()
			}
		}
	}
	npcsLock.RUnlock()
}

//UpdateNPCPositions Loops through the global NPC list and, if they are by a player, updates their path to a new path every so often,
// within their boundaries, and traverses each NPC along said path if necessary.
func UpdateNPCPositions() {
	npcsLock.RLock()
	for _, n := range Npcs {

		nextTile := n.TraversePath()
		if nextTile.LongestDelta(n.Location) > 0 {
			UpdateRegionMob(n, nextTile.CurX(), nextTile.CurY())
			n.SetLocation(nextTile)
			n.Move()
		}
	}
	npcsLock.RUnlock()
}

//ResetNpcUpdateFlags Resets the synchronization update flags for all NPCs in the game world.
func ResetNpcUpdateFlags() {
	npcsLock.RLock()
	for _, n := range Npcs {
		n.TransAttrs.UnsetVar("changed")
		n.TransAttrs.UnsetVar("moved")
		n.TransAttrs.UnsetVar("remove")
	}
	npcsLock.RUnlock()
}
