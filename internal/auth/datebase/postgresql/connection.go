package postgresql

import (
	"log"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
	logger slog.Logger
}

func Сonnect() *DB {

	dsn := "host=localhost user=postgres password=QWERTY dbname=auth_service port=5432 search_path= tokens_info"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("error connect")
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return &DB{db, *logger}

}
