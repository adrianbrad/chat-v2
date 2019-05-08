package db

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(migrationsDir string, db *sql.DB) (err error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir,
		"postgres", driver)
	m.Steps(1)
	return
}
