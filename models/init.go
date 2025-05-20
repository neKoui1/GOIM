package models

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Mongo = InitMongo()

func InitMongo() *mongo.Database {
	
	client, err := mongo.Connect(options.Client().
ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatal("无法连接到MongoDB:", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Fail to Ping mongo db : ", err)
	}

	return client.Database("GOIM")
}
