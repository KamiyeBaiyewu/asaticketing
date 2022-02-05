package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	/* 	dbDriver  = flag.String("db-driver", "postgres", "name of driver used (Postgres)")
	   	dbUser    = flag.String("db-user", "postgres", "User that would used to access database server (Postgres)")
	   	dbSecret  = flag.String("db-secret", "Peaceg419", "Password that would used to access database server (Postgres)")
	   	dbHost    = flag.String("db-host", "db", "Database to connect to")
	   	dbPort    = flag.Int("db-port", 5432, "Connection port for the database server (Postgres)")
	   	dbName    = flag.String("db-name", "asa_tickets_app", "Database to connect to")
	   	dbParam   = flag.String("db-param", "sslmode=disable", "Option on wheter to enable SSL")
	   	dbTimeout = flag.Int64("database-timeout-ms", 2000, "")

	   	// dbURL     = flag.String("database-url", "postgres://postgres:Peaceg419@db:5432/asa_tickets_app?sslmode=disable", "")
	   	dbURL = flag.String("database-url", "postgres://postgres:Peaceg419@db:5432/asa_tickets_app?sslmode=disable", "") */
	dbDriver  string
	dbUser    string
	dbSecret  string
	dbHost    string
	dbPort    int
	dbName    string
	dbParam   string
	dbTimeout int64
	dbURL     string
	dataDir   string
)

// Connect makes a new database Connection.
func Connect() (*sqlx.DB, error) {

	var databaseURL string

	if dbURL != "" {

		databaseURL = dbURL
	} else {
		databaseURL = fmt.Sprintf("%s://%s:%s@%s:%d/%s?%s", dbDriver, dbUser, dbSecret, dbHost, dbPort, dbName, dbParam)
	}

	conn, err := sqlx.Open(dbDriver, databaseURL)

	if err != nil {
		return nil, errors.Wrap(err, "Could not Connect to the Database")
	}
	conn.SetMaxOpenConns(32)

	// check if the database is running
	if err := waitForDB(conn.DB); err != nil {
		return nil, err
	}

	// Migrate Database Schema
	if err := migrateDB(conn.DB); err != nil {
		return nil, errors.Wrap(err, "could not migrate database")
	}

	return conn, nil
}

// New creates a new databse
func New(cfg *config.Info) (Database, error) {

	dbDriver = cfg.Database.Driver
	dbUser = cfg.Database.Username
	dbSecret = cfg.Database.Password
	dbName = cfg.Database.Database
	dbHost = cfg.Database.Hostname
	dbPort = cfg.Database.Port
	dbParam = cfg.Database.Parameter
	dbTimeout = cfg.Database.Timeout
	dbURL = cfg.Database.URL
	dataDir = cfg.DataDirectory

	conn, err := Connect()
	if err != nil {
		return nil, err
	}

	database := &database{conn: conn}
	return database, nil
}

func waitForDB(conn *sql.DB) error {
	ready := make(chan struct{})
	go func() {
		for {
			if err := conn.Ping(); err == nil {
				logrus.Debug("Database Connected")
				close(ready)
				return
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()

	select {
	case <-ready:
		return nil
	case <-time.After(time.Duration(dbTimeout) * time.Millisecond):
		return errors.New("database not ready")
	}
}
