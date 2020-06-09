package db

import (
	"context"

	"github.com/spkaeros/rscgo/pkg/config"
	"github.com/spkaeros/rscgo/pkg/definitions"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

type EntityService interface {
	Objects() []definitions.ScenaryDefinition
	Boundarys() []definitions.BoundaryDefinition
	Tiles() []definitions.TileDefinition
	Items() []definitions.ItemDefinition
	Npcs() []definitions.NpcDefinition
}

var DefaultEntityService *sqlService

func ConnectEntityService() {
	s := newSqlService(config.WorldDriver())
	s.sqlOpen(config.WorldDB())
	DefaultEntityService = s
}

//Objects attempts to load all the scenary object definitions from the SQL service
func (s *sqlService) Objects() (objects []definitions.ScenaryDefinition) {
	s.Lock()
	defer s.Unlock()
	s.context = context.Background()
	db := s.connect(s.context)
	// defer db.Close()
	rows, err := db.QueryContext(s.context, "SELECT id, name, description, LOWER(command_one), LOWER(command_two), type, width, height, modelHeight FROM game_objects ORDER BY id")
	if err != nil {
		log.Warn("Couldn't load entity definitions from sqlService:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		nextDef := definitions.ScenaryDefinition{}
		rows.Scan(&nextDef.ID, &nextDef.Name, &nextDef.Description, &nextDef.Commands[0], &nextDef.Commands[1], &nextDef.CollisionType, &nextDef.Width, &nextDef.Height, &nextDef.ModelHeight)
		objects = append(objects, nextDef)
	}

	return
}

//Boundarys attempts to load all the boundary game object definitions from the SQL service
func (s *sqlService) Boundarys() (boundarys []definitions.BoundaryDefinition) {
	s.Lock()
	defer s.Unlock()
	s.context = context.Background()
	db := s.connect(s.context)
	// defer db.Close()
	rows, err := db.QueryContext(s.context, "SELECT id, name, description, LOWER(command_one), LOWER(command_two), solid, door FROM boundarys ORDER BY id")
	if err != nil {
		log.Warn("Couldn't load entity definitions from sqlService:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		nextDef := definitions.BoundaryDefinition{}
		rows.Scan(&nextDef.ID, &nextDef.Name, &nextDef.Description, &nextDef.Commands[0], &nextDef.Commands[1], &nextDef.Solid, &nextDef.Dynamic)
		boundarys = append(boundarys, nextDef)
	}

	return
}

//Tiles attempts to load all the tile overlay definitions from the SQL service
func (s *sqlService) Tiles() (overlays []definitions.TileDefinition) {
	s.Lock()
	defer s.Unlock()
	s.context = context.Background()
	db := s.connect(s.context)
	// defer db.Close()
	rows, err := db.QueryContext(s.context, "SELECT colour, unknown, objectType FROM tiles")
	if err != nil {
		log.Warn("Couldn't load entity definitions from sqlService:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		nextDef := definitions.TileDefinition{}
		rows.Scan(&nextDef.Color, &nextDef.Visible, &nextDef.Blocked)
		overlays = append(overlays, nextDef)
	}

	return
}

//Items attempts to load all the item definitions from the SQL service
func (s *sqlService) Items() (items []definitions.ItemDefinition) {
	s.Lock()
	defer s.Unlock()
	s.context = context.Background()
	db := s.connect(s.context)
	// defer db.Close()
	rows, err := db.QueryContext(s.context, "SELECT id, name, description, command, base_price, stackable, special, members FROM items ORDER BY id")
	if err != nil {
		log.Warn("Couldn't load entity definitions from sqlService:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		nextDef := definitions.ItemDefinition{}
		rows.Scan(&nextDef.ID, &nextDef.Name, &nextDef.Description, &nextDef.Command, &nextDef.BasePrice, &nextDef.Stackable, &nextDef.Quest, &nextDef.Members)
		items = append(items, nextDef)
	}
	rows.Close()

	rows, err = db.QueryContext(s.context, "SELECT id, skillIndex, level FROM item_wieldable_requirements")
	if err != nil {
		log.Error.Println("Couldn't load entity information from sql database:", err)
		return
	}
	var id, skill, level int
	for rows.Next() {
		rows.Scan(&id, &skill, &level)
		if items[id].Requirements == nil {
			items[id].Requirements = make(map[int]int)
		}
		items[id].Requirements[skill] = level
	}
	rows.Close()

	rows, err = db.QueryContext(s.context, "SELECT id, sprite, type, armour_points, magic_points, prayer_points, range_points, weapon_aim_points, weapon_power_points, pos, femaleOnly FROM item_wieldable")
	if err != nil {
		log.Error.Println("Couldn't load entity information from sql database:", err)
		return
	}
	// TODO: Integrate into ItemDefinition
	for rows.Next() {
		nextDef := definitions.EquipmentDefinition{}
		rows.Scan(&nextDef.ID, &nextDef.Sprite, &nextDef.Type, &nextDef.Armour, &nextDef.Magic, &nextDef.Prayer, &nextDef.Ranged, &nextDef.Aim, &nextDef.Power, &nextDef.Position, &nextDef.Female)
		definitions.Equipment = append(definitions.Equipment, nextDef)
	}

	return
}

//Npcs attempts to load all the npc definitions from the SQL service
func (s *sqlService) Npcs() (npcs []definitions.NpcDefinition) {
	s.Lock()
	defer s.Unlock()
	s.context = context.Background()
	db := s.connect(s.context)
	// defer db.Close()
	rows, err := db.QueryContext(s.context, "SELECT id, name, description, command, hits, attack, strength, defense, hostility FROM npcs ORDER BY id")
	if err != nil {
		log.Warn("Couldn't load entity definitions from sqlService:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		nextDef := definitions.NpcDefinition{}
		rows.Scan(&nextDef.ID, &nextDef.Name, &nextDef.Description, &nextDef.Command, &nextDef.Hits, &nextDef.Attack, &nextDef.Strength, &nextDef.Defense, &nextDef.Hostility)
		npcs = append(npcs, nextDef)
	}

	return
}

//LoadObjectDefinitions Loads game object data into memory for quick access.
func LoadObjectDefinitions() {
	definitions.ScenaryObjects = DefaultEntityService.Objects()
}

//LoadTileDefinitions Loads game tile attribute data into memory for quick access.
func LoadTileDefinitions() {
	definitions.TileOverlays = DefaultEntityService.Tiles()
}

//LoadBoundaryDefinitions Loads game boundary object data into memory for quick access.
func LoadBoundaryDefinitions() {
	definitions.BoundaryObjects = DefaultEntityService.Boundarys()
}

//LoadItemDefinitions Loads game item data into memory for quick access.
func LoadItemDefinitions() {
	definitions.Items = DefaultEntityService.Items()
}

//LoadNpcDefinitions Loads game NPC data into memory for quick access.
func LoadNpcDefinitions() {
	definitions.Npcs = DefaultEntityService.Npcs()
}

//LoadObjectLocations Loads the game objects into memory from the SQLite3 database.
func LoadObjectLocations() {
	database := DefaultEntityService.sqlOpen(config.WorldDB())
	// defer database.Close()
	rows, err := database.Query("SELECT id, direction, boundary, x, y FROM game_object_locations")
	if err != nil {
		log.Warn("Couldn't load SQLite3 database:", err)
		return
	}
	defer rows.Close()
	var id, direction, boundary, x, y int
	for rows.Next() {
		rows.Scan(&id, &direction, &boundary, &x, &y)
		if world.GetObject(x, y) != nil {
			continue
		}
		world.AddObject(world.NewObject(id, direction, x, y, boundary != 0))
	}
}

//LoadNpcLocations Loads the games NPCs into memory from the SQLite3 database.
func LoadNpcLocations() {
	database := DefaultEntityService.sqlOpen(config.WorldDB())
	// defer database.Close()
	rows, err := database.Query("SELECT id, startX, minX, maxX, startY, minY, maxY FROM npc_locations")
	if err != nil {
		log.Warn("Couldn't load SQLite3 database:", err)
		return
	}
	defer rows.Close()
	var id, startX, minX, maxX, startY, minY, maxY int
	for rows.Next() {
		rows.Scan(&id, &startX, &minX, &maxX, &startY, &minY, &maxY)
		world.AddNpc(world.NewNpc(id, startX, startY, minX, maxX, minY, maxY))
	}
}

//LoadItemLocations Loads the games ground items into memory from the SQLite3 database.
func LoadItemLocations() {
	database := DefaultEntityService.sqlOpen(config.WorldDB())
	// defer database.Close()
	rows, err := database.Query("SELECT id, amount, x, y, respawn FROM item_locations")
	if err != nil {
		log.Warn("Couldn't load SQLite3 database:", err)
		return
	}
	defer rows.Close()
	var id, amount, x, y, respawnTime int
	for rows.Next() {
		rows.Scan(&id, &amount, &x, &y, &respawnTime)
		world.AddItem(world.NewPersistentGroundItem(id, amount, x, y, respawnTime))
	}
}

//SaveObjectLocations Clears definitions.db game object locations and repopulates it with the current game locations.
func SaveObjectLocations() int {
	database := DefaultEntityService.sqlOpen(config.WorldDB())
	// defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Info.Println("Error starting transaction for saving object locations:", err)
		return -1
	}

	stmt, err := tx.Exec("DELETE FROM game_object_locations")
	if err != nil {
		tx.Rollback()
		log.Warn("Error clearing object locations to save new ones:", err)
		return -1
	}
	if count, err := stmt.RowsAffected(); count < 1 || err != nil {
		if err != nil {
			log.Warn("Error inserting new game object location to definitions.db:", err)
			return -1
		}
		log.Warn("Rows affected < 1 in game object location insert:%d\n", count)
		return -1
	}

	totalInserts := 0
	for _, v := range world.GetAllObjects() {
		stmt, err := tx.Exec("INSERT INTO game_object_locations(id, direction, x, y, boundary) VALUES(?, ?, ?, ?, ?)", v.ID, v.Direction, v.X(), v.Y(), v.Boundary)
		if err != nil {
			log.Warn("Error inserting game object location to database:", err)
			continue
		}
		if count, err := stmt.RowsAffected(); count < 1 || err != nil {
			if err != nil {
				log.Warn("Error inserting new game object location to definitions.db:", err)
				continue
			}
			log.Warn("Rows affected < 1 in game object location insert:%d\n", count)
			continue
		}
		totalInserts++
	}

	if err := tx.Commit(); err != nil {
		log.Warn("Couldn't commit game object locations:", err)
		return -1
	}

	return totalInserts
}
