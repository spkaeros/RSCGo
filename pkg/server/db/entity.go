package db

import (
	"github.com/spkaeros/rscgo/pkg/server/config"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"strings"
)

func GetEquipmentDefinition(id int) *world.EquipmentDefinition {
	for _, e := range world.EquipmentDefs {
		if e.ID == id {
			return &e
		}
	}

	return nil
}

//LoadObjectDefinitions Loads game object data into memory for quick access.
func LoadObjectDefinitions() {
	database := Open(config.WorldDB())
	defer database.Close()
	rows, err := database.Query("SELECT id, name, description, command_one, command_two, type, width, height, ground_item_var FROM `game_objects`")
	if err != nil {
		log.Error.Println("Couldn't load SQLite3 database:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		nextDef := world.ObjectDefinition{Commands: make([]string, 2)}
		rows.Scan(&nextDef.ID, &nextDef.Name, &nextDef.Description, &nextDef.Commands[0], &nextDef.Commands[1], &nextDef.Type, &nextDef.Width, &nextDef.Height, &nextDef.Length)
		for i, c := range nextDef.Commands {
			nextDef.Commands[i] = strings.ToLower(c)
		}
		world.Objects = append(world.Objects, nextDef)
	}
}

func LoadEquipmentDefinitions() {
	database := Open(config.WorldDB())
	defer database.Close()
	rows, err := database.Query("SELECT id, sprite, type, armour_points, magic_points, prayer_points, range_points, weapon_aim_points, weapon_power_points, pos, femaleOnly FROM `item_wieldable`")
	if err != nil {
		log.Error.Println("Couldn't load SQLite3 database:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		nextDef := world.EquipmentDefinition{}
		rows.Scan(&nextDef.ID, &nextDef.Sprite, &nextDef.Type, &nextDef.Armour, &nextDef.Magic, &nextDef.Prayer, &nextDef.Ranged, &nextDef.Aim, &nextDef.Power, &nextDef.Position, &nextDef.Female)
		world.EquipmentDefs = append(world.EquipmentDefs, nextDef)
	}
}

//LoadTileDefinitions Loads game tile attribute data into memory for quick access.
func LoadTileDefinitions() {
	database := Open(config.WorldDB())
	defer database.Close()
	rows, err := database.Query("SELECT colour, unknown, objectType unknown FROM `tiles`")
	defer rows.Close()
	if err != nil {
		log.Error.Println("Couldn't load SQLite3 database:", err)
		return
	}
	for rows.Next() {
		nextDef := world.TileDefinition{}
		rows.Scan(&nextDef.Color, &nextDef.Visible, &nextDef.ObjectType)
		world.Tiles = append(world.Tiles, nextDef)
	}
}

//LoadBoundaryDefinitions Loads game boundary object data into memory for quick access.
func LoadBoundaryDefinitions() {
	database := Open(config.WorldDB())
	defer database.Close()
	rows, err := database.Query("SELECT id, name, description, command_one, command_two, door_type, unknown FROM `doors` ORDER BY id")
	defer rows.Close()
	if err != nil {
		log.Error.Println("Couldn't load SQLite3 database:", err)
		return
	}
	for rows.Next() {
		nextDef := world.BoundaryDefinition{Commands: make([]string, 2)}
		for i, c := range nextDef.Commands {
			nextDef.Commands[i] = strings.ToLower(c)
		}
		rows.Scan(&nextDef.ID, &nextDef.Name, &nextDef.Description, &nextDef.Commands[0], &nextDef.Commands[1], &nextDef.Traversable, &nextDef.Unknown)
		world.Boundarys = append(world.Boundarys, nextDef)
	}
}

//LoadItemDefinitions Loads game item data into memory for quick access.
func LoadItemDefinitions() {
	database := Open(config.WorldDB())
	defer database.Close()
	rows, err := database.Query("SELECT id, name, description, command, base_price, stackable, special, members FROM `items` ORDER BY id")
	if err != nil {
		log.Error.Println("Couldn't load SQLite3 database:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		nextDef := world.ItemDefinition{}
		rows.Scan(&nextDef.ID, &nextDef.Name, &nextDef.Description, &nextDef.Command, &nextDef.BasePrice, &nextDef.Stackable, &nextDef.Quest, &nextDef.Members)
		world.ItemDefs = append(world.ItemDefs, nextDef)
	}
}

//LoadNpcDefinitions Loads game NPC data into memory for quick access.
func LoadNpcDefinitions() {
	database := Open(config.WorldDB())
	defer database.Close()
	rows, err := database.Query("SELECT id, name, description, command, hits, attack, strength, defense, attackable FROM `npcs` ORDER BY id")
	if err != nil {
		log.Error.Println("Couldn't load SQLite3 database:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		nextDef := world.NpcDefinition{}
		rows.Scan(&nextDef.ID, &nextDef.Name, &nextDef.Description, &nextDef.Command, &nextDef.Hits, &nextDef.Attack, &nextDef.Strength, &nextDef.Defense, &nextDef.Attackable)
		world.NpcDefs = append(world.NpcDefs, nextDef)
	}
}

//LoadObjectLocations Loads the game objects into memory from the SQLite3 database.
func LoadObjectLocations() {
	objectCounter := 0
	database := Open(config.WorldDB())
	defer database.Close()
	rows, err := database.Query("SELECT `id`, `direction`, `boundary`, `x`, `y` FROM `game_object_locations`")
	if err != nil {
		log.Error.Println("Couldn't load SQLite3 database:", err)
		return
	}
	defer rows.Close()
	var id, direction, boundary, x, y int
	for rows.Next() {
		rows.Scan(&id, &direction, &boundary, &x, &y)
		if world.GetObject(x, y) != nil {
			continue
		}
		objectCounter++
		world.AddObject(world.NewObject(id, direction, x, y, boundary != 0))
	}
}

//LoadNpcLocations Loads the game objects into memory from the SQLite3 database.
func LoadNpcLocations() {
	npcCounter := 0
	database := Open(config.WorldDB())
	defer database.Close()
	rows, err := database.Query("SELECT `id`, `startX`, `minX`, `maxX`, `startY`, `minY`, `maxY` FROM `npc_locations`")
	if err != nil {
		log.Error.Println("Couldn't load SQLite3 database:", err)
		return
	}
	defer rows.Close()
	var id, startX, minX, maxX, startY, minY, maxY int
	for rows.Next() {
		rows.Scan(&id, &startX, &minX, &maxX, &startY, &minY, &maxY)
		npcCounter++
		world.AddNpc(world.NewNpc(id, startX, startY, minX, maxX, minY, maxY))
	}
}

//SaveObjectLocations Clears world.db game object locations and repopulates it with the current server locations.
func SaveObjectLocations() int {
	database := Open(config.WorldDB())
	defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Info.Println("Error starting transaction for saving object locations:", err)
		return -1
	}

	stmt, err := tx.Exec("DELETE FROM game_object_locations")
	if err != nil {
		tx.Rollback()
		log.Info.Println("Error clearing object locations to save new ones:", err)
		return -1
	}
	if count, err := stmt.RowsAffected(); count < 1 || err != nil {
		if err != nil {
			log.Warning.Println("Error inserting new game object location to world.db:", err)
			return -1
		}
		log.Warning.Printf("Rows affected < 1 in game object location insert:%d\n", count)
		return -1
	}

	totalInserts := 0
	for _, v := range world.GetAllObjects() {
		stmt, err := tx.Exec("INSERT INTO game_object_locations(id, direction, x, y, boundary) VALUES(?, ?, ?, ?, ?)", v.ID, v.Direction, v.X(), v.Y(), v.Boundary)
		if err != nil {
			log.Warning.Println("Error inserting game object location to database:", err)
			continue
		}
		if count, err := stmt.RowsAffected(); count < 1 || err != nil {
			if err != nil {
				log.Warning.Println("Error inserting new game object location to world.db:", err)
				continue
			}
			log.Warning.Printf("Rows affected < 1 in game object location insert:%d\n", count)
			continue
		}
		totalInserts++
	}

	if err := tx.Commit(); err != nil {
		log.Warning.Println("Couldn't commit game object locations:", err)
		return -1
	}

	return totalInserts
}
