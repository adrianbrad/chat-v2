package db

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func SetMigrationVersion(migrationsDir string, db *sql.DB, version uint) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	migrator, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir,
		"postgres", driver)
	if err != nil {
		return err
	}
	currVersion, _, _ := migrator.Version()
	err = migrator.Migrate(version)
	if currVersion != version {
		if err != nil {
			return err
		}
	}

	return nil
}
