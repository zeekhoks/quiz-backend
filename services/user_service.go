package services

import (
	"github.com/zeekhoks/quiz-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"log"
)

func CreateUser(user models.User) (*mongo.InsertOneResult, error) {
	client := GetConnection()
	collection := GetCollection(client, "users")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)
	result, err := collection.InsertOne(context.TODO(), user)
	return result, err
}

func GetUserByUsername(username string) (models.User, error) {
	client := GetConnection()
	collection := GetCollection(client, "users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	return user, err
}

func UserExists(username string) (bool, error) {
	client := GetConnection()
	collection := GetCollection(client, "users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Println("unable to validate password", err)
	}
	return err == nil
}
