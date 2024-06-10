package services

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var DB *mongo.Client

func ConnectToMongo(URI string) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		log.Println(err)
	}
	DB = client

	err = DB.Ping(ctx, nil)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}

func GetConnection() *mongo.Client {
	return DB
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("Pearson").Collection(collectionName)
	return collection
}
