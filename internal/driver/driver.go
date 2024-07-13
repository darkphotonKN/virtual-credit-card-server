package driver

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenDB(dsn string) (*gorm.DB, error) {

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	log.Println("db:", db)

	if err != nil {
		log.Fatalf("Error when attemping to connect to the database:", err)
		return nil, err
	}

	return db, err
}
