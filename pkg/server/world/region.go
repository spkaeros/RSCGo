package world

import (
	"fmt"
	"time"
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
	Objects *List
}

var regions [HorizontalPlanes][VerticalPlanes]*Region

//WithinWorld Returns true if the tile at x,y is within world boundaries, false otherwise.
func WithinWorld(x, y int) bool {
	return x <= MaxX && x >= 0 && y >= 0 && y <= MaxY
}

//AddPlayer Add a player to the region.
func AddPlayer(p *Player) {
	GetRegion(p.X, p.Y).Players.Add(p)
}

//RemovePlayer Remove a player from the region.
func RemovePlayer(p *Player) {
	GetRegion(p.X, p.Y).Players.Remove(p)
}

//AddObject Add an object to the region.
func AddObject(o *Object) {
	GetRegion(o.X, o.Y).Objects.Add(o)
}

//RemoveObject Remove an object from the region.
func RemoveObject(o *Object) {
	GetRegion(o.X, o.Y).Objects.Remove(o)
}

//EnterDoor Replaces door object with an open door, sleeps for one second, and returns the closed door.
func (player *Player) EnterDoor(door *Object, dest Location) {
	ReplaceObject(door, 11)
	player.SetLocation(dest)
	time.Sleep(time.Second)
	ReplaceObject(GetObject(door.X, door.Y), door.ID)
}

//ReplaceObject Replaces old with a new game object with all of the same characteristics, except it's ID set to newID.
func ReplaceObject(old *Object, newID int) {
	r := GetRegionFromLocation(old.Location)
	r.Objects.Remove(old)
	r.Objects.Add(NewObject(newID, old.Direction, old.X, old.Y, old.Boundary))
}

//GetAllObjects Returns a slice containing all objects in the game world.
func GetAllObjects() (list []*Object) {
	for x := 0; x < MaxX; x += RegionSize {
		for y := 0; y < MaxY; y += RegionSize {
			if r := regions[x/RegionSize][y/RegionSize]; r != nil {
				r.Objects.lock.RLock()
				defer r.Objects.lock.RUnlock()
				for _, o := range r.Objects.List {
					if o, ok := o.(*Object); ok {
						list = append(list, o)
					}
				}
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
			if o.X == x && o.Y == y {
				return o
			}
		}
	}

	return nil
}

//getRegionFromIndex internal function to get a region by its row amd column indexes
func getRegionFromIndex(areaX, areaY int) *Region {
	if areaX < 0 || areaX >= HorizontalPlanes {
		fmt.Println("planeX index out of range")
		return &Region{}
	}
	if areaY < 0 || areaY >= VerticalPlanes {
		fmt.Println("planeY index out of range")
		return &Region{}
	}
	if regions[areaX][areaY] == nil {
		regions[areaX][areaY] = &Region{&List{}, &List{}}
	}
	return regions[areaX][areaY]
}

//GetRegion Returns the region that corresponds with the given coordinates.  If it does not exist yet, it will allocate a new onr and store it for the lifetime of the application in the regions map.
func GetRegion(x, y int) *Region {
	return getRegionFromIndex(x/RegionSize, y/RegionSize)
}

//GetRegionFromLocation Returns the region that corresponds with the given location.  If it does not exist yet, it will allocate a new onr and store it for the lifetime of the application in the regions map.
func GetRegionFromLocation(loc Location) *Region {
	return getRegionFromIndex(loc.X/RegionSize, loc.Y/RegionSize)
}

//SurroundingRegions Returns the regions surrounding the given coordinates.  It wil
func SurroundingRegions(x, y int) (regions [4]*Region) {
	areaX := x / RegionSize
	areaY := y / RegionSize
	regions[0] = getRegionFromIndex(areaX, areaY)
	relX := x % RegionSize
	relY := y % RegionSize
	if relX <= LowerBound {
		if relY <= LowerBound {
			regions[1] = getRegionFromIndex(areaX-1, areaY)
			regions[2] = getRegionFromIndex(areaX-1, areaY-1)
			regions[3] = getRegionFromIndex(areaX, areaY-1)
		} else {
			regions[1] = getRegionFromIndex(areaX-1, areaY)
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
