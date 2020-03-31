package db

import (
	"context"
	"database/sql"

	"github.com/spkaeros/rscgo/pkg/log"

	// Necessary for sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
	// Necessary for postgresql driver
	_ "github.com/lib/pq"
)

//sqlOpen Attempts to connect to the specified address as a database/sql database.
// Returns: upon success, a connected *sql.DB instance accessing the specified SQL
// database, and nil
// upon failure, nil and a meaningful error.
func (db *sqlService) sqlOpen(addr string) *sql.DB {
	database, err := sql.Open(db.Driver, addr) //"file:"+config.DataDir()+addr)
	if err != nil {
		log.Error.Println("Couldn't load database (driver: "+db.Driver+", addr: "+addr+"):", err)
		return nil
	}
	err = database.Ping()
	if err != nil {
		log.Error.Println("Couldn't load database (driver: "+db.Driver+", addr: "+addr+"):", err)
		return nil
	}
	db.database = database
	return database
}

//sqlService A database/sql based persistence service.
// Implements PlayerService interface and sqlService.
type sqlService struct {
	database *sql.DB
	Driver string
}

//newSqlService returns a new sqlService instance attached to the provided *sql.DB
// To obtain a valid *sql.DB, load a database/sql driver and call sqlOpen(driverName, connectAddress string) *sql.DB
func newSqlService(driver string) *sqlService {
	return &sqlService{
		Driver: driver,
	}
}

//connect returns a connection to the services underlying *sql.DB instance upon successful
// connection.  If an error occurs, returns nil.
func (s *sqlService) connect(ctx context.Context) *sql.Conn {
	conn, err := s.database.Conn(ctx)
	if err != nil {
		log.Info.Println("Error connecting to SQLite3 service:", err)
		return nil
	}
	return conn
}
