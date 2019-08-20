package entity

//Player Represents a single player.
type Player struct {
	location  *Location
	state     MobState
	direction Direction
	Username  string
	Password  string
	Index     int
	Path      *Pathway
}

//X Shortcut for Location().X()
func (p *Player) X() int {
	return p.location.x
}

//Y Shortcut for Location().Y()
func (p *Player) Y() int {
	return p.location.y
}

//SetX Shortcut for Location().SetX(int)
func (p *Player) SetX(x int) {
	p.location.x = x
}

//SetY Shortcut for Location().SetY(int)
func (p *Player) SetY(y int) {
	p.location.y = y
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
	if path.CurrentWaypoint >= len(path.WaypointsX) || newLocation.x == -1 || newLocation.y == -1 {
		p.ClearPath()
		return
	}
	p.SetLocation(newLocation)
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
	p.SetCoords(location.x, location.y)
}

//SetCoords Sets the players locations coordinates.
func (p *Player) SetCoords(x, y int) {
	curArea := GetRegion(p.X(), p.Y())
	newArea := GetRegion(x, y)
	if newArea != curArea {
		if _, ok := curArea.Players[p.Index]; ok {
			curArea.RemovePlayer(p)
		}
		newArea.AddPlayer(p)
	}
	p.UpdateDirection(x, y)
	p.location.x = x
	p.location.y = y
}

//AtLocation Returns true if the player is at the specified location, otherwise returns false
func (p *Player) AtLocation(location *Location) bool {
	return p.AtCoords(location.X(), location.Y())
}

//AtCoords Returns true if the player is at the specified coordinates, otherwise returns false
func (p *Player) AtCoords(x, y int) bool {
	return p.X() == x && p.Y() == y
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

//NewPlayer Returns a new player.
func NewPlayer() *Player {
	return &Player{location: &Location{220, 445}, direction: North, state: Idle}
}
