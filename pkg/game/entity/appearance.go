package entity

//AppearanceTable Represents a players appearance.
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

//NewAppearanceTable returns a reference to a new appearance table with specified parameters
func NewAppearanceTable(head, body int, male bool, hair, top, bottom, skin int) AppearanceTable {
	// only one legs, idx 3
	return AppearanceTable{head, body, 3, male, hair, top, bottom, skin}
}

func DefaultAppearance() AppearanceTable {
	return NewAppearanceTable(1, 2, true, 2, 8, 14, 0)
}
