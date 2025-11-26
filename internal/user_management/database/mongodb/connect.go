package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	*mongo.Client
}

type DbConfig struct {
	Ip   string `mapstructure:"ip"`
	Pass string
}

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
