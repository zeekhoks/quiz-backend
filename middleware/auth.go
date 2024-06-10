package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/zeekhoks/quiz-backend/models"
	"github.com/zeekhoks/quiz-backend/services"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"os"
	"strings"
)

func BasicAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		username, password, ok := context.Request.BasicAuth()
		log.Println("Authenticating user", bson.M{"user username": username})
		if !ok {
			log.Println("Unable to authenticate user, failed to parse auth string")
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unable to authenticate user. Check username and password in Authorization header"})
			return
		} else {
			user, err := services.GetUserByUsername(username)
			if err != nil {
				log.Println("Unable to get user with username", err)
				context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unable to find user with this specific username"})
				return
			} else {
				success := services.CheckPasswordHash(password, user.Password)
				if success {
					context.Set("loggedInAccount", user)
					context.Next()
				} else {
					log.Println("Authentication failed")
					context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User password is wrong. Check again"})
					return
				}
			}
		}
	}
}

func AdminCheck() gin.HandlerFunc {
	return func(context *gin.Context) {

		val, _ := context.Get("loggedInAccount")

		user, ok := val.(models.User)

		if !ok {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error. Try again",
			})
			return
		}

		if !user.IsAdmin {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "You don't have access to make this request",
			})
			return
		}
		context.Next()
	}
}

func UserExtractor() gin.HandlerFunc {
	return func(context *gin.Context) {
		authorizationHeader := context.Request.Header.Get("Authorization")

		if len(authorizationHeader) == 0 || !strings.HasPrefix(authorizationHeader, "Bearer ") {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is not in correct format",
			})
			return
		}

		tokenString, _ := strings.CutPrefix(authorizationHeader, "Bearer ")

		signingKey := os.Getenv("SIGNING_KEY")

		token, err := jwt.ParseWithClaims(tokenString, &models.MyUserClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(signingKey), nil
		})

		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token validation failed. Resend valid token",
			})
			return
		}

		if claims, ok := token.Claims.(*models.MyUserClaims); ok && token.Valid {
			context.Set("loggedInAccount", claims.User)
		} else {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token validation failed. Resend valid token",
			})
			return
		}
	}
}
