package server

import (
	"database/sql"

	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/list"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"

	// Necessary for sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

//Objects List of the game objects in the world
var Objects = list.New(16384)

//WorldDatabase SQLite3 connection reference for world data.
var WorldDatabase *sql.DB

//PlayerDatabase SQLite3 connection reference for player data.
var PlayerDatabase *sql.DB

//LoadObjects Loads the game objects into memory from the SQLite3 database.
func LoadObjects() int {
	if WorldDatabase == nil {
		WorldDatabase = Database("world.db")
	}
	rows, err := WorldDatabase.Query("SELECT `id`, `direction`, `type`, `x`, `y` FROM `game_object_locations`")
	if err != nil {
		LogError.Println("Couldn't load SQLite3 database:", err)
		return 0
	}
	var id, direction, kind, x, y int
	counter := 0
	for rows.Next() {
		rows.Scan(&id, &direction, &kind, &x, &y)
		o := entity.NewObject(id, direction, x, y, kind != 0)
		o.Index = Objects.Add(o)
		entity.GetRegion(x, y).AddObject(o)
	}
	return counter
}

//Database Returns an active sqlite3 database reference for the specified database file.
func Database(file string) *sql.DB {
	database, err := sql.Open("sqlite3", TomlConfig.DataDir+file)
	if err != nil {
		LogError.Println("Couldn't load SQLite3 database:", err)
		return nil
	}
	return database
}

//LoadPlayer Loads a player from the SQLite3 database, returns a login response code.
func (c *Client) LoadPlayer(username string, password string) int {
	if PlayerDatabase == nil {
		PlayerDatabase = Database("players.db")
	}

	stmt, err := PlayerDatabase.Prepare("SELECT `player`.`id`, `player`.`x`, `player`.`y`, `player`.`rank`, `player`.`fightmode`, `player`.`lastlogin`, `player`.`lastip`, `player`.`lastskulled`, `player`.`changingappearance`, `player`.`male`, `player`.`fatigue`, `appearance`.`haircolour`, `appearance`.`topcolour`, `appearance`.`trousercolour`, `appearance`.`skincolour`, `appearance`.`head`, `appearance`.`body` FROM `player` INNER JOIN `appearance` ON `appearance`.`playerid` = `player`.`id` WHERE `player`.`username`=? COLLATE NOCASE AND `player`.`password`=?")
	if err != nil {
		LogInfo.Println("LoadPlayer(string,string): Could not prepare query statement for player:", err)
		return 9
	}
	rows, err := stmt.Query(username, password)
	if err != nil {
		LogInfo.Println("LoadPlayer(string,string): Could not execute query statement for player:", err)
		return 9
	}
	var x, y, fightmode int
	if !rows.Next() {
		return 3
	}
	Clients[strutil.Base37(username)] = c
	rows.Scan(&c.player.DatabaseIndex, &x, &y, &c.player.Rank, &fightmode)
	c.player.SetCoords(x, y)
	c.player.SetFightMode(fightmode)
	return 0
}
