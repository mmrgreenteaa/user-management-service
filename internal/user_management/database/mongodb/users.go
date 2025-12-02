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

// usersdb stores name collection in datebase mongodb
const usersdb = "test"

//  NoRecord indicates  about missing records
var NoRecord = errors.New("no record found")

// GetUserId provided by user id user
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

// AddUser add new user in collection.
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
			return fmt.Errorf("failed add a user duplication of data - %w", err)
		}
	}
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed add user - %w", err)
	}

	return nil
}

// DeleteUser deletes user. 
func (db *DB) DeleteUser(userId string) error {
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("failed edit login %w", err)
	}
	filter := bson.D{{Key: "_id", Value: objID}}
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

// EditLogin edit login user.
func (db *DB) EditLogin(userId string, newLogin string) error {
	ctx := context.Background()
	const docName = "test"
	collection := db.Database(usersdb).Collection(docName)

	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("failed edit login %w", err)
	}

	filter := bson.M{"_id": objID}
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
