package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zeekhoks/quiz-backend/models"
	"github.com/zeekhoks/quiz-backend/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"strings"
	"time"
)

func GetDisplayQuestionsByTopicHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		//get params from the request
		params := context.Request.URL.Query()
		if params.Get("topic") == "" {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Topic not provided in URL",
			})
			return
		}

		//get MongoDB client and questions collection
		DB := services.GetConnection()
		questionsCollection := services.GetCollection(DB, "questions")

		//retrieve topic from param
		topic := params.Get("topic")

		//MongoDB filter to search the question related to the `topic` param
		filter := bson.M{"$text": bson.M{"$search": topic}}
		cursor, err := questionsCollection.Find(context, filter)
		if err != nil {
			return
		}
		var questions []models.QuestionUnmarshal
		if err = cursor.All(context, &questions); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if questions == nil || len(questions) == 0 {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "No questions found with this topic",
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"questions": questions,
		})
	}
}

func UploadQuestionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		DB := services.GetConnection()
		questionsCollection := services.GetCollection(DB, "questions")
		topicsCollection := services.GetCollection(DB, "topics")

		file, _ := c.FormFile("questions_file")
		topic := c.PostForm("topic")
		if topic == "" || len(topic) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Topic not provided",
			})
			return
		}

		f, _ := file.Open()
		defer f.Close()
		content := make([]byte, file.Size)
		_, err := f.Read(content)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Server error. Please try again later",
			})
			return
		}

		var questions []models.QuestionUnmarshal
		err = json.Unmarshal(content, &questions)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON file. Please upload valid JSON",
			})
			return
		}
		var interfaces []interface{}
		for _, question := range questions {
			interfaces = append(interfaces, question)
		}

		topicDocument := bson.M{
			"topic": topic,
		}
		_, err = topicsCollection.InsertOne(c, topicDocument)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Server error. Unable to insert topics",
			})
			return
		}

		res, err := questionsCollection.InsertMany(c, interfaces)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Server error. Please try again later",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"number_of_questions_inserted": len(res.InsertedIDs),
			"topic":                        topic,
		})

	}
}

func GenerateQuizHandler() gin.HandlerFunc {
	return func(context *gin.Context) {

		params := context.Request.URL.Query()
		if params.Get("topic") == "" {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Topic not provided in URL",
			})
			return
		}

		DB := services.GetConnection()
		questionsCollection := services.GetCollection(DB, "questions")
		quizCollection := services.GetCollection(DB, "quizzes")

		topic := params.Get("topic")

		filter := bson.M{"$text": bson.M{"$search": topic}}
		cursor, err := questionsCollection.Find(context, filter)
		if err != nil {
			return
		}
		var questions []models.Question
		if err = cursor.All(context, &questions); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(questions) == 0 {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "No questions found with this topic",
			})
			return
		}

		startTime := time.Now()
		endTime := startTime.Add(30 * time.Minute)

		userAny, _ := context.Get("loggedInAccount")

		user := userAny.(models.User)

		quiz := &models.Quiz{
			Topic:         topic,
			User:          user,
			Questions:     questions,
			UserResponses: make([]models.UserResponse, 0),
			Completed:     false,
			StartTime:     startTime,
			EndTime:       endTime,
		}

		res, err := quizCollection.InsertOne(context, &quiz)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Server error. Please try again later",
			})
			return
		}

		quiz.Id = res.InsertedID.(primitive.ObjectID)

		context.JSON(http.StatusOK, gin.H{
			"quiz": quiz,
		})
	}
}

func SubmitAnswerHandler() gin.HandlerFunc {
	return func(context *gin.Context) {

		quizId := context.Param("id")

		DB := services.GetConnection()
		quizCollection := services.GetCollection(DB, "quizzes")

		quizIdParsed, err := primitive.ObjectIDFromHex(quizId)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Quiz ID is in the wrong format",
			})
			return
		}
		res := quizCollection.FindOne(context, bson.M{"_id": quizIdParsed})

		if res.Err() != nil {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Quiz with given ID not found",
			})
			return
		}

		var quiz models.Quiz
		err = res.Decode(&quiz)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server  error",
			})
			return
		}

		userAny, _ := context.Get("loggedInAccount")

		user := userAny.(models.User)

		if user.Username != quiz.User.Username {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Quiz not started by the same user",
			})
			return
		}

		if time.Now().After(quiz.EndTime) {
			quiz.Completed = true
		}

		if quiz.Completed == true {
			_, err = quizCollection.UpdateByID(context, quiz.Id, bson.D{
				{"$set", bson.D{
					{"completed", quiz.Completed},
				}},
			})

			if err != nil {
				context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error. Please try again",
				})
				return
			}
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Quiz has already ended. Start a new quiz",
			})
			return
		}

		body, err := io.ReadAll(context.Request.Body)

		bodyParsed, errs := validateUserResponseBody(body)

		if len(errs) != 0 {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": errs,
			})
			return
		}

		questionIdParsed, err := primitive.ObjectIDFromHex(bodyParsed["question_id"])

		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Question ID in wrong format",
			})
			return
		}
		var question models.Question

		for _, q := range quiz.Questions {
			if questionIdParsed.String() == q.ID.String() {
				found := false
				for _, option := range q.Options {
					if strings.ToLower(option) == bodyParsed["choice"] {
						found = true
					}
				}
				if !found {
					context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
						"error": "User choice is invalid for this current question",
					})
					return
				}
				question = q
			}
		}

		if question.ID.String() != questionIdParsed.String() {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Question with given ID not found in this particular quiz",
			})
			return
		}

		for _, response := range quiz.UserResponses {
			if response.QuestionId.String() == question.ID.String() {
				context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "Question with given ID has already been answered",
				})
				return
			}
		}

		userResponse := models.UserResponse{
			QuestionId:    question.ID,
			Response:      bodyParsed["choice"],
			CorrectAnswer: strings.ToLower(question.CorrectAnswer),
		}

		if userResponse.Response == strings.ToLower(question.CorrectAnswer) {
			userResponse.Result = "Right"
		} else {
			userResponse.Result = "Wrong"
		}

		quiz.UserResponses = append(quiz.UserResponses, userResponse)

		if len(quiz.UserResponses) == len(quiz.Questions) {
			quiz.Completed = true
		}

		updateDocument := bson.D{
			{"$set", bson.D{
				{"user_responses", quiz.UserResponses},
				{"completed", quiz.Completed},
			}},
		}

		_, err = quizCollection.UpdateByID(context, quiz.Id, updateDocument)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error please try again",
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"result": userResponse,
		})
	}
}

func QuizResultHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		quizId := context.Param("id")

		DB := services.GetConnection()
		quizCollection := services.GetCollection(DB, "quizzes")

		quizIdParsed, err := primitive.ObjectIDFromHex(quizId)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Quiz ID is in the wrong format",
			})
			return
		}
		res := quizCollection.FindOne(context, bson.M{"_id": quizIdParsed})

		if res.Err() != nil {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Quiz with given ID not found",
			})
			return
		}

		userAny, _ := context.Get("loggedInAccount")

		user := userAny.(models.User)

		if context.Request.ContentLength != 0 {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Body not expected for this request",
			})
			return
		}

		var quiz models.Quiz
		err = res.Decode(&quiz)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server  error",
			})
			return
		}

		if !user.IsAdmin && user.Username != quiz.User.Username {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "You don't have permissions to view this quiz's result",
			})
			return
		}

		if quiz.Completed != true {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Quiz has not ended yet. Answer all questions to get a result",
			})
			return
		}

		correctAnswers := 0

		for _, response := range quiz.UserResponses {
			if response.Result == "Right" {
				correctAnswers += 1
			}
		}

		context.JSON(http.StatusOK, gin.H{
			"quiz":           quiz,
			"user_responses": quiz.UserResponses,
			"stats": gin.H{
				"total_questions_answered":  len(quiz.UserResponses),
				"number_of_correct_answers": correctAnswers,
				"total_questions":           len(quiz.Questions),
				"percentage":                fmt.Sprintf("%.2f", float64(correctAnswers)/float64(len(quiz.Questions))*100),
			},
		})
	}
}

func validateUserResponseBody(body []byte) (map[string]string, []string) {

	var data map[string]interface{}
	err := json.Unmarshal(body, &data)
	errorStrings := make([]string, 0)

	parsed := make(map[string]string)

	if err != nil {
		errorStrings = append(errorStrings, "JSON is invalid")
		return nil, errorStrings
	}

	keysToValidate := []string{
		"question_id",
		"choice",
	}

	for _, key := range keysToValidate {
		if data[key] == nil {
			errorStrings = append(errorStrings, fmt.Sprintf("%v should be included in the body", key))
			continue
		} else {
			parsed[key] = strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", data[key])))
			if len(parsed[key]) == 0 {
				errorStrings = append(errorStrings, fmt.Sprintf("%v should not be empty", key))
				continue
			}
		}
	}

	if len(errorStrings) != 0 {
		return nil, errorStrings
	} else {
		return parsed, nil
	}

}
