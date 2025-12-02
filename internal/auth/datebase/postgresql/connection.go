// Package postgresql provides tools for managing PostgreSQL database operations.
package postgresql

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB represents a database handle containing the connection pool and a logger.
type DB struct {
	*gorm.DB
	logger slog.Logger
}

// DbConfig holds the configuration settings for the PostgreSQL connection.
type DbConfig struct {
	Host       string `mapstructure:"host"`
	User       string `mapstructure:"user"`
	Pass       string `mapstructure:"password"`
	Name       string `mapstructure:"name"`
	Port       string `mapstructure:"port"`
	SearchPath string `mapstructure:"searchPath"`
}

// Connect establishes a connection to the PostgreSQL database using the provided configuration.
// It returns an initialized DB instance.
func Connect(confg *DbConfig) *DB {

	dsn := "host=%s user=%s password=%s dbname=%s port=%s search_path= %s"
	dsn = fmt.Sprintf(dsn, confg.Host, confg.User, confg.Pass, confg.Name, confg.Port, confg.SearchPath)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("error connect to posgresql")
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return &DB{db, *logger}

}
