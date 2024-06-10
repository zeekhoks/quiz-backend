package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Quiz struct {
	Id            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	User          User               `json:"user" bson:"user"`
	Topic         string             `json:"topic" bson:"topic"`
	Questions     []Question         `json:"questions" bson:"questions"`
	UserResponses []UserResponse     `json:"-" bson:"user_responses"`
	Completed     bool               `json:"-" bson:"completed"`
	StartTime     time.Time          `json:"start_time" bson:"start_time"`
	EndTime       time.Time          `json:"end_time" bson:"end_time"`
}
