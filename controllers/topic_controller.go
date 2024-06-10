package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zeekhoks/quiz-backend/services"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func GetAllTopics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		DB := services.GetConnection()
		questionsCollection := services.GetCollection(DB, "topics")

		cursor, err := questionsCollection.Find(ctx, bson.M{})
		if err != nil {
			return
		}
		var topics []bson.M
		if err = cursor.All(ctx, &topics); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, topics)
	}
}
