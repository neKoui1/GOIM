package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRoom struct {
	UserId    bson.ObjectID `bson:"user_id"`
	RoomId    bson.ObjectID `bson:"room_id"`
	RoomType  int           `bson:"room_type"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

func (UserRoom) CollectionName() string {
	return "user_room"
}

func GetUserRoomByUserIDRoomID(userID bson.ObjectID, roomID bson.ObjectID) (*UserRoom, error) {
	ur := new(UserRoom)
	err := GetMongo().Collection(UserRoom{}.CollectionName()).
		FindOne(context.Background(), bson.D{
			{Key: "user_id", Value: userID},
			{Key: "room_id", Value: roomID},
		}).Decode(ur)
	if err != nil {
		return nil, err
	}
	return ur, nil
}

func GetUserRoomByRoomID(roomID bson.ObjectID) ([]*UserRoom, error) {
	cursor, err := GetMongo().Collection(UserRoom{}.CollectionName()).Find(context.Background(), bson.D{
		{Key: "room_id", Value: roomID},
	})
	if err != nil {
		return nil, err
	}
	urs := make([]*UserRoom, 0)
	for cursor.Next(context.Background()) {
		ur := new(UserRoom)
		err = cursor.Decode(ur)
		if err != nil {
			return nil, err
		}
		urs = append(urs, ur)
	}
	return urs, nil
}

func JudgeUserIsFriend(userId1, userId2 bson.ObjectID) (bool, error) {
	cursor, err := GetMongo().Collection(UserRoom{}.CollectionName()).
		Find(context.Background(), bson.D{
			{Key: "user_id", Value: userId1},
			{Key: "room_type", Value: 1},
		})
	roomIds := make([]bson.ObjectID, 0)
	if err != nil {
		return false, err
	}
	for cursor.Next(context.Background()) {
		ur := new(UserRoom)
		err = cursor.Decode(ur)
		if err != nil {
			return false, err
		}
		roomIds = append(roomIds, ur.RoomId)
	}
	// 获取userId2 私聊房间个数
	cnt, err := GetMongo().Collection(UserRoom{}.CollectionName()).
		CountDocuments(context.Background(), bson.M{
			"user_id": userId2,
			"room_id": bson.M{
				"$in": roomIds,
			},
		})
	if err != nil {
		return false, err
	}
	if cnt > 0 {
		return true, nil
	}
	return false, nil
}

func InsertOneUserRoom(ur *UserRoom) error {
	_, err := GetMongo().Collection(UserRoom{}.CollectionName()).
		InsertOne(context.Background(), ur)
	return err
}
