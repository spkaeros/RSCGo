package entity

//Player Represents a single player.
type Player struct {
	location  Location
	state     MobState
	direction Direction
}

//Location Returns the location of the player
func (p *Player) Location() Location {
	return p.location
}

//SetLocation Sets the players location.
func (p *Player) SetLocation(location Location) {
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
	return &Player{location: Location{220, 445}, direction: North, state: Idle}
}
