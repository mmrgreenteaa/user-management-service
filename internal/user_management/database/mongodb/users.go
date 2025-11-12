package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/mmrgreenteaa/user-management-service/internal/user_management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersdb = "test"

var NoRecord = errors.New("no record")

func (db *DB) GetUserId(login string, pass string) (string, error) {

	ctx := context.Background()
	filter := bson.D{{Key: "login", Value: login}, {Key: "pass", Value: pass}}
	const docName = "test"

	collection := db.Database(usersdb).Collection(docName)

	user := struct {
		Login    string             `bson:"login,omitempty"`
		Password string             `bson:"pass,omitempty"`
		ID       primitive.ObjectID `bson:"_id,omitempty"`
	}{}
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("fail getting a user don't find - %w", err)
		}
		return "", fmt.Errorf("fail getting a user - %w", err)
	}

	return user.ID.Hex(), nil
}

func (db *DB) AddUser(login string, pass string) error {
	//NOTE: шифрование паролей
	filter := bson.D{{Key: "login", Value: login}}
	ctx := context.Background()

	const docName = "test"

	collection := db.Database(usersdb).Collection(docName)
	user := models.User{
		Login:    login,
		Password: pass,
	}

	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return fmt.Errorf("fail getting a user - %w", err)
		}
	}
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed add user - %w", err)
	}

	return nil
}

func (db *DB) DeleteUser(userId string) error {
	filter := bson.D{{Key: "user_id", Value: userId}}
	ctx := context.Background()
	const docName = "test"
	collection := db.Database(usersdb).Collection(docName)
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed delete user %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("failed  delete user document: %w", NoRecord)
	}
	return nil
}

func (db *DB) EditLogin(userId string, newLogin string) error {
	ctx := context.Background()
	const docName = "test"
	collection := db.Database(usersdb).Collection(docName)

	filter := bson.M{"user_id": userId}
	update := bson.M{
		"$set": bson.M{
			"login": newLogin,
		},
	}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("fail update user %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("fail: Document not found: %w", NoRecord)
	}

	return nil
}
