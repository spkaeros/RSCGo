package entity

//Direction Direction within gameworld.
type Direction uint8

const (
	//North Represents north.
	North Direction = iota
	//NorthWest Represents north-west.
	NorthWest
	//West Represents west.
	West
	//SouthWest Represents south-west.
	SouthWest
	//South represents south.
	South
	//SouthEast represents south-east
	SouthEast
	//East Represents east.
	East
	//NorthEast Represents north-east.
	NorthEast
	//MaxX Width of the game
	MaxX = 944
	//MaxY Height of the game
	MaxY = 3776
)

var mobSprites = [][]int{{3, 2, 1}, {4, -1, 0}, {5, 6, 7}}

func (p *Player) UpdateDirection(destX, destY int) {
	xIndex := p.X() - destX + 1
	yIndex := p.Y() - destY + 1
	if xIndex >= 0 && yIndex >= 0 && xIndex < 3 && yIndex < 3 {
		p.direction = Direction(mobSprites[xIndex][yIndex])
	} else {
		p.direction = 0
	}
}

//Location A tile in the game world.
type Location struct {
	x, y int
}

//X Returns the X coordinate of this entity
func (l *Location) X() int {
	return l.x
}

func (l *Location) SetX(x int) {
	l.x = x
}

func (l *Location) SetY(y int) {
	l.y = y
}

//Y Returns the Y coordinate of this entity
func (l *Location) Y() int {
	return l.y
}

//Locatable An interface for locatable entities
type Locatable interface {
	Location() Location
}

//Entity Any data structure that represents something in the game world should be able to implement this interface,
type Entity interface {
	Locatable
	// Returns the direction the given Entity is facing in the game world.
	Direction() int
}
