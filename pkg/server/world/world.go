package world

import (
	"fmt"
	"github.com/spkaeros/rscgo/pkg/server/log"
)

const (
	//RegionSize Represents the size of the region
	RegionSize = 48
	//HorizontalPlanes Represents how many columns of regions there are
	HorizontalPlanes = MaxX/RegionSize + 1
	//VerticalPlanes Represents how many rows of regions there are
	VerticalPlanes = MaxY/RegionSize + 1
	//LowerBound Represents a dividing line in the exact middle of a region
	LowerBound = RegionSize / 2
)

var EngineChannel = make(chan func(), 20)

//region Represents a 48x48 section of map.  The purpose of this is to keep track of entities in the entire world without having to allocate tiles individually, which would make search algorithms slower and utilizes a great deal of memory.
type region struct {
	Players *entityList
	NPCs    *entityList
	Objects *entityList
	Items   *entityList
}

var regions [HorizontalPlanes][VerticalPlanes]*region

//IsValid Returns true if the tile at x,y is within world boundaries, false otherwise.
func WithinWorld(x, y int) bool {
	return x <= MaxX && x >= 0 && y >= 0 && y <= MaxY
}

//AddPlayer Add a player to the region.
func AddPlayer(p *Player) {
	getRegion(p.X(), p.Y()).Players.Add(p)
}

//RemovePlayer Remove a player from the region.
func RemovePlayer(p *Player) {
	getRegion(p.X(), p.Y()).Players.Remove(p)
}

//AddNpc Add a NPC to the region.
func AddNpc(n *NPC) {
	getRegion(n.X(), n.Y()).NPCs.Add(n)
}

//RemoveNpc Remove a NPC from the region.
func RemoveNpc(n *NPC) {
	getRegion(n.X(), n.Y()).NPCs.Remove(n)
}

//AddItem Add a ground item to the region.
func AddItem(i *GroundItem) {
	getRegion(i.X(), i.Y()).Items.Add(i)
}

//GetItem Returns the item at x,y with the specified id.  Returns nil if it can not find the item.
func GetItem(x, y, id int) *GroundItem {
	region := getRegion(x, y)
	region.Items.lock.RLock()
	defer region.Items.lock.RUnlock()
	for _, i := range region.Items.set {
		if i, ok := i.(*GroundItem); ok {
			if i.ID == id && i.X() == x && i.Y() == y {
				return i
			}
		}
	}

	return nil
}

//RemoveItem Remove a ground item to the region.
func RemoveItem(i *GroundItem) {
	getRegion(i.X(), i.Y()).Items.Remove(i)
}

//AddObject Add an object to the region.
func AddObject(o *Object) {
	getRegion(o.X(), o.Y()).Objects.Add(o)
	if !o.Boundary {
		def := Objects[o.ID]
		if def.Type != 1 && def.Type != 2 {
			return
		}
		var width, height int
		if o.Direction == 0 || o.Direction == 4 {
			width = def.Width
			height = def.Height
		} else {
			width = def.Height
			height = def.Width
		}
		for xOffset := 0; xOffset < width; xOffset++ {
			for yOffset := 0; yOffset < height; yOffset++ {
				x, y := o.X()+xOffset, o.Y()+yOffset
				areaX := (2304 + x) % RegionSize
				areaY := (1776 + y - (944 * ((y + 100) / 944))) % RegionSize
				if sectorFromCoords(x, y) == nil {
					return
				}
				if def.Type == 1 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x40
				} else if o.Direction == 0 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x2
					if sectorFromCoords(x-1, y) != nil && (areaX > 0 || areaY >= 48) {
						sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY].CollisionMask |= 0x8
					}
				} else if o.Direction == 2 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x4
					if sectorFromCoords(x, y+1) != nil {
						sectorFromCoords(x, y+1).Tiles[areaX*RegionSize+areaY+1].CollisionMask |= 0x1
					}
				} else if o.Direction == 4 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x8
					if sectorFromCoords(x+1, y) != nil {
						sectorFromCoords(x+1, y).Tiles[(areaX+1)*RegionSize+areaY].CollisionMask |= 0x2
					}
				} else if o.Direction == 6 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x1
					if sectorFromCoords(x, y-1) != nil {
						sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1].CollisionMask |= 0x4
					}
				}
			}
		}
	} else {
		def := BoundaryDefs[o.ID]
		if def.Traversable != 1 {
			return
		}
		x, y := o.X(), o.Y()
		areaX := (2304 + x) % RegionSize
		areaY := (1776 + y - (944 * ((y + 100) / 944))) % RegionSize
		if sectorFromCoords(x, y) == nil {
			return
		}
		if o.Direction == 0 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x1
			if sectorFromCoords(x, y-1) != nil {
				sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1].CollisionMask |= 0x4
			}
		} else if o.Direction == 1 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x2
			if sectorFromCoords(x-1, y) != nil && (areaX > 0 || areaY >= 48) {
				sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY].CollisionMask |= 0x8
			}
		} else if o.Direction == 2 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x10
		} else if o.Direction == 3 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x20
		}
	}
}

//RemoveObject Remove an object from the region.
func RemoveObject(o *Object) {
	getRegion(o.X(), o.Y()).Objects.Remove(o)
	if !o.Boundary {
		def := Objects[o.ID]
		if def.Type != 1 && def.Type != 2 {
			return
		}
		var width, height int
		if o.Direction == 0 || o.Direction == 4 {
			width = def.Width
			height = def.Height
		} else {
			width = def.Height
			height = def.Width
		}
		for xOffset := 0; xOffset < width; xOffset++ {
			for yOffset := 0; yOffset < height; yOffset++ {
				x, y := o.X()+xOffset, o.Y()+yOffset
				areaX := (2304 + x) % RegionSize
				areaY := (1776 + y - (944 * ((y + 100) / 944))) % RegionSize
				if def.Type == 1 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 0x40
				} else if o.Direction == 0 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 2
					if sectorFromCoords(x-1, y) != nil {
						sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY].CollisionMask &= 0xFFFF - 8
					}
				} else if o.Direction == 2 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 4
					if sectorFromCoords(x, y+1) != nil {
						sectorFromCoords(x, y+1).Tiles[areaX*RegionSize+areaY+1].CollisionMask &= 0xFFFF - 1
					}
				} else if o.Direction == 4 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 8
					if sectorFromCoords(x+1, y) != nil {
						sectorFromCoords(x+1, y).Tiles[(areaX+1)*RegionSize+areaY].CollisionMask &= 0xFFFF - 2
					}
				} else if o.Direction == 6 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 1
					if sectorFromCoords(x, y-1) != nil {
						sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1].CollisionMask &= 0xFFFF - 4
					}
				}
			}
		}
	} else {
		def := BoundaryDefs[o.ID]
		if def.Traversable != 1 {
			return
		}
		x, y := o.X(), o.Y()
		areaX := (2304 + x) % RegionSize
		areaY := (1776 + y - (944 * ((y + 100) / 944))) % RegionSize
		if o.Direction == 0 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 1
			if sectorFromCoords(x, y-1) != nil {
				sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1].CollisionMask &= 0xFFFF - 4
			}
		} else if o.Direction == 1 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 2
			if sectorFromCoords(x-1, y) != nil {
				sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY].CollisionMask &= 0xFFFF - 8
			}
		} else if o.Direction == 2 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 0x10
		} else if o.Direction == 3 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 0x20
		}
	}
}

//ReplaceObject Replaces old with a new game object with all of the same characteristics, except it's ID set to newID.
func ReplaceObject(old *Object, newID int) *Object {
	RemoveObject(old)
	object := NewObject(newID, old.Direction, old.X(), old.Y(), old.Boundary)
	AddObject(object)
	return object
}

//GetAllObjects Returns a slice containing all objects in the game world.
func GetAllObjects() (list []*Object) {
	for x := 0; x < MaxX; x += RegionSize {
		for y := 0; y < MaxY; y += RegionSize {
			if r := regions[x/RegionSize][y/RegionSize]; r != nil {
				r.Objects.lock.RLock()
				for _, o := range r.Objects.set {
					if o, ok := o.(*Object); ok {
						list = append(list, o)
					}
				}
				r.Objects.lock.RUnlock()
			}
		}
	}

	return
}

//GetObject If there is an object at these coordinates, returns it.  Otherwise, returns nil.
func GetObject(x, y int) *Object {
	r := getRegion(x, y)
	r.Objects.lock.RLock()
	defer r.Objects.lock.RUnlock()
	for _, o := range r.Objects.set {
		if o, ok := o.(*Object); ok {
			if o.X() == x && o.Y() == y {
				return o
			}
		}
	}

	return nil
}

//GetNpc Returns the NPC with the specified server index.
func GetNpc(index int) *NPC {
	if index > len(Npcs)-1 {
		log.Info.Printf("Index out of bounds in call to GetNpc.  Length:%d, Requested:%d\n", len(Npcs), index)
		return nil
	}
	return Npcs[index]
}

//getRegionFromIndex internal function to get a region by its row amd column indexes
func getRegionFromIndex(areaX, areaY int) *region {
	if areaX < 0 {
		areaX = 0
	}
	if areaX >= HorizontalPlanes {
		fmt.Println("planeX index out of range")
		return &region{&entityList{}, &entityList{}, &entityList{}, &entityList{}}
	}
	if areaY < 0 {
		areaY = 0
	}
	if areaY >= VerticalPlanes {
		fmt.Println("planeY index out of range")
		return &region{&entityList{}, &entityList{}, &entityList{}, &entityList{}}
	}
	if regions[areaX][areaY] == nil {
		regions[areaX][areaY] = &region{&entityList{}, &entityList{}, &entityList{}, &entityList{}}
	}
	return regions[areaX][areaY]
}

//getRegion Returns the region that corresponds with the given coordinates.  If it does not exist yet, it will allocate a new onr and store it for the lifetime of the application in the regions map.
func getRegion(x, y int) *region {
	return getRegionFromIndex(x/RegionSize, y/RegionSize)
}

//getRegionFromLocation Returns the region that corresponds with the given location.  If it does not exist yet, it will allocate a new onr and store it for the lifetime of the application in the regions map.
func getRegionFromLocation(loc *Location) *region {
	return getRegionFromIndex(loc.X()/RegionSize, loc.Y()/RegionSize)
}

//surroundingRegions Returns the regions surrounding the given coordinates.  It wil
func surroundingRegions(x, y int) (regions [4]*region) {
	areaX := x / RegionSize
	areaY := y / RegionSize
	regions[0] = getRegionFromIndex(areaX, areaY)
	relX := x % RegionSize
	relY := y % RegionSize
	if relX <= LowerBound {
		regions[1] = getRegionFromIndex(areaX-1, areaY)
		if relY <= LowerBound {
			regions[2] = getRegionFromIndex(areaX-1, areaY-1)
			regions[3] = getRegionFromIndex(areaX, areaY-1)
		} else {
			regions[2] = getRegionFromIndex(areaX-1, areaY+1)
			regions[3] = getRegionFromIndex(areaX, areaY+1)
		}
	} else if relY <= LowerBound {
		regions[1] = getRegionFromIndex(areaX+1, areaY)
		regions[2] = getRegionFromIndex(areaX+1, areaY-1)
		regions[3] = getRegionFromIndex(areaX, areaY-1)
	} else {
		regions[1] = getRegionFromIndex(areaX+1, areaY)
		regions[2] = getRegionFromIndex(areaX+1, areaY+1)
		regions[3] = getRegionFromIndex(areaX, areaY+1)
	}

	return
}
