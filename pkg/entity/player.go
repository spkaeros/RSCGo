package entity

//Player Represents a single player.
type Player struct {
	location          *Location
	state             MobState
	direction         Direction
	Username          string
	UserBase37        uint64
	Password          string
	Index             int
	Path              *Pathway
	FriendList        []uint64
	LocalPlayers      *LocatableList
	LocalObjects      *LocatableList
	HasMoved          bool
	Removing          bool
	HasSelf           bool
	AppearanceChanged bool
	Appearances       []int
	Skillset          *SkillTable
	DatabaseIndex     int
	Rank              int
	Attributes        map[string]interface{}
	Appearance        *AppearanceTable
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

//NewAppearance Returns a reference to a new appearance with specified parameters
func NewAppearanceTable(head, body int, male bool, hair, top, bottom, skin int) *AppearanceTable {
	return &AppearanceTable{head, body, male, hair, top, bottom, skin}
}

//SkillTable Represents a skill table for a mob.
type SkillTable struct {
	Current    [18]int
	Maximum    [18]int
	Experience [18]int
}

//ArmourPoints Returns the players armour points.
func (p *Player) ArmourPoints() int {
	if _, ok := p.Attributes["armour_points"]; !ok {
		p.Attributes["armour_points"] = 1
	}
	return p.Attributes["armour_points"].(int)
}

//SetArmourPoints Sets the players armour points to i.
func (p *Player) SetArmourPoints(i int) {
	p.Attributes["armour_points"] = i
}

//PowerPoints Returns the players power points.
func (p *Player) PowerPoints() int {
	if _, ok := p.Attributes["power_points"]; !ok {
		p.Attributes["power_points"] = 1
	}
	return p.Attributes["power_points"].(int)
}

//SetPowerPoints Sets the players power points to i
func (p *Player) SetPowerPoints(i int) {
	p.Attributes["power_points"] = i
}

//AimPoints Returns the players aim points
func (p *Player) AimPoints() int {
	if _, ok := p.Attributes["aim_points"]; !ok {
		p.Attributes["aim_points"] = 1
	}
	return p.Attributes["aim_points"].(int)
}

//SetAimPoints Sets the players aim points to i.
func (p *Player) SetAimPoints(i int) {
	p.Attributes["aim_points"] = i
}

//MagicPoints Returns the players magic points
func (p *Player) MagicPoints() int {
	if _, ok := p.Attributes["magic_points"]; !ok {
		p.Attributes["magic_points"] = 1
	}
	return p.Attributes["magic_points"].(int)
}

//SetMagicPoints Sets the players magic points to i
func (p *Player) SetMagicPoints(i int) {
	p.Attributes["magic_points"] = i
}

//PrayerPoints Returns the players prayer points
func (p *Player) PrayerPoints() int {
	if _, ok := p.Attributes["prayer_points"]; !ok {
		p.Attributes["prayer_points"] = 1
	}
	return p.Attributes["prayer_points"].(int)
}

//SetPrayerPoints Sets the players prayer points to i
func (p *Player) SetPrayerPoints(i int) {
	p.Attributes["prayer_points"] = i
}

//RangedPoints Returns the players ranged points.
func (p *Player) RangedPoints() int {
	if _, ok := p.Attributes["ranged_points"]; !ok {
		p.Attributes["ranged_points"] = 1
	}
	return p.Attributes["ranged_points"].(int)
}

//SetRangedPoints Sets the players ranged points tp i.
func (p *Player) SetRangedPoints(i int) {
	p.Attributes["ranged_points"] = i
}

//Fatigue Returns the players current fatigue.
func (p *Player) Fatigue() int {
	if _, ok := p.Attributes["fatigue"]; !ok {
		p.Attributes["fatigue"] = 0
	}
	return p.Attributes["fatigue"].(int)
}

//SetFatigue Sets the players current fatigue to i.
func (p *Player) SetFatigue(i int) {
	p.Attributes["fatigue"] = i
}

//FightMode Returns the players current fight mode.
func (p *Player) FightMode() int {
	if _, ok := p.Attributes["fight_mode"]; !ok {
		p.Attributes["fight_mode"] = 0
	}
	return p.Attributes["fight_mode"].(int)
}

//SetFightMode Sets the players fightmode to i.  0=all,1=attack,2=defense,3=strength
func (p *Player) SetFightMode(i int) {
	p.Attributes["fight_mode"] = i
}

//X Shortcut for Location().X()
func (p *Player) X() int {
	return p.location.X
}

//Y Shortcut for Location().Y()
func (p *Player) Y() int {
	return p.location.Y
}

//TraversePath If the player has a path, calling this method will change the players location to the next location
//  described by said Path data structure.  This should be called no more than once per game tick.
func (p *Player) TraversePath() {
	if p == nil || p.Path == nil {
		return
	}
	path := p.Path
	if p.AtLocation(path.Waypoint(path.CurrentWaypoint)) {
		path.CurrentWaypoint++
	}
	newLocation := path.NextTile(p.X(), p.Y())
	if path.CurrentWaypoint >= len(path.WaypointsX) || newLocation.X == -1 || newLocation.Y == -1 {
		p.ClearPath()
		return
	}
	p.HasMoved = true
	p.SetLocation(newLocation)
}

//NearbyPlayers Returns the nearby players from the current and nearest adjacent regions in a slice.
func (p *Player) NearbyPlayers() (players []*Player) {
	for _, r := range SurroundingRegions(p.X(), p.Y()) {
		for _, p1 := range r.Players {
			if p1.Index != p.Index && p.location.LongestDelta(p1.location) <= 15 {
				players = append(players, p1)
			}
		}
	}

	return
}

//ClearPath Sets the players path to nil, to stop the traversal of the path instantly
func (p *Player) ClearPath() {
	p.Path = nil
}

//Location Returns the location of the player
func (p *Player) Location() *Location {
	return p.location
}

//UpdateDirection Updates the direction the player is facing based on where the player is trying to move, and
// where the player is currently at.
func (p *Player) UpdateDirection(destX, destY int) {
	sprites := [3][3]int{{3, 2, 1}, {4, -1, 0}, {5, 6, 7}}
	xIndex := p.X() - destX + 1
	yIndex := p.Y() - destY + 1
	if xIndex >= 0 && yIndex >= 0 && xIndex < 3 && yIndex < 3 {
		p.direction = Direction(sprites[xIndex][yIndex])
	} else {
		p.direction = 0
	}
}

//SetLocation Sets the players location.
func (p *Player) SetLocation(location *Location) {
	p.SetCoords(location.X, location.Y)
}

//SetCoords Sets the players locations coordinates.
func (p *Player) SetCoords(x, y int) {
	curArea := GetRegion(p.X(), p.Y())
	newArea := GetRegion(x, y)
	if newArea != curArea {
		curArea.RemovePlayer(p)
		newArea.AddPlayer(p)
	}
	p.UpdateDirection(x, y)
	p.location.X = x
	p.location.Y = y
}

//AtLocation Returns true if the player is at the specified location, otherwise returns false
func (p *Player) AtLocation(location *Location) bool {
	return p.AtCoords(location.X, location.Y)
}

//AtCoords Returns true if the player is at the specified coordinates, otherwise returns false
func (p *Player) AtCoords(x, y int) bool {
	return p.location.X == x && p.location.Y == y
}

//State Returns the players state.
func (p *Player) State() MobState {
	return p.state
}

//SetState Sets the players state.
func (p *Player) SetState(state MobState) {
	p.state = state
}

//Direction Returns the players direction.
func (p *Player) Direction() Direction {
	return p.direction
}

//SetDirection Sets the players direction.
func (p *Player) SetDirection(direction Direction) {
	p.direction = direction
}

//NewPlayer Returns a reference to a new player.
func NewPlayer() *Player {
	return &Player{location: &Location{0, 0}, direction: North, state: Idle, Attributes: make(map[string]interface{}), LocalPlayers: &LocatableList{}, LocalObjects: &LocatableList{}, Skillset: &SkillTable{}, Appearance: NewAppearanceTable(1, 2, true, 2, 8, 14, 0)}
}
