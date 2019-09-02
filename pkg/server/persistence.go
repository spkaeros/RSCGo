package server

import (
	"database/sql"

	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/list"
	"bitbucket.org/zlacki/rscgo/pkg/server/errors"

	// Necessary for sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

//Objects List of the game objects in the world
var Objects = list.New(16384)

//LoadObjects Loads the game objects into memory from the SQLite3 database.
func LoadObjects() bool {
	database := OpenDatabase(TomlConfig.Database.WorldDB)
	defer database.Close()
	rows, err := database.Query("SELECT `id`, `direction`, `type`, `x`, `y` FROM `game_object_locations`")
	defer rows.Close()
	if err != nil {
		LogError.Println("Couldn't load SQLite3 database:", err)
		return false
	}
	var id, direction, kind, x, y int
	for rows.Next() {
		rows.Scan(&id, &direction, &kind, &x, &y)
		o := entity.NewObject(id, direction, x, y, kind != 0)
		o.Index = Objects.Add(o)
		entity.GetRegion(x, y).AddObject(o)
	}
	return true
}

//OpenDatabase Returns an active sqlite3 database reference for the specified database file.
func OpenDatabase(file string) *sql.DB {
	database, err := sql.Open("sqlite3", "file:"+TomlConfig.DataDir+file)
	if err != nil {
		LogError.Println("Couldn't load SQLite3 database:", err)
		return nil
	}
	database.SetMaxOpenConns(1)
	return database
}

//LoadPlayer Loads a player from the SQLite3 database, returns a login response code.
func (c *Client) LoadPlayer(usernameHash uint64, password string) int {
	playerID, err := ValidatePlayer(c.player, usernameHash, password)
	if playerID < 0 || err != nil {
		if err.Error() == "Could not find player" {
			// Invalid username/password
			return 3
		}
		// Database error
		return 9
	}
	if i := PlayerAppearance(c.player); i != 0 {
		return i
	}

	Clients[usernameHash] = c
	return 0
}

//ValidatePlayer Sets the player's essential persistent variables from player table from base37 username and password hash.
// Returns 0 if successful, login response code otherwise.
func ValidatePlayer(player *entity.Player, hash uint64, password string) (int, error) {
	database := OpenDatabase(TomlConfig.Database.PlayerDB)
	defer database.Close()

	stmt, err := database.Prepare("SELECT id, x, y, group_id FROM player2 WHERE userhash=? AND password=?")
	defer stmt.Close()
	if err != nil {
		LogInfo.Println("ValidatePlayer(uint64,string): Could not prepare query statement for player:", err)
		return -1, errors.NewDatabaseError(err.Error())
	}
	rows, err := stmt.Query(hash, password)
	defer rows.Close()
	if err != nil {
		LogInfo.Println("ValidatePlayer(uint64,string): Could not execute query statement for player:", err)
		return -1, errors.NewDatabaseError(err.Error())
	}
	if !rows.Next() {
		return -1, errors.NewDatabaseError("Could not find player")
	}
	rows.Scan(&player.DatabaseIndex, &player.Location().X, &player.Location().Y, &player.Rank)
	return player.DatabaseIndex, nil
}

//PlayerAppearance Sets the player's appearance variables from a database search by the player's DatabaseIndex.
// Returns 0 if successful, login response code otherwise.
func PlayerAppearance(player *entity.Player) int {
	database := OpenDatabase(TomlConfig.Database.PlayerDB)
	defer database.Close()
	stmt, err := database.Prepare("SELECT haircolour, topcolour, trousercolour, skincolour, head, body FROM appearance WHERE playerid=?")
	defer stmt.Close()
	if err != nil {
		LogInfo.Println("LoadPlayer(uint64,string): Could not prepare query statement for player appearance:", err)
		return 9
	}
	rows, err := stmt.Query(player.DatabaseIndex)
	defer rows.Close()
	if err != nil {
		LogInfo.Println("LoadPlayer(uint64,string): Could not execute query statement for player appearance:", err)
		return 9
	}
	if !rows.Next() {
		return 17
	}
	rows.Scan(&player.Appearance.Hair, &player.Appearance.Top, &player.Appearance.Bottom, &player.Appearance.Skin, &player.Appearance.Head, &player.Appearance.Body)
	return 0
}

//Save Saves a player to the SQLite3 database.
func (c *Client) Save() {
	db := OpenDatabase(TomlConfig.Database.PlayerDB)
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		LogInfo.Println("Save(): Could not begin transcaction for player update.")
		return
	}
	rs, err := tx.Exec("UPDATE player2 SET x=?, y=? WHERE id=?", c.player.X(), c.player.Y(), c.player.DatabaseIndex)
	count, err := rs.RowsAffected()
	if err != nil {
		LogWarning.Println("Save(): UPDATE failed for player:", err)
		if err := tx.Rollback(); err != nil {
			LogWarning.Println("Save(): Transaction rollback failed:", err)
		}
		return
	}
	if count <= 0 {
		LogInfo.Println("Save(): Affected nothing for player update!")
	}
	if err := tx.Commit(); err != nil {
		LogWarning.Println("Save(): Error committing transaction for player update:", err)
	}
}
