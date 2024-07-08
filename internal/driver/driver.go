package driver

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func OpenDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping() // verify db connection
	if err != nil {
		fmt.Println("Error when attemping to connect to the database:", err)
	}

	return db, err
}
