package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
//dbName = flag.String("database-name", "asa_tickets_app", "")
)

func migrateDB(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "connecting to the database")
	}

	migraitionSource := fmt.Sprintf("file://%sinternal/storage/database/migrations", dataDir)
	migrator, err := migrate.NewWithDatabaseInstance(migraitionSource, dbName, driver)
	if err != nil {
		return errors.Wrap(err, "creating migrator")
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "executing migration")
	}

	version, dirty, err := migrator.Version()
	if err != nil {
		return errors.Wrap(err, "getting migration verison")
	}

	// if the migration is dirty then attempl to reverse it
	if dirty {
		err = migrator.Down()
		if err != nil {
			logrus.Debug("Error reversing failed migration")
		}
	}

	logrus.WithFields(logrus.Fields{
		"version": version,
		"dirty":   dirty,
	}).Debug()
	return nil
}
