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

var router *gin.Engine

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

	// Set up the router
	router = gin.Default()
	router.GET("/questions", GetDisplayQuestionsByTopicHandler())
}

func TestGetDisplayQuestionsByEmptyTopic(t *testing.T) {
	router := gin.Default()
	router.GET("/questions", GetDisplayQuestionsByTopicHandler())

	req, err := http.NewRequest("GET", "/questions", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetDisplayQuestionsByValidTopic(t *testing.T) {
	router := gin.Default()
	router.GET("/questions", GetDisplayQuestionsByTopicHandler())

	req, err := http.NewRequest("GET", "/questions?topic=france", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGenerateQuizHandler(t *testing.T) {
	// Set up the router
	router := gin.Default()
	router.GET("/quiz", GenerateQuizHandler())

	// Test case: Topic not provided in URL
	req, _ := http.NewRequest("GET", "/quiz", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Test case: No questions found with this topic
	req, _ = http.NewRequest("GET", "/quiz?topic=nonexistent", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestSubmitAnswerHandler(t *testing.T) {
	// Set up the router
	router := gin.Default()
	router.POST("/quiz/:id/response", SubmitAnswerHandler())

	// Test case: Quiz ID not provided in URL
	req, _ := http.NewRequest("POST", "/quiz//response", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Test case: Quiz with given ID not found
	req, _ = http.NewRequest("POST", "/quiz/invalidID/response", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
