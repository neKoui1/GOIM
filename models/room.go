package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Room struct {
	Number    string        `bson:"number"`
	Name      string        `bson:"name"`
	Info      string        `bson:"info"`
	UserId    bson.ObjectID `bson:"user_id"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

func (Room) CollectionName() string {
	return "room"
}

func InsertOntRoom(r *Room) error {
	_, err := Mongo.Collection(Room{}.CollectionName()).
		InsertOne(context.Background(), r)
	return err
}
