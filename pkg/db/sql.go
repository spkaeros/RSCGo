package db

import (
	"context"
	"database/sql"
	"sync"

	"github.com/spkaeros/rscgo/pkg/config"
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
func (s *sqlService) sqlOpen(addr string) *sql.DB {
	database, err := sql.Open(s.Driver, addr) //"file:"+config.DataDir()+addr)
	if err != nil {
		log.Error.Println("Couldn't load database (driver: "+s.Driver+", addr: "+addr+"):", err)
		return nil
	}
	err = database.Ping()
	if err != nil {
		log.Error.Println("Couldn't load database (driver: "+s.Driver+", addr: "+addr+"):", err)
		return nil
	}
	s.database = database
	return database
}

//sqlService A database/sql based persistence service.
// Implements PlayerService interface and sqlService.
type sqlService struct {
	database *sql.DB
	conn *sql.Conn
	Driver   string
	context  context.Context
	sync.RWMutex
}

//newSqlService returns a new sqlService instance attached to the provided *sql.DB
// To obtain a valid *sql.DB, load a database/sql driver and call sqlOpen(driverName, connectAddress string) *sql.DB
func newSqlService(driver string) *sqlService {
	return &sqlService{
		Driver: driver,
	}
}

var dbConn *sqlService

//connect returns a connection to the services underlying *sql.DB instance upon successful
// connection.  If an error occurs, returns nil.
func (s *sqlService) connect(ctx context.Context) *sql.Conn {
	if dbConn == nil {
		dbConn = newSqlService(config.PlayerDriver())
		db := dbConn.sqlOpen(config.PlayerDB())
		dbConn.database = db
		c, err := db.Conn(ctx)
		if err != nil {
			return nil
		}
		dbConn.conn = c
	}
	return dbConn.conn
}
