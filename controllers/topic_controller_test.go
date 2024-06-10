package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/zeekhoks/quiz-backend/services"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalln("No .env file available")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatalln("Mongodb URI string not found")
	}

	// connecting to the database
	err = services.ConnectToMongo(uri)
	if err != nil {
		log.Fatalln("Failed to connect to MongoDB")
	} else {
		log.Println("Connected to DB")
	}
}

func TestGetAllTopics(t *testing.T) {
	// Set up the router
	router := gin.Default()
	router.GET("/topics", GetAllTopics())

	// Test case: Get all topics
	req, _ := http.NewRequest("GET", "/topics", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
