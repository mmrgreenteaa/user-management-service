// redis provides jwt storage functions in redis for fast response
package redis

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

var ErrNoRecord = errors.New("redis:key dont found")

// RCongif holds the configuration settings for the redis connection.
type RCongif struct {
	Addr     string  `mapstructure:"listen_addr"`
	Password string  `mapstructure:"password"`
}

// DB represents a database handle containing the connection poo redis
type DB struct {
	*redis.Client
}

// MustConnect connects to redis
func Connect(cfg *RCongif) *DB {
	//"localhost:6379"
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: "",
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Println("Redis connection was refused - %w", err)
	}

	return &DB{
		Client: client,
	}
}

// AddJwt adds new jwt and user id
func (db *DB) AddJwt(jwt string, userId string) error {

	hash := sha256.Sum256([]byte(jwt))
	key := hex.EncodeToString(hash[:])
	err := db.Set(key, userId, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("redis: falied to adds jwt - %w", err)
	}
	return nil
}

// GetJwt  gets the user id belonging to jwt
func (db *DB) GetJwt(jwt string) (string, error) {
	hash := sha256.Sum256([]byte(jwt))
	key := hex.EncodeToString(hash[:])

	val, err := db.Get(key).Result()
	if err == redis.Nil {
		return "", ErrNoRecord
	} else if err != nil {
		return "", fmt.Errorf("redis:failed to get jwt %w", err)
	}

	return val, nil

}
