package server

import (
	"database/sql"
	"strconv"

	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/server/errors"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"

	// Necessary for sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var objectDefCmd = make(map[int]string)
var ObjectDefinitions []ObjectDefinition

//ObjectDefinition This represents a single definition for a single object in the game.
type ObjectDefinition struct {
	ID            int
	Name          string
	Commands      []string
	Description   string
	Type          int
	Width, Height int
	Length        int
}

//LoadObjectDefinitions Loads game object data into memory for quick access.
func LoadObjectDefinitions() {
	database := OpenDatabase(TomlConfig.Database.WorldDB)
	defer database.Close()
	rows, err := database.Query("SELECT id, name, description, command_one, command_two, type, width, height, ground_item_var FROM `game_objects`")
	defer rows.Close()
	if err != nil {
		LogError.Println("Couldn't load SQLite3 database:", err)
		return
	}
	for rows.Next() {
		nextDef := ObjectDefinition{Commands: make([]string, 2)}
		rows.Scan(&nextDef.ID, &nextDef.Name, &nextDef.Description, &nextDef.Commands[0], &nextDef.Commands[1], &nextDef.Type, &nextDef.Width, &nextDef.Height, &nextDef.Length)
		ObjectDefinitions = append(ObjectDefinitions, nextDef)
	}
}

//LoadObjects Loads the game objects into memory from the SQLite3 database.
func LoadObjects() int {
	objectCounter := 0
	database := OpenDatabase(TomlConfig.Database.WorldDB)
	defer database.Close()
	rows, err := database.Query("SELECT `id`, `direction`, `type`, `x`, `y` FROM `game_object_locations`")
	defer rows.Close()
	if err != nil {
		LogError.Println("Couldn't load SQLite3 database:", err)
		return 0
	}
	var id, direction, kind, x, y int
	for rows.Next() {
		rows.Scan(&id, &direction, &kind, &x, &y)
		o := entity.NewObject(id, direction, x, y, kind != 0)
		o.SetIndex(objectCounter)
		objectCounter++
		entity.GetRegion(x, y).AddObject(o)
	}
	return objectCounter
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
func (c *Client) LoadPlayer(usernameHash uint64, password string, loginReply chan byte) {
	validateCredentials := func() error {

		database := OpenDatabase(TomlConfig.Database.PlayerDB)
		defer database.Close()

		stmt, err := database.Prepare("SELECT player.id, player.x, player.y, player.group_id, appearance.haircolour, appearance.topcolour, appearance.trousercolour, appearance.skincolour, appearance.head, appearance.body FROM player2 AS player INNER JOIN appearance AS appearance WHERE appearance.playerid=player.id AND player.userhash=? AND player.password=?")
		defer stmt.Close()
		if err != nil {
			LogInfo.Println("ValidatePlayer(uint64,string): Could not prepare query statement for player:", err)
			loginReply <- byte(3)
			return errors.NewDatabaseError(err.Error())
		}
		rows, err := stmt.Query(usernameHash, password)
		defer rows.Close()
		if err != nil {
			LogInfo.Println("ValidatePlayer(uint64,string): Could not execute query statement for player:", err)
			loginReply <- byte(8)
			return errors.NewDatabaseError(err.Error())
		}
		if !rows.Next() {
			loginReply <- byte(3)
			return errors.NewDatabaseError("Could not find player")
		}
		rows.Scan(&c.player.DatabaseIndex, &c.player.X, &c.player.Y, &c.player.Rank, &c.player.Appearance.Hair, &c.player.Appearance.Top, &c.player.Appearance.Bottom, &c.player.Appearance.Skin, &c.player.Appearance.Head, &c.player.Appearance.Body)
		return nil
	}
	if err := validateCredentials(); err != nil {
		return
	}
	/*	if err := PlayerAppearance(c.player); err != nil {
		return
	}*/
	if err := PlayerAttributes(c.player); err != nil {
		return
	}
	if err := PlayerFriends(c.player); err != nil {
		return
	}
	if err := PlayerIgnore(c.player); err != nil {
		return
	}

	c.player.UserBase37 = usernameHash
	c.player.Username = strutil.DecodeBase37(usernameHash)
	c.player.SetIndex(c.Index)
	Clients[usernameHash] = c
	loginReply <- byte(0)
	return
}

//ValidatePlayer Sets the player's essential persistent variables from player table from base37 username and password hash.
// Returns 0 if successful, login response code otherwise.
func ValidatePlayer(player *entity.Player, hash uint64, password string) error {
	database := OpenDatabase(TomlConfig.Database.PlayerDB)
	defer database.Close()

	stmt, err := database.Prepare("SELECT id, x, y, group_id FROM player2 WHERE userhash=? AND password=?")
	defer stmt.Close()
	if err != nil {
		LogInfo.Println("ValidatePlayer(uint64,string): Could not prepare query statement for player:", err)
		return errors.NewDatabaseError(err.Error())
	}
	rows, err := stmt.Query(hash, password)
	defer rows.Close()
	if err != nil {
		LogInfo.Println("ValidatePlayer(uint64,string): Could not execute query statement for player:", err)
		return errors.NewDatabaseError(err.Error())
	}
	if !rows.Next() {
		return errors.NewDatabaseError("Could not find player")
	}
	rows.Scan(&player.DatabaseIndex, &player.X, &player.Y, &player.Rank)
	return nil
}

//PlayerAppearance Sets the player's appearance variables from a database search by the player's DatabaseIndex.
// Returns nil upon success.
func PlayerAppearance(player *entity.Player) error {
	database := OpenDatabase(TomlConfig.Database.PlayerDB)
	defer database.Close()
	stmt, err := database.Prepare("SELECT haircolour, topcolour, trousercolour, skincolour, head, body FROM appearance WHERE playerid=?")
	defer stmt.Close()
	if err != nil {
		LogInfo.Println("LoadPlayer(uint64,string): Could not prepare query statement for player appearance:", err)
		return errors.NewDatabaseError("Statement could not be prepared.")
	}
	rows, err := stmt.Query(player.DatabaseIndex)
	defer rows.Close()
	if err != nil {
		LogInfo.Println("LoadPlayer(uint64,string): Could not execute query statement for player appearance:", err)
		return errors.NewDatabaseError("Statement could not execute.")
	}
	if !rows.Next() {
		return errors.NewDatabaseError("Could not find player")
	}
	rows.Scan(&player.Appearance.Hair, &player.Appearance.Top, &player.Appearance.Bottom, &player.Appearance.Skin, &player.Appearance.Head, &player.Appearance.Body)
	return nil
}

//PlayerAttributes Sets the player's attribute variables from a database search by the player's DatabaseIndex.
func PlayerAttributes(player *entity.Player) error {
	database := OpenDatabase(TomlConfig.Database.PlayerDB)
	defer database.Close()
	stmt, err := database.Prepare("SELECT name, value FROM player_attr WHERE player_id=?")
	defer stmt.Close()
	if err != nil {
		LogInfo.Println("LoadPlayer(uint64,string): Could not prepare query statement for player attributes:", err)
		return errors.NewDatabaseError("Statement could not be prepared.")
	}
	rows, err := stmt.Query(player.DatabaseIndex)
	defer rows.Close()
	if err != nil {
		LogInfo.Println("LoadPlayer(uint64,string): Could not execute query statement for player attributes:", err)
		return errors.NewDatabaseError("Statement could not execute.")
	}
	for rows.Next() {
		var name, value string
		rows.Scan(&name, &value)
		switch value[0] {
		case 'i':
			val, err := strconv.ParseInt(value[1:], 10, 64)
			if err != nil {
				LogInfo.Printf("Error loading int attribute[%v]: value=%v\n", name, value[1:])
				LogInfo.Println(err)
			}
			player.Attributes[entity.Attribute(name)] = int(val)
			break
		case 'l':
			val, err := strconv.ParseUint(value[1:], 10, 64)
			if err != nil {
				LogInfo.Printf("Error loading long int attribute[%v]: value=%v\n", name, value[1:])
				LogInfo.Println(err)
			}
			player.Attributes[entity.Attribute(name)] = uint(val)
			break
		case 'b':
			val, err := strconv.ParseBool(value[1:])
			if err != nil {
				LogInfo.Printf("Error loading boolean attribute[%v]: value=%v\n", name, value[1:])
				LogInfo.Println(err)
			}
			player.Attributes[entity.Attribute(name)] = val
			break
		}
	}
	return nil
}

//PlayerFriends Loads the player's friends list
func PlayerFriends(player *entity.Player) error {
	database := OpenDatabase(TomlConfig.Database.PlayerDB)
	defer database.Close()
	stmt, err := database.Prepare("SELECT playerhash FROM playerlist WHERE playerid=? AND `type`='friend'")
	defer stmt.Close()
	if err != nil {
		LogInfo.Println("LoadPlayer(uint64,string): Could not prepare query statement for player friends:", err)
		return errors.NewDatabaseError("Statement could not be prepared.")
	}
	rows, err := stmt.Query(player.DatabaseIndex)
	defer rows.Close()
	if err != nil {
		LogInfo.Println("LoadPlayer(uint64,string): Could not execute query statement for player friends:", err)
		return errors.NewDatabaseError("Statement could not execute.")
	}
	for rows.Next() {
		var hash uint64
		rows.Scan(&hash)
		player.FriendList[hash] = ClientFromHash(hash) != nil
		//		player.FriendList = append(player.FriendList, hash)
	}
	return nil
}

//PlayerIgnore Loads the player's ignore list
func PlayerIgnore(player *entity.Player) error {
	database := OpenDatabase(TomlConfig.Database.PlayerDB)
	defer database.Close()
	stmt, err := database.Prepare("SELECT playerhash FROM playerlist WHERE playerid=? AND `type`='ignore'")
	defer stmt.Close()
	if err != nil {
		LogInfo.Println("LoadPlayer(uint64,string): Could not prepare query statement for player ignore:", err)
		return errors.NewDatabaseError("Statement could not be prepared.")
	}
	rows, err := stmt.Query(player.DatabaseIndex)
	defer rows.Close()
	if err != nil {
		LogInfo.Println("LoadPlayer(uint64,string): Could not execute query statement for player ignores:", err)
		return errors.NewDatabaseError("Statement could not execute.")
	}
	for rows.Next() {
		var hash uint64
		rows.Scan(&hash)
		player.IgnoreList = append(player.IgnoreList, hash)
	}
	return nil
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
	saveLocation := func() {
		rs, err := tx.Exec("UPDATE player2 SET x=?, y=? WHERE id=?", c.player.X, c.player.Y, c.player.DatabaseIndex)
		count, err := rs.RowsAffected()
		if err != nil {
			LogWarning.Println("Save(): UPDATE failed for player location:", err)
			if err := tx.Rollback(); err != nil {
				LogWarning.Println("Save(): Transaction location rollback failed:", err)
			}
			return
		}

		if count <= 0 {
			LogInfo.Println("Save(): Affected nothing for location update!")
		}
	}
	saveAppearance := func() {
		// TODO: Should this just be attributes too??  Is that abusing the attributes table?
		appearance := c.player.Appearance
		rs, _ := tx.Exec("UPDATE appearance SET haircolour=?, topcolour=?, trousercolour=?, skincolour=?, head=?, body=? WHERE playerid=?", appearance.Hair, appearance.Top, appearance.Bottom, appearance.Skin, appearance.Head, appearance.Body, c.player.DatabaseIndex)
		count, err := rs.RowsAffected()
		if err != nil {
			LogWarning.Println("Save(): UPDATE failed for player appearance:", err)
			if err := tx.Rollback(); err != nil {
				LogWarning.Println("Save(): Transaction appearance rollback failed:", err)
			}
			return
		}

		if count <= 0 {
			LogInfo.Println("Save(): Affected nothing for appearance update!")
		}
	}
	clearAttributes := func() {
		if _, err := tx.Exec("DELETE FROM player_attr WHERE player_id=?", c.player.DatabaseIndex); err != nil {
			LogWarning.Println("Save(): DELETE failed for player attribute:", err)
			if err := tx.Rollback(); err != nil {
				LogWarning.Println("Save(): Transaction delete attributes rollback failed:", err)
			}
			return
		}
	}
	insertAttribute := func(name string, value interface{}) {
		var val string
		switch value.(type) {
		case int:
			val = "i" + strconv.FormatInt(int64(value.(int)), 10)
		case uint:
			val = "l" + strconv.FormatUint(uint64(value.(uint)), 10)
		case bool:
			if v, ok := value.(bool); v && ok {
				val = "b1"
			} else {
				val = "b0"
			}
		}
		rs, _ := tx.Exec("INSERT INTO player_attr(player_id, name, value) VALUES(?, ?, ?)", c.player.DatabaseIndex, name, val)
		count, err := rs.RowsAffected()
		if err != nil {
			LogWarning.Println("Save(): INSERT failed for player attribute:", err)
			if err := tx.Rollback(); err != nil {
				LogWarning.Println("Save(): Transaction insert attribute rollback failed:", err)
			}
			return
		}

		if count <= 0 {
			LogInfo.Println("Save(): Affected nothing for attribute insertion!")
		}
	}
	clearContactList := func(contactType string) {
		if _, err := tx.Exec("DELETE FROM playerlist WHERE playerid=? AND type=?", c.player.DatabaseIndex, contactType); err != nil {
			LogWarning.Println("Save(): DELETE failed for player friends:", err)
			if err := tx.Rollback(); err != nil {
				LogWarning.Println("Save(): Transaction delete friends rollback failed:", err)
			}
			return
		}
	}
	insertContactList := func(contactType string, hash uint64) {
		rs, _ := tx.Exec("INSERT INTO playerlist(playerid, playerhash, type) VALUES(?, ?, ?)", c.player.DatabaseIndex, hash, contactType)
		count, err := rs.RowsAffected()
		if err != nil {
			LogWarning.Println("Save(): INSERT failed for player friends:", err)
			if err := tx.Rollback(); err != nil {
				LogWarning.Println("Save(): Transaction insert friend rollback failed:", err)
			}
			return
		}

		if count <= 0 {
			LogInfo.Println("Save(): Affected nothing for friend insertion!")
		}
	}
	saveLocation()
	saveAppearance()
	clearAttributes()
	for name, value := range c.player.Attributes {
		insertAttribute(string(name), value)
	}
	clearContactList("friend")
	clearContactList("ignore")
	for hash := range c.player.FriendList {
		insertContactList("friend", hash)
	}
	for _, hash := range c.player.IgnoreList {
		insertContactList("ignore", hash)
	}

	if err := tx.Commit(); err != nil {
		LogWarning.Println("Save(): Error committing transaction for player update:", err)
	}
}
