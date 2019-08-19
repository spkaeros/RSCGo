package entity

import (
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
)

type Region struct {
	Players map[uint64]*Player
}

var regions [20][79]*Region

func (r *Region) AddPlayer(p *Player) {
	r.Players[strutil.Base37(p.Username)] = p
}

func (r *Region) RemovePlayer(p *Player) {
	delete(r.Players, strutil.Base37(p.Username))
}

func regionGet(x, y int) *Region {
	if regions[x][y] == nil {
		regions[x][y] = &Region{Players: make(map[uint64]*Player)}
	}
	return regions[x][y]
}

func GetRegion(x, y int) *Region {
	if regions[x/48][y/48] == nil {
		regions[x/48][y/48] = &Region{Players: make(map[uint64]*Player)}
	}
	return regions[x/48][y/48]
}

func SurroundingRegions(x, y int) (regions []*Region) {
	areaX := x / 48
	areaY := y / 48
	regions = append(regions, regionGet(areaX, areaY))
	relX := x % 48
	relY := y % 48
	if relX <= 24 {
		if relY <= 24 {
			regions = append(regions, regionGet(areaX-1, areaY))
			regions = append(regions, regionGet(areaX-1, areaY-1))
			regions = append(regions, regionGet(areaX, areaY-1))
		} else {
			regions = append(regions, regionGet(areaX-1, areaY))
			regions = append(regions, regionGet(areaX-1, areaY+1))
			regions = append(regions, regionGet(areaX, areaY+1))
		}
	} else if relY <= 24 {
		regions = append(regions, regionGet(areaX+1, areaY))
		regions = append(regions, regionGet(areaX+1, areaY-1))
		regions = append(regions, regionGet(areaX, areaY-1))
	} else {
		regions = append(regions, regionGet(areaX+1, areaY))
		regions = append(regions, regionGet(areaX+1, areaY+1))
		regions = append(regions, regionGet(areaX, areaY+1))
	}

	return
}
