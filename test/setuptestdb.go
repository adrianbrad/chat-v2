package test

import (
	"database/sql"

	"github.com/adrianbrad/chat-v2/configs"
	d "github.com/adrianbrad/chat-v2/db"
)

func SetupTestDB() (db *sql.DB, err error) {
	config, err := configs.LoadDBconfig("../../../configs/test-db-config.yaml")
	if err != nil {
		return
	}
	db, err = d.ConnectDB(config)
	if err != nil {
		return
	}

	err = d.ExecuteSQLfile(db, "../../../db/schema.sql")
	if err != nil {
		return
	}

	err = d.ExecuteSQLfile(db, "../../../db/insert-mock-data.sql")

	return
}
