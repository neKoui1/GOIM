package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Room struct {
	Id        bson.ObjectID `bson:"_id"`
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

func GetRoomByNumber(number string) (*Room, error) {
	room := new(Room)
	err := GetMongo().Collection(Room{}.CollectionName()).
		FindOne(context.Background(), bson.D{
			{Key: "number", Value: number},
		},
		).Decode(room)
	return room, err
}

func InsertOntRoom(r *Room) error {
	_, err := GetMongo().Collection(Room{}.CollectionName()).
		InsertOne(context.Background(), r)
	return err
}
