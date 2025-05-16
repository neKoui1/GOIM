package models

import (
	"time"
)

type User struct {
	ID        string    `bson:"_id"`
	Account   string    `bson:"account"`
	Password  string    `bson:"password"`
	Nickname  string    `bson:"nickname"`
	Gender    bool      `bson:"gender"` // false 女 true 男
	Email     string    `bson:"email"`
	Avatar    string    `bson:"avatar"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (User) CollectionName() string {
	return "user"
}
