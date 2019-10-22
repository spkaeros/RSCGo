package world

import (
	"fmt"
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
}

var regions [HorizontalPlanes][VerticalPlanes]*Region

//WithinWorld Returns true if the tile at x,y is within world boundaries, false otherwise.
func WithinWorld(x, y int) bool {
	return x <= MaxX && x >= 0 && y >= 0 && y <= MaxY
}

//AddPlayer Add a player to the region.
func AddPlayer(p *Player) {
	GetRegion(int(p.X), int(p.Y)).Players.Add(p)
}

//RemovePlayer Remove a player from the region.
func RemovePlayer(p *Player) {
	GetRegion(int(p.X), int(p.Y)).Players.Remove(p)
}

//AddNpc Add a NPC to the region.
func AddNpc(n *NPC) {
	GetRegion(int(n.X), int(n.Y)).NPCs.Add(n)
}

//RemoveNpc Remove a NPC from the region.
func RemoveNpc(n *NPC) {
	GetRegion(int(n.X), int(n.Y)).NPCs.Remove(n)
}

//AddObject Add an object to the region.
func AddObject(o *Object) {
	GetRegion(int(o.X), int(o.Y)).Objects.Add(o)
}

//RemoveObject Remove an object from the region.
func RemoveObject(o *Object) {
	GetRegion(int(o.X), int(o.Y)).Objects.Remove(o)
}

//ReplaceObject Replaces old with a new game object with all of the same characteristics, except it's ID set to newID.
func ReplaceObject(old *Object, newID int) {
	r := GetRegionFromLocation(&old.Location)
	r.Objects.Remove(old)
	r.Objects.Add(NewObject(newID, old.Direction, int(old.X), int(old.Y), old.Boundary))
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
			if o.X == uint32(x) && o.Y == uint32(y) {
				return o
			}
		}
	}

	return nil
}

//getRegionFromIndex internal function to get a region by its row amd column indexes
func getRegionFromIndex(areaX, areaY int) *Region {
	if areaX < 0 {
		areaX = 0
	}
	if areaX >= HorizontalPlanes {
		fmt.Println("planeX index out of range")
		return &Region{&List{}, &List{}, &List{}}
	}
	if areaY < 0 {
		areaY = 0
	}
	if areaY >= VerticalPlanes {
		fmt.Println("planeY index out of range")
		return &Region{&List{}, &List{}, &List{}}
	}
	if regions[areaX][areaY] == nil {
		regions[areaX][areaY] = &Region{&List{}, &List{}, &List{}}
	}
	return regions[areaX][areaY]
}

//GetRegion Returns the region that corresponds with the given coordinates.  If it does not exist yet, it will allocate a new onr and store it for the lifetime of the application in the regions map.
func GetRegion(x, y int) *Region {
	return getRegionFromIndex(x/RegionSize, y/RegionSize)
}

//GetRegionFromLocation Returns the region that corresponds with the given location.  If it does not exist yet, it will allocate a new onr and store it for the lifetime of the application in the regions map.
func GetRegionFromLocation(loc *Location) *Region {
	return getRegionFromIndex(int(loc.X/RegionSize), int(loc.Y/RegionSize))
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
