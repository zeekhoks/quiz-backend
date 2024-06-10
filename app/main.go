package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/zeekhoks/quiz-backend/routes"
	"github.com/zeekhoks/quiz-backend/services"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
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

func main() {
	router := routes.GetRouter()

	err := router.Run(":" + os.Getenv("SERVER_PORT"))

	if err != nil {
		fmt.Printf("Fatal error has occured: %v\n", err)
	}
}
