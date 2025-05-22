package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRoom struct {
	UserId bson.ObjectID `bson:"user_id"`
	RoomId bson.ObjectID `bson:"room_id"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}