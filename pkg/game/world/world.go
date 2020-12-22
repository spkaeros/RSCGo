package world

import (
	"math"
	// "strconv"
	"sync"
	"time"

	// "github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/isaac"
	"github.com/spkaeros/rscgo/pkg/tasks"
	"github.com/spkaeros/rscgo/pkg/game/entity"
	"github.com/spkaeros/rscgo/pkg/log"
	rscRand "github.com/spkaeros/rscgo/pkg/rand"
)

const (
	TicksDay         = 135000
	TicksHour        = 5625
	TicksTwentyMin   = 1875
	TicksMinute      = 100
	TickMinute       = TicksMinute * TickMillis
	TickHour         = TicksHour * TickMillis
	TickDay          = TicksDay * TickMillis
	TickMillis       = 640 * time.Millisecond
	ClientTickMillis = TickMillis >> 5

	//MaxX Width of the game
	MaxX = 944
	// MaxX = 960
	//MaxY Height of the game
	// MaxY = 944 * 4 - 16 // View area for player is 16 tiles wide
	MaxY = 944 * 4 // 4 planes, 944 tiles per plane
	//TODO: I need to investigate other measurement windows to use for the map
	// for instance, wilderness levels increase every 6 tiles
	// The client scene mesh uses 128x128 unit tiles to display things
)

func CurrentTick() int {
	return int(tasks.Ticks.Load())
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

type indexQueue []int

func (q *indexQueue) Push(i int) {
	*q = append([]int{i}, *q...)
}

func (q *indexQueue) Pop() int {
	if len(*q) <= 0 {
		return -1
	}
	idx := (*q)[0]
	if len(*q) == 1 {
		*q = (*q)[:0]
	} else {
		*q = (*q)[1:]
	}
	return idx
}

type PlayersList map[*Player]struct{}

type PlayerList struct {
	sync.RWMutex
	PlayersList
	free   indexQueue
	curIdx int
}

//Players Collection containing all of the active client, by index and username hash, guarded by a mutex
var Players = &PlayerList{free: make(indexQueue, 0, 1250), PlayersList: make(PlayersList)}

//FindHash Returns the client with the base37 username `hash` if it exists and true, otherwise returns nil and false.
func (m *PlayerList) FindHash(hash uint64) (*Player, bool) {
	m.RLock()
	defer m.RUnlock()
	for player := range m.PlayersList {
		if player.UsernameHash() == hash {
			return player, true
		}
	}
	return nil, false
}

//FromIndex Returns the client with the index `index` if it exists and true, otherwise returns nil and false.
func (m *PlayerList) FindIndex(index int) (*Player, bool) {
	if index < 0 {
		return nil, false
	}
	m.RLock()
	defer m.RUnlock()
	for player := range m.PlayersList {
		if player.ServerIndex() == index {
			return player, true
		}
	}
	return nil, false
}

//Find Returns the slot that this player occupies in the set.
func (m *PlayerList) Find(player *Player) int {
	return m.ForEach(func(p *Player) bool {
		return p == player
	})
}

//Contains Returns true if this player is assigned to a slot in the set, otherwise returns false.
func (m *PlayerList) Contains(player *Player) bool {
	return m.Find(player) > -1
}

//ContainsHash Returns true if there is a client mapped to this username hash is in this collection, otherwise returns false.
func (m *PlayerList) ContainsHash(hash uint64) bool {
	_, ret := m.FindHash(hash)
	return ret
}

//Put Finds the lowest available empty slot in the list, and puts the player there.
// This will also assign the players server index variable (*Player.Index) to the assigned slot.
func (m *PlayerList) Put(player *Player) {
	if m.Contains(player) {
		log.Warn("Player list double-entry attempted!")
		return
	}
	player.SetServerIndex(m.nextSlot())
	// for player.ServerIndex() >= len(m.PlayersList) {
		// m.PlayersList = append(m.PlayersList, nil)
	// }
	m.Lock()
	defer m.Unlock()
	m.PlayersList[player] = struct{}{}// = append(m.PlayersList, player)
	// m.PlayersList[player.ServerIndex()] = player
}

//Remove Removes a client from the set.
func (m *PlayerList) Remove(player *Player) {
	slot := player.ServerIndex()
	if slot == -1 {
		log.Warn("Error: Player with non-valid index being removed!")
		return
	}
	player.SetServerIndex(-1)
	m.Lock()
	defer m.Unlock()
	delete(m.PlayersList, player)
	// m.PlayersList[slot] = nil
	// if len(m.PlayersList) == 1 {
		// m.PlayersList = m.PlayersList[:0]
	// } else {
		// if slot < len(m.PlayersList)-1 {
			// m.PlayersList = append(m.PlayersList[:slot], m.PlayersList[slot+1:]...)
		// } else {
			// m.PlayersList = m.PlayersList[:slot]
		// }
	// }
	m.free.Push(slot)
}

//Range Calls action for every active client in the collection.
func (m *PlayerList) Range(action func(*Player)) {
	m.RLock()
	defer m.RUnlock()
	for p := range m.PlayersList {
		if p != nil {
			action(p)
		}
	}
}

//Range Calls action for every active client in the collection.
func (m *PlayerList) ForEach(action func(*Player) bool) int {
	m.RLock()
	defer m.RUnlock()
	for p := range m.PlayersList {
		if p != nil {
			if action(p) {
				return p.ServerIndex()
			}
		}
	}

	return -1
}

//Size Returns the size of the active client collection.
func (m *PlayerList) Size() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.PlayersList)
}

//NextIndex Returns the lowest available index for the client to be mapped to.
func (m *PlayerList) nextSlot() int {
	m.Lock()
	defer m.Unlock()
	idx := m.free.Pop()
	if idx == -1 {
		defer func() { m.curIdx += 1 }()
		return m.curIdx
	}
	return idx
}

func (m *PlayerList) Set() []*Player {
	keys := make([]*Player, m.Size())

	i := 0
	for k := range m.PlayersList {
	    keys[i] = k
	    i++
	}
	return keys
}

func (m *PlayerList) AsyncRange(fn func(*Player)) {
	m.RLock()
	sz := len(m.PlayersList)
	done := make(chan struct{}, sz)
	for player := range m.PlayersList {
		go func(player *Player) {
			fn(player)
			// if i == sz {
				// done <- struct{}{}
			// }
			done <- struct{}{}
		}(player)
	}
	m.RUnlock()
	defer close(done)
	// <-done
	for i := 0; i < sz; i += 1 {
		select {
		case _, ok := <-done:
			if !ok {
				log.Debug("premature AsyncRange done-signal")
				break
			}
			continue
		}
	}
}

//region Represents a 48x48 section of map.  The purpose of this is to keep track of entities in the entire world without having to allocate tiles individually, which would make search algorithms slower and utilizes a great deal of memory.
type region struct {
	x       int
	y       int
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
	Players.Put(p)
	Region(p.X(), p.Y()).Players.Add(p)
	// Players.Range(func(player *Player) {
		// if player.FriendList.Contains(p.Username()) && (!p.FriendBlocked() || p.FriendList.Contains(player.Username())) {
			// player.FriendList.Set(p.Username(), true)
			// player.WritePacket(FriendUpdate(p.UsernameHash(), true))
		// }
// 
		// //		if player.FriendList.Contains(p.Username()) {
		// //			player.WritePacket(FriendUpdate(p.UsernameHash(), p.FriendList.Contains(player.Username()) || !p.FriendBlocked()))
		// //		}
	// })
}

//RemovePlayer Remove a player from the game world.
func RemovePlayer(p *Player) {
	p.SetRegionRemoved()
	Region(p.X(), p.Y()).Players.Remove(p)
	Players.Remove(p)
	// Players.Range(func(player *Player) {
		// if player.FriendList.Contains(p.Username()) && (!p.FriendBlocked() || p.FriendList.Contains(player.Username())) {
			// player.FriendList.Set(p.Username(), false)
			// player.WritePacket(FriendUpdate(p.UsernameHash(), false))
		// }
	// })
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
	region.Items.RLock()
	defer region.Items.RUnlock()
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

type mapBarrier interface {
	// should return true if this barrier is totally impassable, blocking all mobs
	// from walking on it so long as it stands in the way
	// example being things like counters or altars...
	Solid() bool

	// should return true if this barrier only blocks mobs coming at them from specific
	// directions.  Typically results in one to three tiles being blocked as source locations
	// Example being bank doors or member gates, so on..
	Door() bool

	// returns true if both Door and Solid return false
	// just a shortcut for doing this check manually
	Passable() bool

	// How many rows this barrier occupies
	Width() int
	// How many columns this barrier occupies
	Height() int

	// returns true if we have a stored entry for this type of barrier
	// otherwise, returns false
	Defined() bool
}

//AddObject Add an object to the region.
func AddObject(o *Object) {
	Region(o.X(), o.Y()).Objects.Add(o)
	data := o.TypeData()
	if !data.Defined() {
		return
	}
	if data.Passable() {
		return
	}
	if o.Boundary {
		x,y := o.X(),o.Y()
		areaX := (2304+x) % RegionSize
		areaY := (1776+y - (944*((y+100)/944))) % RegionSize
		switch o.Direction {
		case 0:
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipNorth
			if areaX+areaY > 0 {
				sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1] |= ClipSouth
			}
		case 1:
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipEast
			if areaX > 0 || areaY >= 48 {
				sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY] |= ClipWest
			}
		case 2:
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipSwNe
		case 3:
			sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipSeNw
		}
		return
	}
	width, height := data.Width(), data.Height()
	if o.Direction&3 != 0 {
		width, height = height, width
	}
	for dx := 0; dx < width; dx++ {
		for dy := 0; dy < height; dy++ {
			x,y := o.X()+dx,o.Y()+dy
			areaX := (2304+x) % RegionSize
			areaY := (1776+y - (944*((y+100)/944))) % RegionSize
			if len(sectorFromCoords(x, y).Tiles) <= 0 {
				return
			}
			if data.Solid() {
				sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipFullBlock
				continue
			}
			switch int(o.Direction) {
			case 0: // block movement from the west
				// Block the tiles east side
				sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipEast
				// ensure that the neighbors index is valid
				if len(sectorFromCoords(x-1, y).Tiles) > 0 && (areaX > 0 || areaY >= RegionSize) {
					// then block the eastern neighbors west side
					sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY] |= ClipWest
				}
			case 2: // block movement from the north
				// Block the tiles south side
				sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipSouth
				// then block the southern neighbors north side
				sectorFromCoords(x, y+1).Tiles[areaX*RegionSize+areaY+1] |= ClipNorth
			case 4: // block movement from the east
				// Block the tiles west side
				sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipWest
				// then block the western neighbors east side
				if areaX, areaY := (2304+x+1)%RegionSize, (1776+y-(944*((y+100)/944)))%RegionSize; (areaX+1)*RegionSize+areaY > 2304 {
					sectorFromCoords(x+1, y).Tiles[areaX*RegionSize+areaY] |= ClipEast
				}
			case 6: // block movement from the south
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
	// if !o.Boundary {
		// scenary := definitions.ScenaryObjects[o.ID]
		// // type 0 is used when the object causes no collisions of any sort.
		// // type 1 is used when the object fully blocks the tile(s) that it sits on.  Marks tile as fully blocked.
		// // type 2 is used when the object mimics a boundary, e.g for gates and the like.
		// // type 3 is used when the object mimics an opened door-type boundary, e.g opened gates and the like.
		// if scenary.SolidityType == 0 || scenary.SolidityType == 3 {
			// return
		// }
		// width, height := scenary.Height(), scenary.Width()
		// //if o.Direction == 0 || o.Direction == 4 {
		// if o.Direction%4 == 0 {
			// width, height = height, width
		// }
		// for x := o.X(); x < o.X()+width; x++ {
			// for y := o.Y(); y < o.Y()+height; y++ {
				// areaX := (2304 + x) % RegionSize
				// areaY := (1776 + y - (944 * ((y + 100) / 944))) % RegionSize
				// if len(sectorFromCoords(x, y).Tiles) <= 0 {
					// log.Warning.Println("ERROR: Sector with no tiles at:" + strconv.Itoa(x) + "," + strconv.Itoa(y) + " (" + strconv.Itoa(areaX) + "," + strconv.Itoa(areaY) + "\n")
					// return
				// }
				// if scenary.SolidityType == 1 {
					// // Blocks the whole tile.  Can not walk on it from any direction
					// sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipFullBlock
					// continue
				// }
// 
				// // If it's gone this far, collisionType is 2 (directional blocking, e.g gates etc)
				// if o.Direction == byte(North) {
					// // Block the tiles east side
					// sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipEast
					// // ensure that the neighbors index is valid
					// if len(sectorFromCoords(x-1, y).Tiles) > 0 && (areaX > 0 || areaY >= RegionSize) {
						// // then block the eastern neighbors west side
						// sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY] |= ClipWest
					// }
				// } else if o.Direction == byte(West) {
					// // Block the tiles south side
					// sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipSouth
					// // then block the southern neighbors north side
					// sectorFromCoords(x, y+1).Tiles[areaX*RegionSize+areaY+1] |= ClipNorth
				// } else if o.Direction == byte(South) {
					// // Block the tiles west side
					// sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipWest
					// // then block the western neighbors east side
					// if areaX, areaY := (2304+x+1)%RegionSize, (1776+y-(944*((y+100)/944)))%RegionSize; (areaX+1)*RegionSize+areaY > 2304 {
						// sectorFromCoords(x+1, y).Tiles[areaX*RegionSize+areaY] |= ClipEast
					// }
				// } else if o.Direction == byte(East) {
					// // Block the tiles north side
					// sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipNorth
					// // ensure that the neighbors index is valid
					// if len(sectorFromCoords(x, y-1).Tiles) > 0 && areaX+areaY > 0 {
						// // then block the eastern neighbors west side
						// sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1] |= ClipSouth
					// }
				// }
// 
			// }
		// }
	// } else {
		// boundary := definitions.BoundaryObjects[o.ID]
		// if !boundary.Solid {
			// // Doorframes and some other stuff collide with nothing.
			// return
		// }
		// x, y := o.X(), o.Y()
		// areaX := (2304 + x) % RegionSize
		// areaY := (1776 + y - (944 * ((y + 100) / 944))) % RegionSize
		// if len(sectorFromCoords(x, y).Tiles) <= 0 {
			// log.Warn("ERROR: Sector with no tiles at:" + strconv.Itoa(x) + "," + strconv.Itoa(y) + " (" + strconv.Itoa(areaX) + "," + strconv.Itoa(areaY) + "\n")
			// return
		// }
		// if o.Direction == 0 {
			// sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipNorth
			// if areaX+areaY > 0 {
				// sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1] |= ClipSouth
			// }
		// } else if o.Direction == 1 {
			// sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipEast
			// if areaX > 0 || areaY >= 48 {
				// sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY] |= ClipWest
			// }
		// } else if o.Direction == 2 {
			// sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipSwNe
		// } else if o.Direction == 3 {
			// sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] |= ClipSeNw
		// }
	// }
}

//RemoveObject SetRegionRemoved an object from the region.
func RemoveObject(o *Object) {
	Region(o.X(), o.Y()).Objects.Remove(o)
	data := o.TypeData()
	if data.Passable() {
		return
	}
	if !o.Boundary {
		// scenary := definitions.ScenaryObjects[o.ID]
		// type 0 is used when the object causes no collisions of any sort.
		// type 1 is used when the object fully blocks the tile(s) that it sits on.  Marks tile as fully blocked.
		// type 2 is used when the object mimics a boundary, e.g for gates and the like.
		// type 3 is used when the object mimics an opened door-type boundary, e.g opened gates and the like.
		width, height := data.Width(), data.Height()
		//if o.Direction == byte(North) || o.Direction == byte(South) {
		if o.Direction&3 != 0 {
			// reverse measurements for directions 0(North) and 4(South), as scenary measurements
			// are oriented vertically by default
			width,height = height,width
		}
		for dx := 0; dx < width; dx++ {
			for dy := 0; dy < height; dy++ {
				x,y := o.X()+dx,o.Y()+dy
				areaX := (2304 + x) % RegionSize
				areaY := (1776 + y - (944 * ((y + 100) / 944))) % RegionSize
				if data.Solid() {
					// This indicates a solid object.  Impassable and blocks the whole tile.
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] &= ^ClipFullBlock
					continue
				}
				switch o.Direction {
				case 0:
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] &= ^ClipEast
					if sectorFromCoords(x-1, y) != nil {
						sectorFromCoords(x-1, y).Tiles[(areaX-1)*RegionSize+areaY] &= ^ClipWest
					}
				case 2:
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] &= ^ClipSouth
					if sectorFromCoords(x, y+1) != nil {
						sectorFromCoords(x, y+1).Tiles[areaX*RegionSize+areaY+1] &= ^ClipNorth
					}
				case 4:
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] &= ^ClipWest
					if sectorFromCoords(x+1, y) != nil {
						sectorFromCoords(x+1, y).Tiles[(areaX+1)*RegionSize+areaY] &= ^ClipEast
					}
				case 6:
					sectorFromCoords(x, y).Tiles[areaX*RegionSize+areaY] &= ^ClipNorth
					if sectorFromCoords(x, y-1) != nil {
						sectorFromCoords(x, y-1).Tiles[areaX*RegionSize+areaY-1] &= ^ClipSouth
					}
				}
			}
		}
	} else {
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
func GetAllObjects() (list []entity.Entity) {
	regionLock.RLock()
	defer regionLock.RUnlock()
	for _, xR := range regions {
		for _, yR := range xR {
			if yR != nil {
				yR.Objects.Range(func(e entity.Entity) {
					list = append(list, e)
				})
			}
		}
			
	}
	return list
}

//GetObject If there is an object at these coordinates, returns it.  Otherwise, returns nil.
func GetObject(x, y int) *Object {
	r := Region(x, y)
	r.Objects.RLock()
	defer r.Objects.RUnlock()
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
			Region(x, y).NPCs.RangeNpcs(func(n *NPC) bool {
				if n.ID == id && n.LongestDelta(point) < minDelta {
					minDelta = n.LongestDelta(point)
					npc = n
				}
				return false
			})
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
			Region(x, y).NPCs.RangeNpcs(func(n *NPC) bool {
				if n.ID == id && n.LongestDelta(point) < minDelta {
					minDelta = n.LongestDelta(point)
					npc = n
				}
				return false
			})
		}
	}
	return npc
}

var regionLock = sync.RWMutex{}
// 
// func init() {
	// regionLock.Lock()
	// defer regionLock.Unlock()
// 
	// for x := 0; x < MaxX/RegionSize+1; x += 1 {
		// for y := 0; y < MaxY/RegionSize+1; y += 1 {
			// regions[x][y] = &region{x*RegionSize, y*RegionSize, NewMobList(), NewMobList(), &entityList{}, &entityList{}}
		// }
	// }
// }

// internal function to get a region by its row amd column indexes
func get(x, y int) *region {
	x = int(math.Min(math.Max(0, float64(x/RegionSize)), HorizontalPlanes - 1))
	y = int(math.Min(math.Max(0, float64(y/RegionSize)), VerticalPlanes - 1))
	regionLock.Lock()
	defer regionLock.Unlock()
	if regions[x][y] == nil {
		regions[x][y] = &region{x, y, NewMobList(), NewMobList(), &entityList{}, &entityList{}}
	}
	return regions[x][y]
}

//Region Returns the region that corresponds with the given coordinates.  If it does not exist yet, it will allocate a new onr and store it for the lifetime of the application in the regions map.
func Region(x, y int) *region {
	return get(x, y)
}

//VisibleRegions Returns the regions surrounding the given coordinates.
func VisibleRegions(x, y int) (regions [4]*region) {
	regionX, regionY := x, y
	regions[0] = get(regionX, regionY)
	regions[1] = get(regionX+RegionSize, regionY)
	regions[2] = get(regionX, regionY+RegionSize)
	regions[3] = get(regionX+RegionSize, regionY+RegionSize)
	// if x % RegionSize <= LowerBound {
		// regionX -= RegionSize
	// } else {
		// regionX += RegionSize
	// }
	// regions[1] = get(regionX, regionY)
	// regions[2] = get(x, regionY)
	// if y % RegionSize <= LowerBound {
		// regionY -= RegionSize+(RegionSize/2)
	// } else {
		// regionY += RegionSize+(RegionSize/2)
	// }
	// regions[3] = get(regionX, regionY)
	
	return
}
func VisibleRegionsFrom(l entity.Location) (regions [4]*region) {
	return VisibleRegions(l.X(), l.Y())
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
	if minPercent > maxPercent {
		maxPercent, minPercent = minPercent, maxPercent
	}
	threshold := uint8(math.Max(minPercent, math.Min(maxPercent, percent)) / maxPercent * 256.0)
	return ChanceByte(int(threshold))
}

//ChanceByte Grabs a single 8-bit unsigned byte out of the rscgo/rand pkg, and returns true if it's less than or equals the provided threshold.
func ChanceByte(threshold int) bool {
	return rscRand.Byte() <= uint8(threshold)
}

//Chance should return true (percent)% of the time, and false (100-percent)% of the time.
// It uses ISAAC64+ to provide randomness.
//
// percent defines the percentage of chance for this check to pass.
func Chance(percent float64) bool {
	return ChanceByte(int(percent/100.0*256.0))
	// return BoundedChance(percent, 0.0, 100.0)
}

//probWeights
type IntProbabilitys = map[int]float64

//Statistical
func Statistical(rng *isaac.ISAAC, options IntProbabilitys) int {
	if rng == nil {
		rng = rscRand.Rng
	}

	total := 0.0
	for _, p := range options {
		total += p
	}

	rolled := rng.Float64() * total
	prob := 0.0
	for i, p := range options {
		prob += p
		if rolled <= prob {
			log.Debug("Chose", i, "; hit", rolled, "probability =", prob, "/", total)
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
func WeightedChoice(choices IntProbabilitys) int {
	return Statistical(rscRand.Rng, choices)
}
