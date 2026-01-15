package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

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

	conn := fmt.Sprintf("mongodb://%s", confg.Ip)
	opts := options.Client().ApplyURI(conn)

	for i := 0; i < 4; i++ {
		log.Printf("Попытка подключения к MongoDB %d/4...", i+1)

		client, err := mongo.Connect(context.Background(), opts)
		if err != nil {
			log.Printf("Ошибка в конфигурации URI: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = client.Ping(ctx, nil)
		cancel()

		if err == nil {
			log.Println("Успешное подключение к MongoDB")
			return &DB{Client: client}
		}

		log.Printf("База не ответила на Ping: %v", err)

		_ = client.Disconnect(context.Background())

		time.Sleep(10 * time.Second)
	}

	log.Fatal("Не удалось подключиться к MongoDB после 4 попыток")
	return nil
}
