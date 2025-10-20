package postgresql

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func Сonnect() *DB {

	dsn := "host=localhost user=postgres password=QWERTY dbname=auth_service port=5432 search_path= tokens_info"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("error connect")
	}

	return &DB{db}

}
