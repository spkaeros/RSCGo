package server

import (
	"database/sql"
	"os"

	// Necessary for sqlite3 driver
	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/list"
	_ "github.com/mattn/go-sqlite3"
)

//Objects List of the game objects in the world
var Objects = list.New(16384)

//LoadObjects Loads the game objects into memory from the SQLite3 database.
func LoadObjects() bool {
	database, err := sql.Open("sqlite3", DataDirectory+string(os.PathSeparator)+"world.db")
	if err != nil {
		LogError.Println("Couldn't load SQLite3 database:", err)
		return false
	}
	rows, err := database.Query("SELECT `id`, `direction`, `type`, `x`, `y` FROM `game_object_locations`")
	if err != nil {
		LogError.Println("Couldn't load SQLite3 database:", err)
		return false
	}
	var id, direction, kind, x, y int
	counter := 0
	for rows.Next() {
		rows.Scan(&id, &direction, &kind, &x, &y)
		o := entity.NewObject(id, direction, x, y, kind == 1)
		o.Index = counter
		counter++
		objects := entity.GetRegion(x, y).Objects
		entity.GetRegion(x, y).Objects = append(objects, o)
	}
	return true
}
