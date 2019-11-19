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

//Region Represents a 48x48 section of map.  The purpose of this is to keep track of entities in the entire world without having to allocate tiles individually, which would make search algorithms slower and utilizes a great deal of memory.
type Region struct {
	Players *List
	NPCs    *List
	Objects *List
	Items   *List
}

var regions [HorizontalPlanes][VerticalPlanes]*Region

//IsValid Returns true if the tile at x,y is within world boundaries, false otherwise.
func WithinWorld(x, y int) bool {
	return x <= MaxX && x >= 0 && y >= 0 && y <= MaxY
}

//AddPlayer Add a player to the region.
func AddPlayer(p *Player) {
	GetRegion(int(p.X.Load()), int(p.Y.Load())).Players.Add(p)
}

//RemovePlayer Remove a player from the region.
func RemovePlayer(p *Player) {
	GetRegion(int(p.X.Load()), int(p.Y.Load())).Players.Remove(p)
}

//AddNpc Add a NPC to the region.
func AddNpc(n *NPC) {
	GetRegion(int(n.X.Load()), int(n.Y.Load())).NPCs.Add(n)
}

//RemoveNpc Remove a NPC from the region.
func RemoveNpc(n *NPC) {
	GetRegion(int(n.X.Load()), int(n.Y.Load())).NPCs.Remove(n)
}

//AddItem Add a ground item to the region.
func AddItem(i *GroundItem) {
	GetRegion(int(i.X.Load()), int(i.Y.Load())).Items.Add(i)
}

//GetItem Returns the item at x,y with the specified id.  Returns nil if it can not find the item.
func GetItem(x, y, id int) *GroundItem {
	region := GetRegion(x, y)
	region.Items.lock.RLock()
	defer region.Items.lock.RUnlock()
	for _, i := range region.Items.List {
		if i, ok := i.(*GroundItem); ok {
			if i.ID == id && int(i.X.Load()) == x && int(i.Y.Load()) == y {
				return i
			}
		}
	}

	return nil
}

//RemoveItem Remove a ground item to the region.
func RemoveItem(i *GroundItem) {
	GetRegion(int(i.X.Load()), int(i.Y.Load())).Items.Remove(i)
}

//AddObject Add an object to the region.
func AddObject(o *Object) {
/*	x, y := int(o.X.Load()), int(o.Y.Load())
	areaX := (2304+x) % 48
	areaY := (1776+y-(944*((y+100)/944))) % 48
	if Sectors[strutil.JagHash(SectorName(x, y))] != nil {
		Sectors[strutil.JagHash(SectorName(x, y))].Tiles[areaX*48+areaY].CollisionMask |= 64
	}*/
	GetRegion(int(o.X.Load()), int(o.Y.Load())).Objects.Add(o)
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
				x, y := int(o.X.Load())+xOffset, int(o.Y.Load())+yOffset
				areaX := (2304+x) % RegionSize
				areaY := (1776+y-(944*((y+100)/944))) % RegionSize
				if SectorFromCoords(x, y) == nil {
					return
				}
				if def.Type == 1 {
					SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x40
				} else if o.Direction == 0 {
					SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x2
					if SectorFromCoords(x - 1, y) != nil && (areaX > 0 || areaY >= 48) {
						SectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY].CollisionMask |= 0x8
					}
				} else if o.Direction == 2 {
					SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x4
					if SectorFromCoords(x, y + 1) != nil {
						SectorFromCoords(x, y + 1).Tiles[areaX*RegionSize+areaY + 1].CollisionMask |= 0x1
					}
				} else if o.Direction == 4 {
					SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x8
					if SectorFromCoords(x + 1, y) != nil {
						SectorFromCoords(x+1, y).Tiles[(areaX+1)*RegionSize+areaY].CollisionMask |= 0x2
					}
				} else if o.Direction == 6 {
					SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x1
					if SectorFromCoords(x, y - 1) != nil {
						SectorFromCoords(x, y - 1).Tiles[areaX*RegionSize+areaY - 1].CollisionMask |= 0x4
					}
				}
			}
		}
	} else {
		def := Boundarys[o.ID]
		if def.Traversable != 1 {
			return
		}
		x, y := int(o.X.Load()), int(o.Y.Load())
		areaX := (2304+x) % RegionSize
		areaY := (1776+y-(944*((y+100)/944))) % RegionSize
		if SectorFromCoords(x, y) == nil {
			return
		}
		if o.Direction == 0 {
			SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x1
			if SectorFromCoords(x, y - 1) != nil {
				SectorFromCoords(x, y - 1).Tiles[areaX*RegionSize+areaY-1].CollisionMask |= 0x4
			}
		} else if o.Direction == 1 {
			SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x2
			if SectorFromCoords(x - 1, y) != nil && (areaX > 0 || areaY >= 48) {
				SectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY].CollisionMask |= 0x8
			}
		} else if o.Direction == 2 {
			SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x10
		} else if o.Direction == 3 {
			SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= 0x20
		}
	}
}

//RemoveObject Remove an object from the region.
func RemoveObject(o *Object) {
	GetRegion(int(o.X.Load()), int(o.Y.Load())).Objects.Remove(o)
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
				x, y := int(o.X.Load())+xOffset, int(o.Y.Load())+yOffset
				areaX := (2304+x) % RegionSize
				areaY := (1776+y-(944*((y+100)/944))) % RegionSize
				if def.Type == 1 {
					SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 0x40
				} else if o.Direction == 0 {
					SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 2
					if SectorFromCoords(x - 1, y) != nil {
						SectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY].CollisionMask &= 0xFFFF - 8
					}
				} else if o.Direction == 2 {
					SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 4
					if SectorFromCoords(x, y + 1) != nil {
						SectorFromCoords(x, y + 1).Tiles[areaX*RegionSize+areaY + 1].CollisionMask &= 0xFFFF - 1
					}
				} else if o.Direction == 4 {
					SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 8
					if SectorFromCoords(x + 1, y) != nil {
						SectorFromCoords(x+1, y).Tiles[(areaX+1)*RegionSize+areaY].CollisionMask &= 0xFFFF - 2
					}
				} else if o.Direction == 6 {
					SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 1
					if SectorFromCoords(x, y - 1) != nil {
						SectorFromCoords(x, y - 1).Tiles[areaX*RegionSize+areaY - 1].CollisionMask &= 0xFFFF - 4
					}
				}
			}
		}
	} else {
		def := Boundarys[o.ID]
		if def.Traversable != 1 {
			return
		}
		x, y := int(o.X.Load()), int(o.Y.Load())
		areaX := (2304+x) % RegionSize
		areaY := (1776+y-(944*((y+100)/944))) % RegionSize
		if o.Direction == 0 {
			SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 1
			if SectorFromCoords(x, y - 1) != nil {
				SectorFromCoords(x, y - 1).Tiles[areaX*RegionSize+areaY-1].CollisionMask &= 0xFFFF - 4
			}
		} else if o.Direction == 1 {
			SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 2
			if SectorFromCoords(x - 1, y) != nil {
				SectorFromCoords(x - 1, y).Tiles[(areaX-1)*RegionSize+areaY].CollisionMask &= 0xFFFF - 8
			}
		} else if o.Direction == 2 {
			SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 0x10
		} else if o.Direction == 3 {
			SectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= 0xFFFF - 0x20
		}
	}
}

//ReplaceObject Replaces old with a new game object with all of the same characteristics, except it's ID set to newID.
func ReplaceObject(old *Object, newID int) *Object {
	RemoveObject(old)
	object := NewObject(newID, old.Direction, int(old.X.Load()), int(old.Y.Load()), old.Boundary)
	AddObject(object)
	return object
}

//GetAllObjects Returns a slice containing all objects in the game world.
func GetAllObjects() (list []*Object) {
	for x := 0; x < MaxX; x += RegionSize {
		for y := 0; y < MaxY; y += RegionSize {
			if r := regions[x/RegionSize][y/RegionSize]; r != nil {
				r.Objects.lock.RLock()
				for _, o := range r.Objects.List {
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
	r := GetRegion(x, y)
	r.Objects.lock.RLock()
	defer r.Objects.lock.RUnlock()
	for _, o := range r.Objects.List {
		if o, ok := o.(*Object); ok {
			if o.X.Load() == uint32(x) && o.Y.Load() == uint32(y) {
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
func getRegionFromIndex(areaX, areaY int) *Region {
	if areaX < 0 {
		areaX = 0
	}
	if areaX >= HorizontalPlanes {
		fmt.Println("planeX index out of range")
		return &Region{&List{}, &List{}, &List{}, &List{}}
	}
	if areaY < 0 {
		areaY = 0
	}
	if areaY >= VerticalPlanes {
		fmt.Println("planeY index out of range")
		return &Region{&List{}, &List{}, &List{}, &List{}}
	}
	if regions[areaX][areaY] == nil {
		regions[areaX][areaY] = &Region{&List{}, &List{}, &List{}, &List{}}
	}
	return regions[areaX][areaY]
}

//GetRegion Returns the region that corresponds with the given coordinates.  If it does not exist yet, it will allocate a new onr and store it for the lifetime of the application in the regions map.
func GetRegion(x, y int) *Region {
	return getRegionFromIndex(x/RegionSize, y/RegionSize)
}

//GetRegionFromLocation Returns the region that corresponds with the given location.  If it does not exist yet, it will allocate a new onr and store it for the lifetime of the application in the regions map.
func GetRegionFromLocation(loc *Location) *Region {
	return getRegionFromIndex(int(loc.X.Load()/RegionSize), int(loc.Y.Load()/RegionSize))
}

//SurroundingRegions Returns the regions surrounding the given coordinates.  It wil
func SurroundingRegions(x, y int) (regions [4]*Region) {
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
