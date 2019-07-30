package entity

type Player struct {
	location Location
	state MobState
	direction Direction
}

func (p *Player) Location() Location {
	return p.location
}

func (p *Player) SetLocation(location Location) {
	p.location = location
}

func (p *Player) State() MobState {
	return p.state
}

func (p *Player) SetState(state MobState) {
	p.state = state
}

func (p *Player) Direction() Direction {
	return p.direction
}

func (p *Player) SetDirection(direction Direction) {
	p.direction = direction
}

func NewPlayer() *Player {
	return &Player{location: Location{220, 445}, direction: North, state: Idle}
}