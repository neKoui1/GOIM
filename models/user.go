package models

import (
	"GOIM/helper"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	UserStatusOnline    = 1
	UserStatusOffline   = 0
	UserStatusInvisible = 2
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	Account   string        `bson:"account" json:"account"`
	Password  string        `bson:"password" json:"-"`
	Nickname  string        `bson:"nickname" json:"nickname"`
	Gender    bool          `bson:"gender" json:"gender"` // false 女 true 男
	Email     string        `bson:"email" json:"email"`
	Avatar    string        `bson:"avatar" json:"avatar"`
	Status    int           `bson:"status" json:"status"` // 0 离线 1 在线
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
	LastLogin time.Time     `bson:"last_login" json:"last_login"`
}

func (User) CollectionName() string {
	return "user"
}

func (u *User) SetLastLoginNow() error {
	filter := bson.D{
		{Key: "account", Value: u.Account},
	}
	update := bson.D{
		{
			Key: "$set", Value: bson.D{
				{Key: "last_login", Value: time.Now()},
				{Key: "status", Value: 1},
			},
		},
	}
	_, err := GetMongo().Collection(User{}.CollectionName()).
		UpdateOne(context.Background(), filter, update)

	return err
}

func GetUserList() ([]*User, error) {
	userList := make([]*User, 0)
	cursor, err := GetMongo().Collection(User{}.CollectionName()).
		Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &userList)
	if err != nil {
		return nil, err
	}
	return userList, nil
}

func GetUserByAccountPassword(account, password string) (*User, error) {
	u := new(User)
	err := GetMongo().Collection(User{}.CollectionName()).
		FindOne(context.Background(), bson.D{
			{Key: "account", Value: account},
		}).Decode(u)
	if err != nil {
		return nil, err
	}
	if !helper.CheckPassword(password, u.Password) {
		return nil, fmt.Errorf("密码错误")
	}
	return u, nil
}

func GetUserByAccount(account string) (*User, error) {
	u := new(User)
	err := GetMongo().Collection(User{}.CollectionName()).FindOne(
		context.Background(), bson.D{
			{Key: "account", Value: account},
		},
	).Decode(u)
	return u, err
}

func GetUserByID(ID bson.ObjectID) (*User, error) {
	u := new(User)

	err := GetMongo().Collection(User{}.CollectionName()).
		FindOne(context.Background(), bson.D{
			{Key: "_id", Value: ID},
		}).Decode(u)
	return u, err
}

func GetUserCountByEmail(email string) (int64, error) {
	return GetMongo().Collection(User{}.CollectionName()).
		CountDocuments(context.Background(), bson.D{
			{Key: "email", Value: email},
		})
}

func GetUserCountByAccount(account string) (int64, error) {
	return GetMongo().Collection(User{}.CollectionName()).
		CountDocuments(context.Background(), bson.D{
			{Key: "account", Value: account},
		})
}

func InsertOntUser(u *User) error {
	_, err := GetMongo().Collection(User{}.CollectionName()).
		InsertOne(context.Background(), u)
	return err
}
