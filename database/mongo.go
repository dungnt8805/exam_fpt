package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func MongoInit() {
	var err error
	mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal(err)
	}
}

func GetMgoDB() *mongo.Database {
	if mongoClient == nil {
		MongoInit()
	}
	return mongoClient.Database("exam_fpt")
}
