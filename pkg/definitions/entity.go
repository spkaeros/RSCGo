package definitions

//ItemDefinition This represents a single definition for a single item in the game.
type ItemDefinition struct {
	ID           int
	Name         string
	Description  string
	Command      string
	BasePrice    int
	Stackable    bool
	Quest        bool
	Members      bool
	Requirements map[int]int
}

//Items This holds the defining characteristics for all of the game's items, ordered by ID.
var Items []ItemDefinition

//Item returns the associated item definition or nil if none.
func Item(id int) ItemDefinition {
	for _, i := range Items {
		if i.ID == id {
			return i
		}
	}

	return ItemDefinition{ID: -1}
}

//EquipmentDefinition a container for the equipment items in the game
type EquipmentDefinition struct {
	ID       int
	Sprite   int
	Type     int
	Armour   int
	Magic    int
	Prayer   int
	Ranged   int
	Aim      int
	Power    int
	Position int
	Female   bool
}

//Equipment contains all of the equipment related data for the game
var Equipment []EquipmentDefinition

//Equip returns the associated equipment definition or nil if none.
func Equip(id int) *EquipmentDefinition {
	for _, e := range Equipment {
		if e.ID == id {
			return &e
		}
	}

	return nil
}

//DefaultDrop returns the default item ID all mobs should drop on death
const DefaultDrop = 20

//NpcDefinition This is a representation of a single NPC in the game.
type NpcDefinition struct {
	ID          int
	Name        string
	Description string
	Command     string
	Hits        int
	Attack      int
	Strength    int
	Defense     int
	Hostility   int
}

//Npcs This is used to cache the persistent NPC data in RAM for quick access when needed.
var Npcs []NpcDefinition

//Npc returns the associated NPC definition or nil if none.
func Npc(id int) NpcDefinition {
	for _, e := range Npcs {
		if e.ID == id {
			return e
		}
	}

	return NpcDefinition{ID: -1}
}

//ObjectDefinition This represents a single definition for a single object in the game.
type ScenaryDefinition struct {
	ID            int
	Name          string
	Commands      [2]string
	Description   string
	CollisionType int
	Width, Height int
	ModelHeight   int
}

//ScenaryObjects This is used to cache the persistent scenary object data in RAM for quick access when needed.
var ScenaryObjects []ScenaryDefinition

func Scenary(id int) ScenaryDefinition {
	for _, o := range ScenaryObjects {
		if o.ID == id {
			return o
		}
	}

	return ScenaryDefinition{ID: -1}
}

//TileDefinition Representation of a tile overlay.
type TileDefinition struct {
	Color   int
	Visible int
	Blocked int
}

//TileOverlays Cache for tile overlays.
var TileOverlays []TileDefinition

func TileOverlay(id int) TileDefinition {
	if id < len(TileOverlays) && id > 0 {
		return TileOverlays[id]
	}

	return TileDefinition{Blocked: 1, Visible: 0, Color: 987654321}
}

//BoundaryDefinition This represents a single definition for a single boundary object in the game.
type BoundaryDefinition struct {
	ID          int
	Name        string
	Commands    [2]string
	Description string
	Dynamic     bool
	Solid       bool
}

//BoundaryObjectss This holds the defining characteristics for all of the game's boundary scene objects, ordered by ID.
var BoundaryObjects []BoundaryDefinition

func Boundary(id int) BoundaryDefinition {
	for _, b := range BoundaryObjects {
		if b.ID == id {
			return b
		}
	}

	return BoundaryDefinition{ID: -1}
}

const (
	OverlayBlank = iota
	//OverlayGravel Used for roads, ID 1
	OverlayGravel
	//OverlayWater Used for regular water, ID 2
	OverlayWater
	//OverlayWoodFloor Used for the floors of buildings, ID 3
	OverlayWoodFloor
	//OverlayBridge Used for bridges, suspends wood floor over water, ID 4
	OverlayBridge
	//OverlayStoneFloor Used for the floors of buildings, ID 5
	OverlayStoneFloor
	//OverlayRedCarpet Used for the floors of buildings, ID 6
	OverlayRedCarpet
	//OverlayDarkWater Used for dark, swampy water, ID 7
	OverlayDarkWater
	//OverlayBlack Used for empty parts of upper planes, ID 8
	OverlayBlack
	//OverlayWhite Used as a separator, e.g for edge of water, mountains, etc.  ID 9
	OverlayWhite
	//OverlayBlack2 Not sure where it is used, ID 10
	OverlayBlack2
	//OverlayLava Used in dungeons and on Karamja/Crandor as lava, ID 11
	OverlayLava
	//OverlayBridge2 Used for a specific type of bridge, ID 12
	OverlayBridge2
	//OverlayBlueCarpet Used for the floors of buildings, ID 13
	OverlayBlueCarpet
	//OverlayPentagram Used for certain questing purposes, ID 14
	OverlayPentagram
	//OverlayPurpleCarpet Used for the floors of buildings, ID 15
	OverlayPurpleCarpet
	//OverlayBlack3 Not sure what it is used for, ID 16, traversable
	OverlayBlack3
	//OverlayStoneFloorLight Used for the entrance to temple of ikov, ID 17
	OverlayStoneFloorLight
	//OverlayUnknown Not sure what this is yet, ID 18
	OverlayUnknown
	//OverlayBlack4 Not sure what it is used for, ID 19
	OverlayBlack4
	//OverlayAgilityLog Blank suspended tile over blackness for agility challenged, ID 20
	OverlayAgilityLog
	//OverlayAgilityLog Blank suspended tile over blackness for agility challenged, ID 21
	OverlayAgilityLog2
	//OverlayUnknown2 Not sure what this is yet, ID 22
	OverlayUnknown2
	//OverlaySandFloor Used for sand floor, ID 23
	OverlaySandFloor
	//OverlayMudFloor Used for mud floor, ID 24
	OverlayMudFloor
	//OverlaySandFloor Used for water floor, ID 25
	OverlayWaterFloor
)

//blockedOverlays An array filled with any overlay types that mobs aren't able to walk to from another tile
var blockedOverlays = [...]int{OverlayWater, OverlayDarkWater, OverlayBlack, OverlayWhite, OverlayLava, OverlayBlack2, OverlayBlack3, OverlayBlack4}
