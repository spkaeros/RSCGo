package db

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/spkaeros/rscgo/pkg/crypto"
	"github.com/spkaeros/rscgo/pkg/game/entity"

	"github.com/spkaeros/rscgo/pkg/config"
	"github.com/spkaeros/rscgo/pkg/errors"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

//PlayerService An interface for manipulating player save data.
type PlayerService interface {
	PlayerCreate(string, string, string) bool
	PlayerNameExists(username string) bool
	PlayerHasRecoverys(uint64) bool
	PlayerValidLogin(uint64, string) bool
	PlayerChangePassword(uint64, string) bool
	PlayerLoadRecoverys(uint64) []string
	PlayerLoad(*world.Player) bool
	PlayerSave(*world.Player)
}

//NewPlayerServiceSql Returns a new SqlPlayerService to manage the specified *sql.DB instance, configured against
// the default players database.
func NewPlayerServiceSql() PlayerService {
	s := newSqlService(config.PlayerDriver())
	s.sqlOpen(config.PlayerDB())
	return s
}

//DefaultPlayerService the default player save managing service in use by the game server
// Currently using an sqlService.
var DefaultPlayerService PlayerService

//PlayerCreate Creates a new entry in the player SQLite3 database with the specified credentials.
// Returns true if successful, otherwise returns false.
func (s *sqlService) PlayerCreate(username, password, ip string) bool {
	db := s.connect(context.Background())
	defer db.Close()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Info.Println("SQLiteService Could not begin transaction:", err)
		return false
	}

	var playerID int
	if s.Driver != "postgres" {
		stmt, err := tx.Exec("INSERT INTO player(username, userhash, password, x, y, group_id) VALUES($1, $2, $3, 220, 445, 0)", username, strutil.Base37.Encode(username), crypto.Hash(password))
		if err != nil {
			log.Info.Println("SQLiteService Could not insert new player profile information:", err)
			return false
		}
		pID, err := stmt.LastInsertId()
		if err != nil || playerID < 0 {
			tx.Rollback()
			log.Info.Printf("PlayerCreate(): Could not retrieve player database ID(%d):\n%v", playerID, err)
			return false
		}
		playerID = int(pID)
	} else {
		stmt := tx.QueryRow("INSERT INTO player(username, userhash, password, x, y, group_id) VALUES($1, $2, $3, 220, 445, 0) RETURNING id", username, strutil.Base37.Encode(username), crypto.Hash(password))
		err = stmt.Scan(&playerID)
		if err != nil || playerID < 0 {
			tx.Rollback()
			log.Info.Printf("PlayerCreate(): Could not retrieve player database ID(%d):\n%v", playerID, err)
			return false
		}
	}
	
	_, err = tx.Exec("INSERT INTO appearance VALUES($1, 2, 8, 14, 0, 1, 2)", playerID)
	if err != nil {
		log.Info.Println("PlayerCreate(): Could not insert new player profile information:", err)
		return false
	}
	_, err = tx.Exec("INSERT INTO player_attr VALUES($1, 'lastIP', $2)", playerID, "s" + ip)
	if err != nil {
		log.Info.Println("PlayerCreate(): Could not insert new player profile information:", err)
		return false
	}
	if err := tx.Commit(); err != nil {
		log.Warning.Println("PlayerCreate(): Error committing transaction for new player:", err)
		return false
	}

	return true
}

//PlayerNameExists Returns true if there is a player with the name 'username' in the player database, otherwise returns false.
func (s *sqlService) PlayerNameExists(username string) bool {
	database := s.connect(context.Background())
	defer database.Close()
	stmt, err := database.QueryContext(context.Background(), "SELECT id FROM player WHERE userhash=$1", strutil.Base37.Encode(username))
	if err != nil {
		log.Info.Println("UsernameTaken: Could not query player profile information:", err)
		// return true just to be safe since we could not check
		return true
	}
	defer stmt.Close()
	return stmt.Next()
}

//PlayerValidLogin Returns true if it finds a user with this username hash and password in the database, otherwise returns false
func (s *sqlService) PlayerValidLogin(userHash uint64, password string) bool {
	database := s.connect(context.Background())
	defer database.Close()
	rows, err := database.QueryContext(context.Background(), "SELECT id FROM player WHERE userhash=$1 AND password=$2", userHash, password)
	if err != nil {
		log.Info.Println("Validate: Could not validate user credentials:", err)
		return false
	}
	defer rows.Close()
	return rows.Next()
}

//PlayerChangePassword Updates the players password to password in the database.
func (s *sqlService) PlayerChangePassword(userHash uint64, password string) bool {
	database := s.connect(context.Background())
	defer database.Close()
	stmt, err := database.ExecContext(context.Background(), "UPDATE player SET password=$1 WHERE userhash=$2", password, userHash)
	if err != nil {
		log.Info.Println("PlayerChangePassword: Could not update player password:", err)
		return false
	}
	count, err := stmt.RowsAffected()
	if count <= 0 || err != nil {
		log.Info.Println("PlayerChangePassword: Could not update player password:", err)
		return false
	}
	return true
}

//PlayerHasRecoverys Returns true if this username has recovery questions assigned to it, otherwise returns false.
func (s *sqlService) PlayerHasRecoverys(userHash uint64) bool {
	database := s.connect(context.Background())
	defer database.Close()
	rows, err := database.QueryContext(context.Background(), "SELECT question1 FROM recovery_questions WHERE userhash=$1", userHash)
	if err != nil {
		log.Info.Println("PlayerHasRecoverys: Could not search for recovery questions:", err)
		return false
	}
	defer rows.Close()
	return rows.Next()
}

//PlayerLoadRecoverys Retrieves the recovery questions assigned to this username if any, otherwise returns nil
func (s *sqlService) PlayerLoadRecoverys(userHash uint64) []string {
	database := s.connect(context.Background())
	defer database.Close()
	rows, err := database.QueryContext(context.Background(), "SELECT question1, question2, question3, question4, question5 FROM recovery_questions WHERE userhash=$1", userHash)
	if err != nil {
		log.Info.Println("PlayerLoadRecoverys: Could not find recovery questions:", err)
		return nil
	}
	defer rows.Close()

	var question1, question2, question3, question4, question5 string
	if rows.Next() {
		err := rows.Scan(&question1, &question2, &question3, &question4, &question5)
		if err != nil {
			log.Info.Println("PlayerLoadRecoverys: Could not scan recovery questions to variables:", err)
			return nil
		}
		return []string{question1, question2, question3, question4, question5}
	}

	return nil
}

//SaveRecoveryQuestions Saves new recovery questions to the database.
func (s *sqlService) SaveRecoveryQuestions(userHash uint64, questions []string, answers []uint64) {

}

//PlayerLoad Loads a player from the SQLite3 database, returns a login response code.
// Returns: true on success, false on failure
func (s *sqlService) PlayerLoad(player *world.Player) bool {
	loadProfile := func() error {
		database := s.connect(context.Background())
		defer database.Close()
		rows, err := database.QueryContext(context.Background(), "SELECT player.id, player.x, player.y, player.group_id, appearance.haircolour, appearance.topcolour, appearance.trousercolour, appearance.skincolour, appearance.head, appearance.body FROM player INNER JOIN appearance ON appearance.playerid=player.id AND player.userhash=$1", player.UsernameHash())
		if err != nil {
			log.Info.Println("Load error: Could not prepare statement:", err)
			return errors.NewDatabaseError(err.Error())
		}
		defer rows.Close()
		if !rows.Next() {
			return errors.NewDatabaseError("Could not find player")
		}
		var x, y, rank, dbID int
		rows.Scan(&dbID, &x, &y, &rank, &player.Appearance.HeadColor, &player.Appearance.BodyColor, &player.Appearance.LegsColor, &player.Appearance.SkinColor, &player.Appearance.Head, &player.Appearance.Body)
		//	player.Location = world.NewLocation(x, y)
		//	player.Teleport(220, 445)
		player.TransAttrs.SetVar("dbID", dbID)
		player.TransAttrs.SetVar("rank", rank)
		player.Equips[0] = player.Appearance.Head
		player.Equips[1] = player.Appearance.Body
		player.SetX(x)
		player.SetY(y)
		return nil
	}
	loadAttributes := func() error {
		database := s.connect(context.Background())
		defer database.Close()

		rows, err := database.QueryContext(context.Background(), "SELECT name, value FROM player_attr WHERE player_id=$1", player.DatabaseID())
		if err != nil {
			log.Info.Println("Load error: Could not prepare statement:", err)
			return errors.NewDatabaseError(err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			var name, value string
			rows.Scan(&name, &value)
			switch value[0] {
			case 'i':
				val, err := strconv.ParseInt(value[1:], 10, 64)
				if err != nil {
					log.Info.Printf("Error loading int attribute[%v]: value=%v\n", name, value[1:])
					log.Info.Println(err)
				}
				player.Attributes.SetVar(name, int(val))
				break
			case 'l':
				val, err := strconv.ParseUint(value[1:], 10, 64)
				if err != nil {
					log.Info.Printf("Error loading long int attribute[%v]: value=%v\n", name, value[1:])
					log.Info.Println(err)
				}
				player.Attributes.SetVar(name, uint(val))
				break
			case 'b':
				val, err := strconv.ParseBool(value[1:])
				if err != nil {
					log.Info.Printf("Error loading boolean attribute[%v]: value=%v\n", name, value[1:])
					log.Info.Println(err)
				}
				player.Attributes.SetVar(name, val)
				break
			case 's':
				player.Attributes.SetVar(name, value[1:])
				break
			case 'd':
				t, err := time.ParseDuration(value[1:])
				if err != nil {
					continue
				}
				player.Attributes.SetVar(name, time.Now().Add(t))
			case 't':
				t, err := time.ParseInLocation(time.RFC822, value[1:], time.Local)
				if err != nil {
					continue
				}
				player.Attributes.SetVar(name, t)
			}
		}
		return nil
	}
	loadContactList := func(list string) error {
		database := s.connect(context.Background())
		defer database.Close()

		rows, err := database.QueryContext(context.Background(), "SELECT playerhash FROM contacts WHERE playerid=$1 AND type=$2", player.DatabaseID(), list)
		if err != nil {
			log.Info.Println("Load error: Could not prepare statement:", err)
			return errors.NewDatabaseError(err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			var hash uint64
			rows.Scan(&hash)
			switch list {
			case "friend":
				player.FriendList.Add(strutil.Base37.Decode(hash))
			case "ignore":
				player.IgnoreList = append(player.IgnoreList, hash)
			}
		}
		return nil
	}
	loadInventory := func() error {
		database := s.connect(context.Background())
		defer database.Close()
		rows, err := database.QueryContext(context.Background(), "SELECT itemid, amount, wielded FROM inventory WHERE playerid=$1", player.DatabaseID())
		if err != nil {
			log.Info.Println("Load error: Could not prepare statement:", err)
			return errors.NewDatabaseError(err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			var id, amt int
			wielded := false
			rows.Scan(&id, &amt, &wielded)
			index := player.Inventory.Add(id, amt)
			if e := world.GetEquipmentDefinition(id); e != nil && wielded {
				player.Inventory.Get(index).Worn = true
				player.Equips[e.Position] = e.Sprite
				player.SetAimPoints(player.AimPoints() + e.Aim)
				player.SetPowerPoints(player.PowerPoints() + e.Power)
				player.SetArmourPoints(player.ArmourPoints() + e.Armour)
				player.SetMagicPoints(player.MagicPoints() + e.Magic)
				player.SetPrayerPoints(player.PrayerPoints() + e.Prayer)
				player.SetRangedPoints(player.RangedPoints() + e.Ranged)
			}
		}
		return nil
	}
	loadBank := func() error {
		database := s.connect(context.Background())
		defer database.Close()
		rows, err := database.QueryContext(context.Background(), "SELECT itemid, amount FROM bank WHERE playerid=$1", player.DatabaseID())
		if err != nil {
			log.Info.Println("Load error: Could not prepare statement:", err)
			return errors.NewDatabaseError(err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			var id, amt int
			rows.Scan(&id, &amt)
			player.Bank().Add(id, amt)
		}
		return nil
	}
	loadStats := func() error {
		database := s.connect(context.Background())
		defer database.Close()
		rows, err := database.QueryContext(context.Background(), "SELECT cur, exp FROM stats WHERE playerid=$1 ORDER BY num", player.DatabaseID())
		if err != nil {
			log.Info.Println("Load error: Could not prepare statement:", err)
			return errors.NewDatabaseError(err.Error())
		}
		defer rows.Close()
		i := 0
		for rows.Next() {
			var cur, exp int
			rows.Scan(&cur, &exp)
			player.Skills().SetCur(i, cur)
			player.Skills().SetMax(i, entity.ExperienceToLevel(exp))
			player.Skills().SetExp(i, exp)
			i++
		}
		return nil
	}

	// If this fails, then the login information was incorrect, and we don't need to do anything else
	if err := loadProfile(); err != nil {
		return false
	}
	if err := loadAttributes(); err != nil {
		return false
	}
	if err := loadContactList("friend"); err != nil {
		return false
	}
	if err := loadContactList("ignore"); err != nil {
		return false
	}
	if err := loadInventory(); err != nil {
		return false
	}
	if err := loadBank(); err != nil {
		return false
	}
	if err := loadStats(); err != nil {
		return false
	}
	return true
}

//PlayerSave Saves a player to the SQLite3 database.
func (s *sqlService) PlayerSave(player *world.Player) {
	db := s.connect(context.Background())
	defer db.Close()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Info.Println("Save(): Could not begin transcaction for player update.")
		return
	}
	updateLocation := func() {
		rs, err := tx.Exec("UPDATE player SET x=$1, y=$2 WHERE id=$3", player.X(), player.Y(), player.DatabaseID())
		if err != nil {
			log.Warning.Println("Save(): UPDATE failed for player location:", err)
			if err := tx.Rollback(); err != nil {
				log.Warning.Println("Save(): Transaction location rollback failed:", err)
			}
			return
		}
		count, _ := rs.RowsAffected()

		if count <= 0 {
			log.Info.Println("Save(): Affected nothing for location update!")
		}
	}
	updateAppearance := func() {
		// TODO: Should this just be attributes too??  Is that abusing the attributes table?
		appearance := player.Appearance
		rs, _ := tx.Exec("UPDATE appearance SET haircolour=$1, topcolour=$2, trousercolour=$3, skincolour=$4, head=$5, body=$6 WHERE playerid=$7", appearance.HeadColor, appearance.BodyColor, appearance.LegsColor, appearance.SkinColor, appearance.Head, appearance.Body, player.DatabaseID())
		count, err := rs.RowsAffected()
		if err != nil {
			log.Warning.Println("Save(): UPDATE failed for player appearance:", err)
			if err := tx.Rollback(); err != nil {
				log.Warning.Println("Save(): Transaction appearance rollback failed:", err)
			}
			return
		}

		if count <= 0 {
			log.Info.Println("Save(): Affected nothing for appearance update!")
		}
	}
	clearAttributes := func() {
		if _, err := tx.Exec("DELETE FROM player_attr WHERE player_id=$1", player.DatabaseID()); err != nil {
			log.Warning.Println("Save(): DELETE failed for player attribute:", err)
			if err := tx.Rollback(); err != nil {
				log.Warning.Println("Save(): Transaction delete attributes rollback failed:", err)
			}
			return
		}
	}
	insertAttribute := func(name string, value interface{}) bool {
		var val string
		switch value.(type) {
		case int64:
			val = "i" + strconv.FormatInt(value.(int64), 10)
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
		case string:
			if v, ok := value.(string); ok {
				val = "s" + v
			}
		case time.Time:
			if strings.HasSuffix(name, "Timer") {
				// Save timers as duration
				if v, ok := value.(time.Time); ok {
					val = "d" + time.Until(v).String()
				}
			} else {
				if v, ok := value.(time.Time); ok {
					val = "t" + v.Format(time.RFC822)
				}
			}
		}
		rs, _ := tx.Exec("INSERT INTO player_attr(player_id, name, value) VALUES($1, $2, $3)", player.DatabaseID(), name, val)
		count, err := rs.RowsAffected()
		if err != nil {
			log.Warning.Println("Save(): INSERT failed for player attribute:", err)
			if err := tx.Rollback(); err != nil {
				log.Warning.Println("Save(): Transaction insert attribute rollback failed:", err)
			}
			return false
		}

		if count <= 0 {
			log.Info.Println("Save(): Affected nothing for attribute insertion!")
		}
		return false
	}
	clearContactList := func(contactType string) {
		if _, err := tx.Exec("DELETE FROM contacts WHERE playerid=$1 AND type=$2", player.DatabaseID(), contactType); err != nil {
			log.Warning.Println("Save(): DELETE failed for player friends:", err)
			if err := tx.Rollback(); err != nil {
				log.Warning.Println("Save(): Transaction delete friends rollback failed:", err)
			}
			return
		}
	}
	insertContact := func(contactType string, hash uint64) {
		rs, _ := tx.Exec("INSERT INTO contacts(playerid, playerhash, type) VALUES($1, $2, $3)", player.DatabaseID(), hash, contactType)
		count, err := rs.RowsAffected()
		if err != nil {
			log.Warning.Println("Save(): INSERT failed for player friends:", err)
			if err := tx.Rollback(); err != nil {
				log.Warning.Println("Save(): Transaction insert friend rollback failed:", err)
			}
			return
		}

		if count <= 0 {
			log.Info.Println("Save(): Affected nothing for friend insertion!")
		}
	}
	clearItems := func() {
		if _, err := tx.Exec("DELETE FROM inventory WHERE playerid=$1", player.DatabaseID()); err != nil {
			log.Warning.Println("Save(): DELETE failed for player inventory:", err)
			if err := tx.Rollback(); err != nil {
				log.Warning.Println("Save(): Transaction delete inventory rollback failed:", err)
			}
			return
		}
	}
	clearBank := func() {
		if _, err := tx.Exec("DELETE FROM bank WHERE playerid=$1", player.DatabaseID()); err != nil {
			log.Warning.Println("Save(): DELETE failed for player bank:", err)
			if err := tx.Rollback(); err != nil {
				log.Warning.Println("Save(): Transaction delete bank rollback failed:", err)
			}
			return
		}
	}
	insertItem := func(id, amt int, worn bool) {
		rs, _ := tx.Exec("INSERT INTO inventory(playerid, itemid, amount, wielded) VALUES($1, $2, $3, $4)", player.DatabaseID(), id, amt, worn)
		count, err := rs.RowsAffected()
		if err != nil {
			log.Warning.Println("Save(): INSERT failed for player items:", err)
			if err := tx.Rollback(); err != nil {
				log.Warning.Println("Save(): Transaction insert item rollback failed:", err)
			}
			return
		}

		if count <= 0 {
			log.Info.Println("Save(): Affected nothing for item insertion!")
		}
	}
	insertBank := func(id, amt int) {
		rs, _ := tx.Exec("INSERT INTO bank(playerid, itemid, amount) VALUES($1, $2, $3)", player.DatabaseID(), id, amt)
		count, err := rs.RowsAffected()
		if err != nil {
			log.Warning.Println("Save(): INSERT failed for player bank items:", err)
			if err := tx.Rollback(); err != nil {
				log.Warning.Println("Save(): Transaction insert bank item rollback failed:", err)
			}
			return
		}

		if count <= 0 {
			log.Info.Println("Save(): Affected nothing for bank item insertion!")
		}
	}
	clearStats := func() {
		if _, err := tx.Exec("DELETE FROM stats WHERE playerid=$1", player.DatabaseID()); err != nil {
			log.Warning.Println("Save(): DELETE failed for player stats:", err)
			if err := tx.Rollback(); err != nil {
				log.Warning.Println("Save(): Transaction delete stats rollback failed:", err)
			}
			return
		}
	}
	insertStat := func(idx, cur, exp int) {
		rs, _ := tx.Exec("INSERT INTO stats(playerid, num, cur, exp) VALUES($1, $2, $3, $4)", player.DatabaseID(), idx, cur, exp)
		count, err := rs.RowsAffected()
		if err != nil {
			log.Warning.Println("Save(): INSERT failed for player stats:", err)
			if err := tx.Rollback(); err != nil {
				log.Warning.Println("Save(): Transaction insert stat rollback failed:", err)
			}
			return
		}

		if count <= 0 {
			log.Info.Println("Save(): Affected nothing for stat insertion!")
		}
	}
	clearAttributes()
	clearContactList("friend")
	clearContactList("ignore")
	clearItems()
	clearBank()
	clearStats()

	updateLocation()
	updateAppearance()
	player.Attributes.Range(insertAttribute)
	player.FriendList.ForEach(func(s string, b bool) bool {
		insertContact("friend", strutil.Base37.Encode(s))
		return false
	})
	for _, hash := range player.IgnoreList {
		insertContact("ignore", hash)
	}
	for stat := 0; stat < 18; stat++ {
		insertStat(stat, player.Skills().Current(stat), player.Skills().Experience(stat))
	}
	player.Inventory.Range(func(item *world.Item) bool {
		insertItem(item.ID, item.Amount, item.Worn)
		return true
	})
	player.Bank().Range(func(item *world.Item) bool {
		insertBank(item.ID, item.Amount)
		return true
	})

	if err := tx.Commit(); err != nil {
		log.Warning.Println("Save(): Error committing transaction for player update:", err)
	}
}
