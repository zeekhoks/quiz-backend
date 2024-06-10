package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/zeekhoks/quiz-backend/models"
	"github.com/zeekhoks/quiz-backend/services"
	"log"
	"net/http"
	"os"
	"time"
)

func CreateNewUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User

		if err := ctx.Bind(&user); err != nil {
			log.Println("Failed to bind incoming payload with Gin", err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userExists, _ := services.UserExists(user.Username)
		if userExists {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}

		createdUser, err := services.CreateUser(user)
		if err != nil {
			log.Println("Failed to create user", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Unable to create an user"})
			return
		}

		ctx.JSON(http.StatusCreated, createdUser)
	}
}

func LoginHandler() gin.HandlerFunc {
	return func(context *gin.Context) {

		val, _ := context.Get("loggedInAccount")

		user, ok := val.(models.User)

		if !ok {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error. Try again",
			})
			return
		}

		claims := models.MyUserClaims{
			user,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 45).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		signingKey := os.Getenv("SIGNING_KEY")
		tokenString, err := token.SignedString([]byte(signingKey))
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error. Try again later",
			})
			return
		} else {
			context.JSON(http.StatusOK, gin.H{
				"user":  user,
				"token": tokenString,
			})
		}
	}
}
