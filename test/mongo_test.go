package test

import (
	"context"
	"fmt"
	"testing"

	"GOIM/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestFindOne(t *testing.T) {

	client, err := mongo.Connect(options.Client().
		ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			t.Fatal(err)
		}
	}()

	coll := client.Database("GOIM").Collection("user")

	var user models.User
	err = coll.FindOne(context.TODO(),
		bson.D{}).Decode(&user)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(user)
}

func TestFindAllUsers(t *testing.T) {
	client, err := mongo.Connect(options.Client().
		ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			t.Fatal(err)
		}
	}()

	coll := client.Database("GOIM").Collection("user")
	userList := make([]*models.User, 0)
	cursor, err := coll.Find(context.Background(), bson.D{})
	if err != nil {
		t.Fatal(err)
	}
	err = cursor.All(context.Background(), &userList)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*userList[0])
	for _, v := range userList {
		t.Log(v.Account)
		t.Log(v.Nickname)
		t.Log(v.Gender)
	}
}

func TestFind(t *testing.T) {
	client, err := mongo.Connect(
		options.Client().ApplyURI("mongodb://localhost:27017"),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			t.Fatal(err)
		}
	}()

	coll := client.Database("GOIM").Collection("user_room")

	cursor, err := coll.Find(context.Background(), bson.D{})

	userRooms := make([]*models.UserRoom, 0)
	for cursor.Next(context.Background()) {
		ur := new(models.UserRoom)
		err := cursor.Decode(ur)
		if err != nil {
			t.Fatal(err)
		}
		userRooms = append(userRooms, ur)
	}
	for _, v := range userRooms {
		fmt.Println("User Room => ", v)
	}
}
