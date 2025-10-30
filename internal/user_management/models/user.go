package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Login    string             `bson:"login,omitempty"`
	Password string             `bson:"pass,omitempty"`
}
