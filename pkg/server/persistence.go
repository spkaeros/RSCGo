package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/config"
	"bitbucket.org/zlacki/rscgo/pkg/server/db"
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"

	// Necessary for sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

//LoadObjectLocations Loads the game objects into memory from the SQLite3 database.
func LoadObjectLocations() int {
	objectCounter := 0
	database := db.Open(config.WorldDB())
	defer database.Close()
	rows, err := database.Query("SELECT `id`, `direction`, `type`, `x`, `y` FROM `game_object_locations`")
	defer rows.Close()
	if err != nil {
		log.Error.Println("Couldn't load SQLite3 database:", err)
		return 0
	}
	var id, direction, boundary, x, y int
	for rows.Next() {
		rows.Scan(&id, &direction, &boundary, &x, &y)
		if world.GetObject(x, y) != nil {
			continue
		}
		objectCounter++
		world.AddObject(world.NewObject(id, direction, x, y, boundary != 0))
	}
	return objectCounter
}

//SaveObjectLocations Clears world.db game object locations and repopulates it with the current server locations.
func SaveObjectLocations() int {
	database := db.Open(config.WorldDB())
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
		stmt, err := tx.Exec("INSERT INTO game_object_locations(id, direction, x, y, type) VALUES(?, ?, ?, ?, ?)", v.ID, v.Direction, v.X, v.Y, v.Boundary)
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
