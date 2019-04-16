package db

import (
	"database/sql"
	"fmt"

	"github.com/adrianbrad/chat-v2/configs"

	_ "github.com/lib/pq"
)

func ConnectDB(dbConfig configs.DBconfig) (db *sql.DB, err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port,
		dbConfig.User, dbConfig.Pass, dbConfig.Name)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return
	}
	err = db.Ping()
	if err != nil {
		return
	}
	return
}
