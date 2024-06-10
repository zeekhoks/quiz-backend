package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `json:"first_name" bson:"first_name"`
	LastName  string             `json:"last_name" bson:"last_name"`
	Username  string             `json:"username" bson:"username"`
	Password  string             `json:"-" bson:"password"`
	IsAdmin   bool               `json:"is_admin" bson:"is_admin"`
}
