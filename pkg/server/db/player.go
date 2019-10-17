package db

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/config"
	"bitbucket.org/zlacki/rscgo/pkg/server/errors"
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
	"bitbucket.org/zlacki/rscgo/pkg/strutil"
)

//CreatePlayer Creates a new entry in the player SQLite3 database with the specified credentials.
// Returns true if successful, otherwise returns false.
func CreatePlayer(username, password string) bool {
	database := Open(config.PlayerDB())
	defer database.Close()

	tx, err := database.Begin()
	if err != nil {
		log.Info.Println("CreatePlayer(): Could not begin transaction for new player.")
		return false
	}

	s, err := tx.Exec("INSERT INTO player2(username, userhash, password, x, y, group_id) VALUES(?, ?, ?, 220, 445, 0)", username, strutil.Base37(username), password)
	if err != nil {
		log.Info.Println("CreatePlayer(): Could not insert new player profile information:", err)
		return false
	}
	playerID, err := s.LastInsertId()
	if err != nil || playerID < 0 {
		log.Info.Printf("CreatePlayer(): Could not retrieve player database ID(got %d):\n%v", playerID, err)
		return false
	}
	_, err = tx.Exec("INSERT INTO appearance VALUES(?, 2, 8, 14, 0, 1, 2)", playerID)
	if err != nil {
		log.Info.Println("CreatePlayer(): Could not insert new player profile information:", err)
		return false
	}
	if err := tx.Commit(); err != nil {
		log.Warning.Println("CreatePlayer(): Error committing transaction for new player:", err)
		return false
	}

	return true
}

//UsernameExists Returns true if there is a player with the name 'username' in the player database, otherwise returns false.
func UsernameExists(username string) bool {
	database := Open(config.PlayerDB())
	defer database.Close()
	s, err := database.Query("SELECT id FROM player2 WHERE userhash=?", strutil.Base37(username))
	defer s.Close()
	if err != nil {
		log.Info.Println("UsernameTaken: Could not query player profile information:", err)
		// return true just to be safe since we could not check
		return true
	}
	return s.Next()
}

//ValidateCredentials Looks for a player with the specified credentials in the player database.  Returns nil if it finds the player, otherwise returns an error.
func ValidateCredentials(usernameHash uint64, password string, loginReply chan byte, player *world.Player) error {
	database := Open(config.PlayerDB())
	defer database.Close()

	rows, err := database.Query("SELECT player.id, player.x, player.y, player.group_id, appearance.haircolour, appearance.topcolour, appearance.trousercolour, appearance.skincolour, appearance.head, appearance.body FROM player2 AS player INNER JOIN appearance AS appearance WHERE appearance.playerid=player.id AND player.userhash=? AND player.password=?", usernameHash, password)
	defer rows.Close()
	if err != nil {
		log.Info.Println("ValidatePlayer(uint64,string): Could not prepare query statement for player:", err)
		loginReply <- byte(3)
		return errors.NewDatabaseError(err.Error())
	}
	if !rows.Next() {
		loginReply <- byte(3)
		return errors.NewDatabaseError("Could not find player")
	}
	rows.Scan(&player.DatabaseIndex, &player.X, &player.Y, &player.Rank, &player.Appearance.Hair, &player.Appearance.Top, &player.Appearance.Bottom, &player.Appearance.Skin, &player.Appearance.Head, &player.Appearance.Body)
	return nil
}
