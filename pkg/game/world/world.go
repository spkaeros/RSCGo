package world

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/spkaeros/rscgo/pkg/config"
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

//BroadcastLogin Broadcasts the login status of player to the whole game.
func (m *PlayerMap) BroadcastLogin(player *Player, online bool) {
	m.Range(func(rangedPlayer *Player) {
		if player.Friends(rangedPlayer.UsernameHash()) {
			if !rangedPlayer.FriendBlocked() || rangedPlayer.Friends(rangedPlayer.UsernameHash()) {
				player.FriendList[rangedPlayer.UsernameHash()] = online
			}
		}
		if rangedPlayer.Friends(player.UsernameHash()) {
			if !player.FriendBlocked() || player.Friends(rangedPlayer.UsernameHash()) {
				rangedPlayer.FriendList[player.UsernameHash()] = online
				rangedPlayer.SendPacket(FriendUpdate(player.UsernameHash(), online))
			}
		}
	})
}

//region Represents a 48x48 section of map.  The purpose of this is to keep track of entities in the entire world without having to allocate tiles individually, which would make search algorithms slower and utilizes a great deal of memory.
type region struct {
	Players *entityList
	NPCs    *entityList
	Objects *entityList
	Items   *entityList
}

var regions [HorizontalPlanes][VerticalPlanes]*region

// A convencience type for tickable task closures, as typing func signatures out as return signatures gets tiresome.
type Task func() bool
type taskSet map[string]Task

type TaskCollection struct {
	taskSet
	sync.RWMutex
}

func (t *TaskCollection) Range(fn func(string, Task)) {
	t.RLock()
	for name, task := range t.taskSet {
		fn(name, task)
	}
	t.RUnlock()
}

func (t *TaskCollection) ExecuteSequentially() {
	var removed []string
	t.Lock()
	for name, task := range t.taskSet {
		//start := time.Now()
		if task() {
			removed = append(removed, name)
		}
		//log.Info.Printf("tickTask--%s; finished executing in %v", name, time.Since(start))
	}
	for _, taskName := range removed {
		delete(Tickables.taskSet, taskName)
	}
	t.Unlock()
}

func (t *TaskCollection) Add(name string, fn Task) {
	t.Lock()
	t.taskSet[name] = fn
	t.Unlock()
}

func (t *TaskCollection) Get(name string) Task {
	t.RLock()
	defer t.RUnlock()
	return t.taskSet[name]
}

func (t *TaskCollection) Remove(name string) {
	t.Lock()
	delete(t.taskSet, name)
	t.Unlock()
}

var Tickables = &TaskCollection{
	taskSet: make(taskSet),
}

//IsValid Returns true if the tile at x,y is within world boundaries, false otherwise.
func WithinWorld(x, y int) bool {
	return x <= MaxX && x >= 0 && y >= 0 && y <= MaxY
}

//AddPlayer Add a player to the region.
func AddPlayer(p *Player) {
	getRegion(p.X(), p.Y()).Players.Add(p)
}

//RemovePlayer SetRegionRemoved a player from the region.
func RemovePlayer(p *Player) {
	getRegion(p.X(), p.Y()).Players.Remove(p)
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
		def := ObjectDefs[o.ID]
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
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipFullBlock
				} else if o.Direction == 0 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipEast
					if sectorFromCoords(x-1, y) != nil && (areaX > 0 || areaY >= 48) {
						sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY].CollisionMask |= ClipWest
					}
				} else if o.Direction == 2 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipSouth
					if sectorFromCoords(x, y+1) != nil {
						sectorFromCoords(x, y+1).Tiles[areaX*RegionSize+areaY+1].CollisionMask |= ClipNorth
					}
				} else if o.Direction == 4 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipWest
					if sectorFromCoords(x+1, y) != nil {
						sectorFromCoords(x+1, y).Tiles[(areaX+1)*RegionSize+areaY].CollisionMask |= ClipEast
					}
				} else if o.Direction == 6 {
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipNorth
					if sectorFromCoords(x, y-1) != nil {
						sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1].CollisionMask |= ClipSouth
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
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipNorth
			if sectorFromCoords(x, y-1) != nil {
				sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1].CollisionMask |= ClipSouth
			}
		} else if o.Direction == 1 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipEast
			if sectorFromCoords(x-1, y) != nil && (areaX > 0 || areaY >= 48) {
				sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY].CollisionMask |= ClipWest
			}
		} else if o.Direction == 2 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipDiag1
		} else if o.Direction == 3 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask |= ClipDiag2
		}
	}
}

//RemoveObject SetRegionRemoved an object from the region.
func RemoveObject(o *Object) {
	getRegion(o.X(), o.Y()).Objects.Remove(o)
	if !o.Boundary {
		def := ObjectDefs[o.ID]
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
		def := BoundaryDefs[o.ID]
		if def.Traversable != 1 {
			return
		}
		x, y := o.X(), o.Y()
		areaX := (2304 + x) % RegionSize
		areaY := (1776 + y - (944 * ((y + 100) / 944))) % RegionSize
		if o.Direction == 0 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= ^ClipNorth
			if sectorFromCoords(x, y-1) != nil {
				sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1].CollisionMask &= ^ClipSouth
			}
		} else if o.Direction == 1 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= ^ClipEast
			if sectorFromCoords(x-1, y) != nil {
				sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY].CollisionMask &= ^ClipWest
			}
		} else if o.Direction == 2 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= ^ClipDiag1
		} else if o.Direction == 3 {
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY].CollisionMask &= ^ClipDiag2
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

//BoundedChance should return true (percent)% of the time, and false (100-percent)% of the time.
// It uses ISAAC64+ to provide randomness. If percent is not inside the specified boundaries, it will be
// ignored and the appropriate boundary will replace it.
//
// percent defines the percentage of chance for this check to pass, as long as it's between minPercent and maxPercent.
//
// minPercent defines the minimum percentage of chance for this check to pass.  if percent is lower than this, it will
// be ignored and this will be used.
//
// maxPercent defines the maximum percentage of chance for this check to pass.  if percent is higher than this, it will
// be ignored and this will be used.
func BoundedChance(percent float64, minPercent, maxPercent float64) bool {
	if minPercent < 0.0 {
		minPercent = 0.0
	}
	if maxPercent > 100.0 {
		maxPercent = 100.0
	}
	if minPercent > maxPercent {
		maxPercent, minPercent = minPercent, maxPercent
	}
	return float64(rscRand.Uint8()) <= math.Max(minPercent, math.Min(maxPercent, percent))/100.0*256.0
}

//Chance should return true (percent)% of the time, and false (100-percent)% of the time.
// It uses ISAAC64+ to provide randomness.
//
// percent defines the percentage of chance for this check to pass.
func Chance(percent float64) bool {
	return BoundedChance(percent, 0.0, 100.0)
}

//WeightedChance Awesome API call takes map[retVal]probability as input and returns a statistically weighted randomized retVal as output.
//
// The input's mapped value assigned to each key is its return probability, out of the total sum of all return probabilities.
// You can determine the percentage chance of any given input entry being returned by: probability/sumOfAllProbabilities*100
// E.g, if the sum of all probabilities is 100, and you have a total probability of 100, where the first retVal maps to 25.0, the chance it will be returned is 25%
//
// You can make the total anything. Useful for anything that needs to return certain values deterministically more often than others, but randomly.
func WeightedChoice(choices map[int]float64) int {
	total := 0.0
	totalProb := 0.0
	for ret, probability := range choices {
		if probability > 100 && config.Verbose() {
			log.Warning.Println("Probability of a single WeightedChance entry is over 100%:{entryVal:", ret, "; probability:", probability, "}; it is likely you are not using this function properly.")
		}
		totalProb += probability
	}

	// We determine the real upper limit of the cumulative probability.
	// if the sum of all probabilities adds up to 100.0% or less, then the upper limit will be 65535
	// if it's less, any unaccounted for percentage will return -1.
	// if it's more, it'll scale 65535 up by however much percentage it needs to handle all inputs provided, and subsequently
	// all individual probabilities must be scaled down by the same figure to account for the difference.
	// e.g passing 3 values in with 50.0 probability on each of them will result in each entry returning at a
	// rate of 33.333~% instead of 50% each, because 50/150=33.333~
	cumulativeProbability := math.Max(1.0, totalProb/100)
	upperBound := cumulativeProbability * math.MaxUint16
	hit := float64(rscRand.Int31N(1, int(upperBound)))
	if config.Verbosity >= 3 {
		log.Info.Printf("WeightedChoice: Upper bound for total probability:%d (total was: %.2f%% of 65535; the RNG lower bound) {\n", int(upperBound), cumulativeProbability*100)
		log.Info.Printf("\tRolled: %d/%d;\n", int(hit), int(upperBound))
		defer log.Info.Println("};")
	}
	for choice, prob := range choices {
		newProb := prob / 100 * math.MaxUint16 / upperBound
		if config.Verbosity >= 3 {
			log.Info.Printf("\tentry{val:%d; hit range between %d - %d (%.2f%% chance)}", choice, int(total), int(total+(newProb*upperBound)), newProb*100)
		}
		total += newProb * upperBound
		if hit < total {
			return choice
		}
	}
	if config.Verbosity >= 3 {
		log.Info.Println("Rolled value did not return anything!")
	}
	return -1
}

//MeleeExperience returns how much combat experience to award for killing an opponent with melee.
func MeleeExperience(victim MobileEntity) float64 {
	return float64((victim.Skills().CombatLevel()*2.0)+10.0) * 1.5
}

//CombatPrefix Returns the chat prefix to colorize combat levels in right click menus and such.
// The color fades red as the target compares better than you, or fades green as the target compares worse than you.
// White indicates an equal target.
func CombatPrefix(delta int) string {
	// They're stronger
	if delta < -9 {
		return "@red@"
	}
	if delta < -6 {
		return "@or3@"
	}
	if delta < -3 {
		return "@or2@"
	}
	if delta < 0 {
		return "@or1@"
	}

	// They're weaker
	if delta > 9 {
		return "@gre@"
	}
	if delta > 6 {
		return "@gr3@"
	}
	if delta > 3 {
		return "@gr2@"
	}
	if delta > 0 {
		return "@gr1@"
	}

	// They're the same
	return "@whi@"
}
