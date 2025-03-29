package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Collection
var ctx context.Context

func init() {

	clientOptions := options.Client().ApplyURI("mongodb+srv://igorDitto:6ofqT6GgVlUFoUWm@cluster0.4dalr.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	db = client.Database("finance").Collection("transactions")
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
}
