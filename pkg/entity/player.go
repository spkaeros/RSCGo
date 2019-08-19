package entity

//Player Represents a single player.
type Player struct {
	location  *Location
	state     MobState
	direction Direction
	Username  string
	Password  string
	Path      *Pathway
}

type Pathway struct {
	StartX, StartY  int
	WaypointsX      []int
	WaypointsY      []int
	CurrentWaypoint int
}

func (p *Pathway) WaypointXoffset(w int) int {
	if w >= len(p.WaypointsX) || w == -1 {
		return 0
	}
	return p.WaypointsX[w]
}

func (p *Pathway) WaypointX(w int) int {
	return p.StartX + p.WaypointXoffset(w)
}

func (p *Pathway) WaypointYoffset(w int) int {
	if w >= len(p.WaypointsY) || w == -1 {
		return 0
	}
	return p.WaypointsY[w]
}

func (p *Pathway) WaypointY(w int) int {
	return p.StartY + p.WaypointYoffset(w)
}

func (p *Pathway) NextTile(startX, startY int) (nextCoords [2]int) {
	destX := p.WaypointX(p.CurrentWaypoint)
	destY := p.WaypointY(p.CurrentWaypoint)
	nextCoords = [2]int{-1, -1}
	if startX > destX {
		nextCoords[0] = startX - 1
	} else if startX < destX {
		nextCoords[0] = startX + 1
	} else {
		nextCoords[0] = destX
	}
	if startY > destY {
		nextCoords[1] = startY - 1
	} else if startY < destY {
		nextCoords[1] = startY + 1
	} else {
		nextCoords[1] = destY
	}
	return nextCoords
}

func (p *Player) X() int {
	return p.location.x
}

func (p *Player) Y() int {
	return p.location.y
}

func (p *Player) SetX(x int) {
	p.location.x = x
}

func (p *Player) SetY(y int) {
	p.location.y = y
}

func (p *Player) TraversePath() {
	if p.Path.CurrentWaypoint == -1 {
		if p.Path.StartX == p.X() && p.Path.StartY == p.Y() {
			p.Path.CurrentWaypoint = 0
		} else {
			nextCoords := p.Path.NextTile(p.X(), p.Y())
			if nextCoords[0] != -1 && nextCoords[1] != -1 {
				p.SetX(nextCoords[0])
				p.SetY(nextCoords[1])
			}
		}
	}
	if p.Path.CurrentWaypoint > -1 {
		if p.X() == p.Path.WaypointX(p.Path.CurrentWaypoint) && p.Y() == p.Path.WaypointY(p.Path.CurrentWaypoint) {
			p.Path.CurrentWaypoint++
		}
		if p.Path.CurrentWaypoint < len(p.Path.WaypointsX) {
			nextCoords := p.Path.NextTile(p.X(), p.Y())
			if nextCoords[0] != -1 && nextCoords[1] != -1 {
				p.SetX(nextCoords[0])
				p.SetY(nextCoords[1])
			}
		} else {
			p.Path = nil
		}
	}
}

//Location Returns the location of the player
func (p *Player) Location() *Location {
	return p.location
}

//SetLocation Sets the players location.
func (p *Player) SetLocation(location *Location) {
	p.location = location
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
