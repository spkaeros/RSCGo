package db

import (
	"database/sql"

	"github.com/spkaeros/rscgo/pkg/game/config"
	"github.com/spkaeros/rscgo/pkg/log"

	// Necessary for sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

//Open Returns an active sqlite3 database reference for the specified database file.
func Open(file string) *sql.DB {
	database, err := sql.Open("sqlite3", "file:"+config.DataDir()+file)
	if err != nil {
		log.Error.Println("Couldn't load SQLite3 database:", err)
		return nil
	}
	database.SetMaxOpenConns(1)
	return database
}
