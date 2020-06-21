package world

import (
	"fmt"
	"math"
	"strconv"
	"math/rand"
	"sync"
	"time"

	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/log"
	rscRand "github.com/spkaeros/rscgo/pkg/rand"
	"github.com/spkaeros/rscgo/pkg/isaac"

	"go.uber.org/atomic"
)

const (
	TicksDay       = 135000
	TicksHour      = 5625
	TicksTwentyMin = 1875
	TicksMinute    = 100

	//MaxX Width of the game
	MaxX = 944
	//MaxY Height of the game
	MaxY = 3776
)

var (
	Ticks = atomic.NewUint64(0)
)


func CurrentTick() int {
	return int(Ticks.Load())
}

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

//UpdateTime a point in time in the future to log all active players out and shut down the game for updates.
// Before the command is issued to set this time, it is initialized to time.Time{} zero value.
var UpdateTime time.Time

type PlayerList struct {
	players [1250]*Player
	curIdx int
	free []int
	sync.RWMutex
}

//Players Collection containing all of the active client, by index and username hash, guarded by a mutex
var Players = &PlayerList{free: make([]int, 0, 1250)}

//FindHash Returns the client with the base37 username `hash` if it exists and true, otherwise returns nil and false.
func (m *PlayerList) FindHash(hash uint64) (*Player, bool) {
	m.RLock()
	defer m.RUnlock()
	for _, p := range m.players {
		if p.UsernameHash() == hash {
			return p, true
		}
	}
	return nil, false
}

//FromIndex Returns the client with the index `index` if it exists and true, otherwise returns nil and false.
func (m *PlayerList) FindIndex(index int) (*Player, bool) {
	m.RLock()
	defer m.RUnlock()
	return m.players[index], m.players[index] != nil
}

//Find Returns the slot that this player occupies in the set.
func (m *PlayerList) Find(player *Player) int {
	m.RLock()
	defer m.RUnlock()
	for i, v := range m.players {
		if v == player {
			return i
		}
	}

	return -1
}

//Contains Returns true if this player is assigned to a slot in the set, otherwise returns false.
func (m *PlayerList) Contains(player *Player) bool {
	return m.Find(player) > -1
}

//ContainsHash Returns true if there is a client mapped to this username hash is in this collection, otherwise returns false.
func (m *PlayerList) ContainsHash(hash uint64) bool {
	player, ret := m.FindHash(hash)
	return ret && player != nil
}


//Put Finds the lowest available empty slot in the list, and puts the player there.
// This will also assign the players server index variable (*Player.Index) to the assigned slot.
func (m *PlayerList) Put(player *Player) {
	if i := m.Find(player); i > -1 {
		log.Debug("Player list double-put attempted; old index is", i)
		return
	}
	player.Index = m.nextSlot()
	m.Lock()
	defer m.Unlock()
	m.players[player.Index] = player
}

//Remove Removes a client from the set.
func (m *PlayerList) Remove(player *Player) {
	freedSlot := player.ServerIndex()
	log.Debug("free slot: ",freedSlot)
	m.Lock()
	defer m.Unlock()
	m.players[freedSlot].Index = -1
	m.players[freedSlot] = nil
	m.free = append(m.free, freedSlot)
	
}

//Range Calls action for every active client in the collection.
func (m *PlayerList) Range(action func(*Player)) {
	m.RLock()
	defer m.RUnlock()
	for _, p := range m.players {
		if p != nil {
			action(p)
		}
	}
}

//Size Returns the size of the active client collection.
func (m *PlayerList) Size() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.players)
}

//NextIndex Returns the lowest available index for the client to be mapped to.
func (m *PlayerList) nextSlot() int {
	m.Lock()
	defer m.Unlock()
	if len(m.free) > 0 {
		idx := m.free[0]
		if len(m.free) > 1 {
			m.free = m.free[1:]
		} else {
			m.free = []int{}
		}
		return idx
	}
	defer func() { m.curIdx += 1 }()
	return m.curIdx
}

func (m *PlayerList) AsyncRange(fn func(*Player)) {
	w := sync.WaitGroup{}
	m.RLock()
	defer m.RUnlock()
	for _, p := range m.players {
		if p == nil {
			continue
		}
		if p != nil {
			w.Add(1)
			go func() {
				fn(p)
				w.Done()
			}()
		}
	}
	w.Wait()
}

//region Represents a 48x48 section of map.  The purpose of this is to keep track of entities in the entire world without having to allocate tiles individually, which would make search algorithms slower and utilizes a great deal of memory.
type region struct {
	x int
	y int
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

//AddPlayer Add a player to a region of the game world.
func AddPlayer(p *Player) {
	Region(p.X(), p.Y()).Players.Add(p)
	Players.Put(p)
	Players.Range(func(player *Player) {
		if player.FriendList.Contains(p.Username()) && (!p.FriendBlocked() || p.FriendList.Contains(player.Username())) {
			player.FriendList.Set(p.Username(), true)
			player.SendPacket(FriendUpdate(p.UsernameHash(), true))
		}
		
//		if player.FriendList.Contains(p.Username()) {
//			player.SendPacket(FriendUpdate(p.UsernameHash(), p.FriendList.Contains(player.Username()) || !p.FriendBlocked()))
//		}
	})
}

//RemovePlayer Remove a player from the game world.
func RemovePlayer(p *Player) {
	p.SetRegionRemoved()
	Players.Remove(p)
	Region(p.X(), p.Y()).Players.Remove(p)
	Players.Range(func(player *Player) {
		if player.FriendList.Contains(p.Username()) && (!p.FriendBlocked() || p.FriendList.Contains(player.Username())) {
			player.FriendList.Set(p.Username(), false)
			player.SendPacket(FriendUpdate(p.UsernameHash(), false))
		}
	})
}

//AddNpc Add a NPC to the region.
func AddNpc(n *NPC) {
	Region(n.X(), n.Y()).NPCs.Add(n)
}

//RemoveNpc SetRegionRemoved a NPC from the region.
func RemoveNpc(n *NPC) {
	Region(n.X(), n.Y()).NPCs.Remove(n)
}

//AddItem Add a ground item to the region.
func AddItem(i *GroundItem) {
	Region(i.X(), i.Y()).Items.Add(i)
}

//GetItem Returns the item at x,y with the specified id.  Returns nil if it can not find the item.
func GetItem(x, y, id int) *GroundItem {
	
	region := Region(x, y)
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
	Region(i.X(), i.Y()).Items.Remove(i)
}

//AddObject Add an object to the region.
func AddObject(o *Object) {
	Region(o.X(), o.Y()).Objects.Add(o)
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
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipFullBlock
					continue
				}

				// If it's gone this far, collisionType is 2 (directional blocking, e.g gates etc)
				if o.Direction == byte(North) {
					// Block the tiles east side
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipEast
					// ensure that the neighbors index is valid
					if len(sectorFromCoords(x-1, y).Tiles) > 0 && (areaX > 0 || areaY >= RegionSize) {
						// then block the eastern neighbors west side
						sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY] |= ClipWest
					}
				} else if o.Direction == byte(West) {
					// Block the tiles south side
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipSouth
					// then block the southern neighbors north side
					sectorFromCoords(x, y+1).Tiles[areaX*RegionSize+areaY+1] |= ClipNorth
				} else if o.Direction == byte(South) {
					// Block the tiles west side
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipWest
					// then block the western neighbors east side
					if areaX, areaY := (2304+x+1)%RegionSize, (1776+y-(944*((y+100)/944)))%RegionSize; (areaX+1)*RegionSize+areaY > 2304 {
						sectorFromCoords(x+1, y).Tiles[areaX*RegionSize+areaY] |= ClipEast
					}
				} else if o.Direction == byte(East) {
					// Block the tiles north side
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipNorth
					// ensure that the neighbors index is valid
					if len(sectorFromCoords(x, y-1).Tiles) > 0 && areaX+areaY > 0 {
						// then block the eastern neighbors west side
						sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1] |= ClipSouth
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
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipNorth
			if areaX+areaY > 0 {
				sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1] |= ClipSouth
			}
		} else if o.Direction == 1 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipEast
			if areaX > 0 || areaY >= 48 {
				sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY] |= ClipWest
			}
		} else if o.Direction == 2 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipSwNe
		} else if o.Direction == 3 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipSeNw
		}
	}
}

//RemoveObject SetRegionRemoved an object from the region.
func RemoveObject(o *Object) {
	Region(o.X(), o.Y()).Objects.Remove(o)
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
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] &= ^ClipFullBlock
				} else if o.Direction == 0 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] &= ^ClipEast
					if sectorFromCoords(x-1, y) != nil {
						sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY] &= ^ClipWest
					}
				} else if o.Direction == 2 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] &= ^ClipSouth
					if sectorFromCoords(x, y+1) != nil {
						sectorFromCoords(x, y+1).Tiles[areaX*RegionSize+areaY+1] &= ^ClipNorth
					}
				} else if o.Direction == 4 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] &= ^ClipWest
					if sectorFromCoords(x+1, y) != nil {
						sectorFromCoords(x+1, y).Tiles[(areaX+1)*RegionSize+areaY] &= ^ClipEast
					}
				} else if o.Direction == 6 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] &= ^ClipNorth
					if sectorFromCoords(x, y-1) != nil {
						sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1] &= ^ClipSouth
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
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] &= ^ClipNorth
			if areaX+areaY > 0 {
				sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1] &= ^ClipSouth
			}
		} else if o.Direction == 1 { // Horizontal wall ('__','‾‾') East<->West
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] &= ^ClipEast
			if areaX > 0 || areaY >= 48 {
				sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY] &= ^ClipWest
			}
		} else if o.Direction == 2 { // Diagonal wall ('\','‾|','|_') Southwest<->Northeast
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] &= ^ClipSwNe
		} else if o.Direction == 3 { // Diagonal wall ('/','|‾','_|') Southeast<->Northwest
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] &= ^ClipSeNw
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
	regionLock.RLock()
	defer regionLock.RUnlock()
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
	r := Region(x, y)
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
	regionLock.RLock()
	defer regionLock.RUnlock()
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

var regionLock = sync.RWMutex{}

// internal function to get a region by its row amd column indexes
func get(x, y int) *region {
	if x < 0 {
		x = 0
	}
	if x >= HorizontalPlanes {
		fmt.Println("planeX index out of range", x)
		x = HorizontalPlanes-1
	}
	if y < 0 {
		y = 0
	}
	if y >= VerticalPlanes {
		fmt.Println("planeY out of range", y)
		y = VerticalPlanes-1
	}
	regionLock.Lock()
	defer regionLock.Unlock()
	if regions[x][y] == nil {
		regions[x][y] = &region{x, y, &MobList{}, &MobList{}, &entityList{}, &entityList{}}
	}
	return regions[x][y]
}

//Region Returns the region that corresponds with the given coordinates.  If it does not exist yet, it will allocate a new onr and store it for the lifetime of the application in the regions map.
func Region(x, y int) *region {

	regionLock.Lock()
	defer regionLock.Unlock()
	if regions[x/RegionSize][y/RegionSize] == nil {
		regions[x/RegionSize][y/RegionSize] = &region{x, y, &MobList{}, &MobList{}, &entityList{}, &entityList{}}
	}
	return regions[x/RegionSize][y/RegionSize]}

//surroundingRegions Returns the regions surrounding the given coordinates.  It wil
func (r *region) neighbors() (regions [4]*region) {
	regions[0] = r
	regionX := r.x % RegionSize
	regionY := r.y % RegionSize
	if regionX <= LowerBound {
		regions[1] = get((r.x/RegionSize)-1, r.y/RegionSize)
		if regionY <= LowerBound {
			regions[2] = get((r.x/RegionSize)-1, (r.y/RegionSize)-1)
			regions[3] = get((r.x/RegionSize), (r.y/RegionSize)-1)
		} else {
			regions[2] = get((r.x/RegionSize)-1, (r.y/RegionSize)+1)
			regions[3] = get((r.x/RegionSize), (r.y/RegionSize)+1)
		}
	} else if regionY <= LowerBound {
		regions[1] = get((r.x/RegionSize)+1, (r.y/RegionSize))
		regions[2] = get((r.x/RegionSize)+1, (r.y/RegionSize)-1)
		regions[3] = get((r.x/RegionSize), (r.y/RegionSize)-1)
	} else {
		regions[1] = get((r.x/RegionSize)+1, (r.y/RegionSize))
		regions[2] = get((r.x/RegionSize)+1, (r.y/RegionSize)+1)
		regions[3] = get((r.x/RegionSize), (r.y/RegionSize)+1)
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

//probWeights
type IntProbabilitys = map[int]float64

//Statistical 
func Statistical(rng *rand.Rand, options IntProbabilitys) int {
	if rng == nil {
		rng = rand.New(isaac.New(uint64(time.Now().UnixNano())))
	}

	total := 0.0
	for _, p := range options {
		total += p
	}

	rolled := rng.Float64()*total
	prob := 0.0
	for i, p := range options {
		prob += p
		if rolled <= prob {
			log.Debug("Chose", i, "; hit", rolled,"probability =", prob, "/", total)
			return i
		}
	}
	return -1
}

//WeightedChoice Awesome API call takes map[retVal]probability as input and returns a statistically weighted randomized retVal as output.
//
// The input's mapped value assigned to each key is its return probability, out of the total sum of all return probabilities.
// You can determine the percentage chance of any given input entry being returned by: probability/sumOfAllProbabilities*100
// E.g, if the sum of all probabilities is 100, and you have a total probability of 100, where the first retVal maps to 25.0, the chance it will be returned is 25%
//
// You can make the total anything. Useful for anything that needs to return certain values deterministically more often than others, but randomly.
func WeightedChoice(choices map[int]float64) int {

	return Statistical(rscRand.Rng, choices)
}
/*
func init() {
for i := 0; i < 50; i+=1 {
	WeightedChoice( map[int]float64 {
		1: 25.0,
		3: 25.5,
		6: 66.66,
		21: 25.0,
		23: 25.5,
		26: 66.66,
	})
	}
}
*/