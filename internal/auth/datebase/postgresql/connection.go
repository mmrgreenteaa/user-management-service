package postgresql

import (
	"fmt"
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

type DbConfig struct {
	Host       string `mapstructure:"host"`
	User       string `mapstructure:"user"`
	Pass       string `mapstructure:"password"`
	Name       string `mapstructure:"name"`
	Port       string `mapstructure:"port"`
	SearchPath string `mapstructure:"searchPath"`
}

func Сonnect(confg *DbConfig) *DB {
	
	dsn := "host=%s user=%s password=%s dbname=%s port=%s search_path= %s"
 dsn = 	fmt.Sprintf(dsn, confg.Host, confg.User, confg.Pass, confg.Name, confg.Port, confg.SearchPath)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("error connect to posgresql")
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return &DB{db, *logger}

}
