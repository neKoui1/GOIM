package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Message struct {
	UserId    bson.ObjectID `bson:"user_id"`
	RoomId    bson.ObjectID `bson:"room_id"`
	Data      string        `bson:"data"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

func (Message) CollectionName() string {
	return "message"
}

func InsertOntMessage(msg *Message) error {
	msg.CreatedAt = time.Now()
	msg.UpdatedAt = time.Now()
	_, err := Mongo.Collection(Message{}.CollectionName()).
		InsertOne(context.Background(), msg)
	return err
}

func GetMessageListByRoomID(roomID bson.ObjectID, limit *int64, skip *int64) ([]*Message, error) {
	data := make([]*Message, 0)
	findOptions := options.Find()
	findOptions.SetLimit(*limit).SetSkip(*skip).SetSort(bson.D{
		{Key: "created_time", Value: -1},
	})
	cursor, err := Mongo.Collection(Message{}.CollectionName()).Find(
		context.Background(),
		bson.M{
			"room_id": roomID,
		},
		findOptions,
	)
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		msg := new(Message)
		err = cursor.Decode(msg)
		if err != nil {
			return nil, err
		}
		data = append(data, msg)
	}
	return data, nil
}
