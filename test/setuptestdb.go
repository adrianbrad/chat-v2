package test

import (
	"database/sql"
	"os"
	"strings"

	"github.com/adrianbrad/chat-v2/configs"
	d "github.com/adrianbrad/chat-v2/db"
)

func SetupTestDB() (db *sql.DB, err error) {
	currDir, _ := os.Getwd()
	chatRoot := currDir[:strings.LastIndex(currDir, "/chat-v2/")+len("/chat-v2/")]
	config, err := configs.LoadDBconfig(chatRoot + "/configs/test-db-config.yaml")
	if err != nil {
		return
	}
	db, err = d.ConnectDB(config)
	if err != nil {
		return
	}

	err = d.ExecuteSQLfile(db, chatRoot+"/db/schema.sql")
	if err != nil {
		return
	}

	err = d.ExecuteSQLfile(db, chatRoot+"/db/insert-mock-data.sql")

	return
}
