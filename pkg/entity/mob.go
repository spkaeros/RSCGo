package entity

//MobState Mob state.
type MobState uint8

const (
	//Idle The default MobState, means doing nothing.
	Idle MobState = iota
	//Walking The mob is walking.
	Walking
	//Banking The mob is banking.
	Banking
	//Chatting The mob is chatting with a NPC
	Chatting
	//MenuChoosing The mob is in a query menu
	MenuChoosing
	//Trading The mob is negotiating a trade.
	Trading
	//Dueling The mob is negotiating a duel.
	Dueling
	//Fighting The mob is fighting.
	Fighting
	//Batching The mob is performing a skill that repeats itself an arbitrary number of times.
	Batching
	//Sleeping The mob is using a bed or sleeping bag, and trying to solve a CAPTCHA
	Sleeping
)

//Mob Represents a mobile entity within the game world.
type Mob struct {
	Entity
	State      MobState
	Skillset   *SkillTable
	Path       *Pathway
	Attributes AttributeList
	TransAttrs AttributeList
}

//Direction Returns the mobs direction.
func (m *Mob) Direction() int {
	return m.Attributes.VarInt("direction", 0)
}

//SetDirection Sets the mobs direction.
func (m *Mob) SetDirection(direction int) {
	m.Attributes["direction"] = direction
}

//SetPath Sets the mob's current pathway to path.  If path is nil, effectively resets the mobs path.
func (m *Mob) SetPath(path *Pathway) {
	m.Path = path
}

//ResetPath Sets the mobs path to nil, to stop the traversal of the path instantly
func (m *Mob) ResetPath() {
	m.SetPath(nil)
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

//AtLocation Returns true if the mob is at the specified location, otherwise returns false
func (m *Mob) AtLocation(location *Location) bool {
	return m.AtCoords(location.X, location.Y)
}

//AtCoords Returns true if the mob is at the specified coordinates, otherwise returns false
func (m *Mob) AtCoords(x, y int) bool {
	return m.X == x && m.Y == y
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
	curArea := GetRegion(m.X, m.Y)
	newArea := GetRegion(x, y)
	if newArea != curArea {
		curArea.Players.Remove(m)
		newArea.Players.Add(m)
	}
	m.UpdateDirection(x, y)
	m.X = x
	m.Y = y
}

//TODO: Probably remove the Attribute type-alias.
//Attribute Type-alias for attribute names.  Might not need this, was just so to provide methods for them, don't think I'm doing it anymore
type Attribute string

//AttributeList A type alias for a map of strings to empty interfaces, to hold generic mob information for easy serialization and to provide dynamic insertion/deletion of new mob properties easily
type AttributeList map[Attribute]interface{}

//SetVar Sets the attribute mapped at name to value in the attribute map.
func (attributes AttributeList) SetVar(name Attribute, value interface{}) {
	attributes[name] = value
}

//VarInt If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes AttributeList) VarInt(name Attribute, zero int) int {
	if _, ok := attributes[name].(int); !ok {
		attributes[name] = zero
	}

	return attributes[name].(int)
}

//VarLong If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes AttributeList) VarLong(name Attribute, zero uint64) uint64 {
	if _, ok := attributes[name].(uint64); !ok {
		attributes[name] = zero
	}

	return attributes[name].(uint64)
}

//VarBool If there is an attribute assigned to the specified name, returns it.  Otherwise, returns zero
func (attributes AttributeList) VarBool(name Attribute, zero bool) bool {
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
