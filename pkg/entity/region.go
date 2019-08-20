package entity

const (
	RegionSize       = 64
	HorizontalPlanes = MaxX/RegionSize + 1
	VerticalPlanes   = MaxY/RegionSize + 1
	LowerBound       = RegionSize / 2
)

//Region Represents a 64x64 section of map.  The purpose of this is to keep track of entities in the entire world without having to allocate tiles individually, which would make search algorithms slower and utilizes a great deal of memory.
type Region struct {
	Players map[int]*Player
}

var regions [HorizontalPlanes][VerticalPlanes]*Region

//AddPlayer Add a player to the region.
func (r *Region) AddPlayer(p *Player) {
	r.Players[p.Index] = p
}

//RemovePlayer Remove a player from the region.
func (r *Region) RemovePlayer(p *Player) {
	delete(r.Players, p.Index)
}

//getRegionFromIndex internal function to get a region by its area coordinate indexes
func getRegionFromIndex(areaX, areaY int) *Region {
	if regions[areaX][areaY] == nil {
		regions[areaX][areaY] = &Region{Players: make(map[int]*Player)}
	}
	return regions[areaX][areaY]
}

//GetRegion Returns the region that corresponds with the given coordinates.  If it does not exist yet, it will allocate a new onr and store it for the lifetime of the application in the regions map.
func GetRegion(x, y int) *Region {
	return getRegionFromIndex(x/RegionSize, y/RegionSize)
}

//GetRegionFromLocation Returns the region that corresponds with the given location.  If it does not exist yet, it will allocate a new onr and store it for the lifetime of the application in the regions map.
func GetRegionFromLocation(loc *Location) *Region {
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
