package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB represents a database handle containing the connection poo mongodb
type DB struct {
	*mongo.Client
}

// DbCongig holds the configuration settings for the mongodb connection.
type DbConfig struct {
	Ip   string `mapstructure:"ip"`
	Pass string
}

// Connect establishes a connection to the mongodb database using the provided configuration.
// It returns an initialized DB instance.
func Сonnect(confg *DbConfig) *DB {
	ctx := context.TODO()
	conn := fmt.Sprintf("mongodb://%s", confg.Ip)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{client}
}
