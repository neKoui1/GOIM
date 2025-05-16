package models

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InitMongo() *mongo.Database {
	client, err := mongo.Connect(options.Client().ApplyURI(
		"mongodb://localhost:27017"))
	if err != nil {
		log.Println("Connect MongoDB error: ", err)
		return nil
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Println("Disconnect MongoDB error: ", err)
		}
	}()
	return client.Database("GOIM")
}
