package db

import (
	"database/sql"

	"bitbucket.org/zlacki/rscgo/pkg/server/config"
	"bitbucket.org/zlacki/rscgo/pkg/server/log"

	// Necessary for sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

//OpenDatabase Returns an active sqlite3 database reference for the specified database file.
func OpenDatabase(file string) *sql.DB {
	database, err := sql.Open("sqlite3", "file:"+config.DataDir()+file)
	if err != nil {
		log.Error.Println("Couldn't load SQLite3 database:", err)
		return nil
	}
	database.SetMaxOpenConns(1)
	return database
}
