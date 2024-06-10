package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserResponse struct {
	QuestionId    primitive.ObjectID `json:"question_id" bson:"question_id"`
	Response      string             `json:"response" bson:"response"`
	Result        string             `json:"result" bson:"result"`
	CorrectAnswer string             `json:"correct_answer" bson:"correct_answer"`
}
