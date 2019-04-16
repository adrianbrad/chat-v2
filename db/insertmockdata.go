package db

import (
	"database/sql"
	"io/ioutil"
)

func ExecuteSQLfile(db *sql.DB, sqlFilePath string) (err error) {
	m, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		return
	}
	sql := string(m)
	_, err = db.Exec(sql)
	return
}
