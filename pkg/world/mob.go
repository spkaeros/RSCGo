package world

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
)

//Mob Represents a mobile entity within the game world.
type Mob struct {
	Entity
	State      MobState
	Skillset   *SkillTable
	Path       *Pathway
	TransAttrs AttributeList
}

//Direction Returns the mobs direction.
func (m *Mob) Direction() int {
	return m.TransAttrs.VarInt("direction", 0)
}

//SetDirection Sets the mobs direction.
func (m *Mob) SetDirection(direction int) {
	m.TransAttrs["direction"] = direction
}

//SetPath Sets the mob's current pathway to path.  If path is nil, effectively resets the mobs path.
func (m *Mob) SetPath(path *Pathway) {
	m.Path = path
}

//ResetPath Sets the mobs path to nil, to stop the traversal of the path instantly
func (m *Mob) ResetPath() {
	m.Path = nil
}

//TraversePath If the mob has a path, calling this method will change the mobs location to the next location described by said Path data structure.  This should be called no more than once per game tick.
func (m *Mob) TraversePath() {
	if m.Path == nil {
		return
	}
	path := m.Path
	if m.AtLocation(path.Waypoint(path.CurrentWaypoint)) {
		path.CurrentWaypoint++
	}
	newLocation := path.NextTile(m.X, m.Y)
	if path.CurrentWaypoint >= len(path.WaypointsX) || newLocation.X == -1 || newLocation.Y == -1 {
		m.ResetPath()
		return
	}
	m.TransAttrs["plrmoved"] = true
	m.SetLocation(newLocation)
}

//FinishedPath Returns true if the mobs path is nil or if we are already on the path's next tile.
func (m *Mob) FinishedPath() bool {
	if m.Path == nil {
		return true
	}
	next := m.Path.NextTile(m.X, m.Y)
	return m.AtLocation(&next)
}

//UpdateDirection Updates the direction the mob is facing based on where the mob is trying to move, and where the mob is currently at.
func (m *Mob) UpdateDirection(destX, destY int) {
	sprites := [3][3]int{{3, 2, 1}, {4, -1, 0}, {5, 6, 7}}
	xIndex := m.X - destX + 1
	yIndex := m.Y - destY + 1
	if xIndex >= 0 && yIndex >= 0 && xIndex < 3 && yIndex < 3 {
		m.SetDirection(sprites[xIndex][yIndex])
	} else {
		m.SetDirection(int(North))
	}
}

//SetLocation Sets the mobs location.
func (m *Mob) SetLocation(location Location) {
	m.SetCoords(location.X, location.Y)
}

//SetCoords Sets the mobs locations coordinates.
func (m *Mob) SetCoords(x, y int) {
	m.UpdateDirection(x, y)
	m.X = x
	m.Y = y
}

//AttributeList A type alias for a map of strings to empty interfaces, to hold generic mob information for easy serialization and to provide dynamic insertion/deletion of new mob properties easily
type AttributeList map[string]interface{}

//SetVar Sets the attribute mapped at name to value in the attribute map.
func (attributes AttributeList) SetVar(name string, value interface{}) {
	attributes[name] = value
}

//VarInt If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes AttributeList) VarInt(name string, zero int) int {
	if _, ok := attributes[name].(int); !ok {
		attributes[name] = zero
	}

	return attributes[name].(int)
}

//VarLong If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes AttributeList) VarLong(name string, zero uint64) uint64 {
	if _, ok := attributes[name].(uint64); !ok {
		attributes[name] = zero
	}

	return attributes[name].(uint64)
}

//VarBool If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes AttributeList) VarBool(name string, zero bool) bool {
	if _, ok := attributes[name].(bool); !ok {
		attributes[name] = zero
	}

	return attributes[name].(bool)
}

//AppearanceTable Represents a mobs appearance.
type AppearanceTable struct {
	Head   int
	Body   int
	Male   bool
	Hair   int
	Top    int
	Bottom int
	Skin   int
}

//NewAppearanceTable Returns a reference to a new appearance table with specified parameters
func NewAppearanceTable(head, body int, male bool, hair, top, bottom, skin int) *AppearanceTable {
	return &AppearanceTable{head, body, male, hair, top, bottom, skin}
}

//SkillTable Represents a skill table for a mob.
type SkillTable struct {
	Current    [18]int
	Maximum    [18]int
	Experience [18]int
}
