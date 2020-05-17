package world

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/spkaeros/rscgo/pkg/config"
	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/log"
	rscRand "github.com/spkaeros/rscgo/pkg/rand"
)

const (
	TicksDay       = 135000
	TicksHour      = 5625
	TicksTwentyMin = 1875
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

//OrderedDirections This is an array containing all of the directions a mob can walk in, ordered by path finder precedent.
var OrderedDirections = [...]int{2, 6, 0, 4, 3, 5, 1, 7}

func (l Location) Step(direction int) Location {
	loc := l.Clone()
	if direction == 0 || direction == 1 || direction == 7 {
		loc.y.Dec()
	} else if direction == 4 || direction == 5 || direction == 3 {
		loc.y.Inc()
	}
	if direction == 1 || direction == 2 || direction == 3 {
		loc.x.Inc()
	} else if direction == 5 || direction == 6 || direction == 7 {
		loc.x.Dec()
	}
	return loc
}

//UpdateTime a point in time in the future to log all active players out and shut down the game for updates.
// Before the command is issued to set this time, it is initialized to time.Time{} zero value.
var UpdateTime time.Time

type PlayerMap struct {
	usernames map[uint64]*Player
	indices   map[int]*Player
	lock      sync.RWMutex
}

//Players Collection containing all of the active client, by index and username hash, guarded by a mutex
var Players = &PlayerMap{usernames: make(map[uint64]*Player), indices: make(map[int]*Player)}

//FromUserHash Returns the client with the base37 username `hash` if it exists and true, otherwise returns nil and false.
func (m *PlayerMap) FromUserHash(hash uint64) (*Player, bool) {
	m.lock.RLock()
	result, ok := m.usernames[hash]
	m.lock.RUnlock()
	return result, ok
}

//ContainsHash Returns true if there is a client mapped to this username hash is in this collection, otherwise returns false.
func (m *PlayerMap) ContainsHash(hash uint64) bool {
	_, ret := m.FromUserHash(hash)
	return ret
}

//FromIndex Returns the client with the index `index` if it exists and true, otherwise returns nil and false.
func (m *PlayerMap) FromIndex(index int) (*Player, bool) {
	m.lock.RLock()
	result, ok := m.indices[index]
	m.lock.RUnlock()
	return result, ok
}

//Add Puts a client into the map.
func (m *PlayerMap) Put(player *Player) {
	nextIndex := m.NextIndex()
	m.lock.Lock()
	player.Index = nextIndex
	m.usernames[player.UsernameHash()] = player
	m.indices[nextIndex] = player
	m.lock.Unlock()
}

//Remove Removes a client from the map.
func (m *PlayerMap) Remove(player *Player) {
	m.lock.Lock()
	delete(m.usernames, player.UsernameHash())
	delete(m.indices, player.Index)
	m.lock.Unlock()
}

//Range Calls action for every active client in the collection.
func (m *PlayerMap) Range(action func(*Player)) {
	m.lock.RLock()
	for _, c := range m.indices {
		if c != nil && c.Connected() {
			action(c)
		}
	}
	m.lock.RUnlock()
}

//Size Returns the size of the active client collection.
func (m *PlayerMap) Size() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.usernames)
}

//NextIndex Returns the lowest available index for the client to be mapped to.
func (m *PlayerMap) NextIndex() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for i := 0; i < config.MaxPlayers(); i++ {
		if _, ok := m.indices[i]; !ok {
			return i
		}
	}
	return -1
}

//region Represents a 48x48 section of map.  The purpose of this is to keep track of entities in the entire world without having to allocate tiles individually, which would make search algorithms slower and utilizes a great deal of memory.
type region struct {
	Players *MobList
	NPCs    *MobList
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
	Players.Put(p)
	Players.Range(func(player *Player) {
		if player.FriendList.Contains(p.Username()) && (!p.FriendBlocked() || p.FriendsWith(player.UsernameHash())) {
			player.FriendList.ToggleStatus(p.Username())
			player.SendPacket(FriendUpdate(p.UsernameHash(), true))
		}
		
//		if player.FriendList.Contains(p.Username()) {
//			player.SendPacket(FriendUpdate(p.UsernameHash(), p.FriendList.Contains(player.Username()) || !p.FriendBlocked()))
//		}
	})
}

//RemovePlayer SetRegionRemoved a player from the region.
func RemovePlayer(p *Player) {
	//	p.UpdateStatus(false)
	getRegion(p.X(), p.Y()).Players.Remove(p)
	Players.Remove(p)
	Players.Range(func(player *Player) {
		if player.FriendList.Contains(p.Username()) && (!p.FriendBlocked() || p.FriendsWith(player.UsernameHash())) {
			player.FriendList.ToggleStatus(p.Username())
			player.SendPacket(FriendUpdate(p.UsernameHash(), false))
		}
	})
	p.SetRegionRemoved()
}

//AddNpc Add a NPC to the region.
func AddNpc(n *NPC) {
	getRegion(n.X(), n.Y()).NPCs.Add(n)
}

//RemoveNpc SetRegionRemoved a NPC from the region.
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

//RemoveItem SetRegionRemoved a ground item to the region.
func RemoveItem(i *GroundItem) {
	getRegion(i.X(), i.Y()).Items.Remove(i)
}

//AddObject Add an object to the region.
func AddObject(o *Object) {
	getRegion(o.X(), o.Y()).Objects.Add(o)
	if !o.Boundary {
		scenary := definitions.ScenaryObjects[o.ID]
		// type 0 is used when the object causes no collisions of any sort.
		// type 1 is used when the object fully blocks the tile(s) that it sits on.  Marks tile as fully blocked.
		// type 2 is used when the object mimics a boundary, e.g for gates and the like.
		// type 3 is used when the object mimics an opened door-type boundary, e.g opened gates and the like.
		if scenary.CollisionType%3 == 0 {
			return
		}
		width := scenary.Height
		height := scenary.Width
		//if o.Direction == 0 || o.Direction == 4 {
		if o.Direction%4 == 0 {
			width = scenary.Width
			height = scenary.Height
		}
		for x := o.X(); x < o.X()+width; x++ {
			for y := o.Y(); y < o.Y()+height; y++ {
				areaX := (2304 + x) % RegionSize
				areaY := (1776 + y - (944 * ((y + 100) / 944))) % RegionSize
				if len(sectorFromCoords(x, y).Tiles) <= 0 {
					log.Warning.Println("ERROR: Sector with no tiles at:" + strconv.Itoa(x) + "," + strconv.Itoa(y) + " (" + strconv.Itoa(areaX) + "," + strconv.Itoa(areaY) + "\n")
					return
				}
				if scenary.CollisionType == 1 {
					// Blocks the whole tile.  Can not walk on it from any direction
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipFullBlock
					continue
				}

				// If it's gone this far, collisionType is 2 (directional blocking, e.g gates etc)
				if o.Direction == byte(North) {
					// Block the tiles east side
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipEast
					// ensure that the neighbors index is valid
					if len(sectorFromCoords(x-1, y).Tiles) > 0 && (areaX > 0 || areaY >= RegionSize) {
						// then block the eastern neighbors west side
						sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY].CollisionMask |= ClipWest
					}
				} else if o.Direction == byte(West) {
					// Block the tiles south side
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipSouth
					// then block the southern neighbors north side
					sectorFromCoords(x, y+1).Tiles[areaX*RegionSize+areaY+1].CollisionMask |= ClipNorth
				} else if o.Direction == byte(South) {
					// Block the tiles west side
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipWest
					// then block the western neighbors east side
					if areaX, areaY := (2304+x+1)%RegionSize, (1776+y-(944*((y+100)/944)))%RegionSize; (areaX+1)*RegionSize+areaY > 2304 {
						sectorFromCoords(x+1, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipEast
					}
				} else if o.Direction == byte(East) {
					// Block the tiles north side
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipNorth
					// ensure that the neighbors index is valid
					if len(sectorFromCoords(x, y-1).Tiles) > 0 && areaX+areaY > 0 {
						// then block the eastern neighbors west side
						sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1].CollisionMask |= ClipSouth
					}
				}

			}
		}
	} else {
		boundary := definitions.BoundaryObjects[o.ID]
		if !boundary.Solid {
			// Doorframes and some other stuff collide with nothing.
			return
		}
		x, y := o.X(), o.Y()
		areaX := (2304 + x) % RegionSize
		areaY := (1776 + y - (944 * ((y + 100) / 944))) % RegionSize
		if len(sectorFromCoords(x, y).Tiles) <= 0 {
			log.Warn("ERROR: Sector with no tiles at:" + strconv.Itoa(x) + "," + strconv.Itoa(y) + " (" + strconv.Itoa(areaX) + "," + strconv.Itoa(areaY) + "\n")
			return
		}
		if o.Direction == 0 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipNorth
			if areaX+areaY > 0 {
				sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1].CollisionMask |= ClipSouth
			}
		} else if o.Direction == 1 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipEast
			if areaX > 0 || areaY >= 48 {
				sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY].CollisionMask |= ClipWest
			}
		} else if o.Direction == 2 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipSwNe
		} else if o.Direction == 3 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipSeNw
		}
	}
}

//RemoveObject SetRegionRemoved an object from the region.
func RemoveObject(o *Object) {
	getRegion(o.X(), o.Y()).Objects.Remove(o)
	if !o.Boundary {
		scenary := definitions.ScenaryObjects[o.ID]
		// type 0 is used when the object causes no collisions of any sort.
		// type 1 is used when the object fully blocks the tile(s) that it sits on.  Marks tile as fully blocked.
		// type 2 is used when the object mimics a boundary, e.g for gates and the like.
		// type 3 is used when the object mimics an opened door-type boundary, e.g opened gates and the like.
		if scenary.CollisionType%3 == 0 {
			return
		}
		width := scenary.Height
		height := scenary.Width

		//if o.Direction == byte(North) || o.Direction == byte(South) {
		if o.Direction%4 == 0 {
			// reverse measurements for directions 0(North) and 4(South), as scenary measurements
			// are oriented vertically by default
			width = scenary.Width
			height = scenary.Height
		}
		for x := o.X(); x < o.X()+width; x++ {
			for y := o.Y(); y < o.Y()+height; y++ {
				areaX := (2304 + x) % RegionSize
				areaY := (1776 + y - (944 * ((y + 100) / 944))) % RegionSize
				if scenary.CollisionType == 1 {
					// This indicates a solid object.  Impassable and blocks the whole tile.
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= ^ClipFullBlock
				} else if o.Direction == 0 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= ^ClipEast
					if sectorFromCoords(x-1, y) != nil {
						sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY].CollisionMask &= ^ClipWest
					}
				} else if o.Direction == 2 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= ^ClipSouth
					if sectorFromCoords(x, y+1) != nil {
						sectorFromCoords(x, y+1).Tiles[areaX*RegionSize+areaY+1].CollisionMask &= ^ClipNorth
					}
				} else if o.Direction == 4 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= ^ClipWest
					if sectorFromCoords(x+1, y) != nil {
						sectorFromCoords(x+1, y).Tiles[(areaX+1)*RegionSize+areaY].CollisionMask &= ^ClipEast
					}
				} else if o.Direction == 6 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= ^ClipNorth
					if sectorFromCoords(x, y-1) != nil {
						sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1].CollisionMask &= ^ClipSouth
					}
				}
			}
		}
	} else {
		// Wall or door location
		boundary := definitions.BoundaryObjects[o.ID]
		if !boundary.Solid {
			return
		}

		x, y := o.X(), o.Y()
		areaX := (2304 + x) % RegionSize
		areaY := (1776 + y - (944 * ((y + 100) / 944))) % RegionSize
		if o.Direction == 0 { // Vertical wall ('| ',' |') North<->South
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= ^ClipNorth
			if areaX+areaY > 0 {
				sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1].CollisionMask &= ^ClipSouth
			}
		} else if o.Direction == 1 { // Horizontal wall ('__','‾‾') East<->West
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= ^ClipEast
			if areaX > 0 || areaY >= 48 {
				sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY].CollisionMask &= ^ClipWest
			}
		} else if o.Direction == 2 { // Diagonal wall ('\','‾|','|_') Southwest<->Northeast
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= ^ClipSwNe
		} else if o.Direction == 3 { // Diagonal wall ('/','|‾','_|') Southeast<->Northwest
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= ^ClipSeNw
		}
	}
}

//ReplaceObject Replaces old with a new game object with all of the same characteristics, except it's ID set to newID.
func ReplaceObject(old *Object, newID int) *Object {
	RemoveObject(old)
	object := NewObject(newID, int(old.Direction), old.X(), old.Y(), old.Boundary)
	AddObject(object)
	return object
}

//GetAllObjects Returns a slice containing all objects in the game
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

//GetNpc Returns the NPC with the specified game index.
func GetNpc(index int) *NPC {
	m := Npcs.Get(index)
	if m == nil {
		return nil
	}
	return m.(*NPC)
}

//NpcNearest looks for the NPC with the given ID, that is the closest to the given coordinates
// and then returns it.
// Returns nil if it can not find an NPC to fit the given description.
func NpcNearest(id, x, y int) *NPC {
	point := NewLocation(x, y)
	minDelta := 16
	var npc *NPC
	for x := 0; x < MaxX; x += RegionSize {
		for y := 0; y < MaxY; y += RegionSize {
			if r := regions[x/RegionSize][y/RegionSize]; r != nil {
				r.NPCs.RangeNpcs(func(n *NPC) bool {
					if n.ID == id && n.LongestDelta(point) < minDelta {
						minDelta = n.LongestDelta(point)
						npc = n
					}
					return false
				})
			}
		}

	}
	return npc
}

//NpcVisibleFrom looks for any NPC with the given ID, that is within a 16x16 square
// surrounding the given coordinates, and then returns it.
// Returns nil if it can not find an NPC to fit the given parameters.
func NpcVisibleFrom(id, x, y int) *NPC {
	point := NewLocation(x, y)
	minDelta := 16
	var npc *NPC
	for x := 0; x < MaxX; x += RegionSize {
		for y := 0; y < MaxY; y += RegionSize {
			if r := regions[x/RegionSize][y/RegionSize]; r != nil {
				r.NPCs.RangeNpcs(func(n *NPC) bool {
					if n.ID == id && n.LongestDelta(point) < minDelta {
						minDelta = n.LongestDelta(point)
						npc = n
					}
					return false
				})
			}
		}

	}
	return npc
}

//getRegionFromIndex internal function to get a region by its row amd column indexes
func getRegionFromIndex(areaX, areaY int) *region {
	if areaX < 0 {
		areaX = 0
	}
	if areaX >= HorizontalPlanes {
		fmt.Println("planeX index out of range")
		return &region{&MobList{}, &MobList{}, &entityList{}, &entityList{}}
	}
	if areaY < 0 {
		areaY = 0
	}
	if areaY >= VerticalPlanes {
		fmt.Println("planeY index out of range")
		return &region{&MobList{}, &MobList{}, &entityList{}, &entityList{}}
	}
	if regions[areaX][areaY] == nil {
		regions[areaX][areaY] = &region{&MobList{}, &MobList{}, &entityList{}, &entityList{}}
	}
	return regions[areaX][areaY]
}

//getRegion Returns the region that corresponds with the given coordinates.  If it does not exist yet, it will allocate a new onr and store it for the lifetime of the application in the regions map.
func getRegion(x, y int) *region {
	return getRegionFromIndex(x/RegionSize, y/RegionSize)
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

//BoundedChance is a statistically random function that should return true percent/maxPercent% of the time.
// This should return true approximately percent/maxPercent of the time, and false (maxPercent-percent)/maxPercent of
// the time.
// percent defines the percentage of chance for this check to pass, clamped to the provided boundaries.
// minPercent is the minimum allowed value of percent.  If percent is larger than minPercent, then minPercent is used
// in its place.
// maxPercent is the maximum allowed value of percent, and also the denominator used when scaling the percentage to
// a 1-byte value
// Returns: randByte <= max(min(percent, maxPercent), minPercent)/maxPercent*256.0, (with uniform randomness)
func BoundedChance(percent float64, minPercent, maxPercent float64) bool {
	if minPercent < 0.0 {
		minPercent = 0.0
	}
//	if maxPercent > 100.0 {
//		maxPercent = 100.0
//	}
	if minPercent > maxPercent {
		maxPercent, minPercent = minPercent, maxPercent
	}
	return rscRand.Rng.Float64() <= math.Max(minPercent, math.Min(maxPercent, percent))/maxPercent
}

//Chance should return true (percent)% of the time, and false (100-percent)% of the time.
// It uses ISAAC64+ to provide randomness.
//
// percent defines the percentage of chance for this check to pass.
func Chance(percent float64) bool {
	return BoundedChance(percent, 0.0, 100.0)
}

//WeightedChoice Awesome API call takes map[retVal]probability as input and returns a statistically weighted randomized retVal as output.
//
// The input's mapped value assigned to each key is its return probability, out of the total sum of all return probabilities.
// You can determine the percentage chance of any given input entry being returned by: probability/sumOfAllProbabilities*100
// E.g, if the sum of all probabilities is 100, and you have a total probability of 100, where the first retVal maps to 25.0, the chance it will be returned is 25%
//
// You can make the total anything. Useful for anything that needs to return certain values deterministically more often than others, but randomly.
func WeightedChoice(choices map[int]float64) int {
	total := 0.0
	totalProb := 0.0
	for _, probability := range choices {
		totalProb += probability
	}

	// We determine the real upper limit of the cumulative probability.
	// if the sum of all probabilities adds up to 100.0% or less, then the upper limit will be 65535
	// if it's less, any unaccounted for percentage will return -1.
	// if it's more, it'll scale 65535 up by however much percentage it needs to handle all inputs provided, and subsequently
	// all individual probabilities must be scaled down by the same figure to account for the difference.
	// e.g passing 3 values in with 50.0 probability on each of them will result in each entry returning at a
	// rate of 33.333~% instead of 50% each, because 50/150=33.333~
//	cumulativeProbability := math.Max(1.0, totalProb/100)
//	upperBound := cumulativeProbability * math.MaxUint16
//	hit := float64(rscRand.Int31N(1, int(upperBound)))
	hit := rscRand.Rng.Float64()*totalProb
	if config.Verbosity >= 3 {
//		log.Debugf("WeightedChoice: Upper bound for total probability:%d {\n", int(totalProb))
		log.Debugf("\nRolled: %d/%d;\n", int(hit), int(totalProb))
//		defer log.Info.Println("};")
	}
	for choice, prob := range choices {
//		newProb := prob / 100 * math.MaxUint16 / upperBound
		if config.Verbosity >= 3 {
			log.Debugf("\tentry{val:%d; hit range between %d - %d (%.2f%% chance)}\n", choice, int(total), int(total+(prob)), prob/totalProb*100)
		}
		total += prob
		if hit < total {
			log.Debugf("Hit:%v\n", choice)
			return choice
		}
	}
	if config.Verbosity >= 3 {
		log.Debugln("Rolled value did not return anything!")
	}
	return -1
}
