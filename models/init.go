package models

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	Mongo        *mongo.Database
	mongoClient  *mongo.Client
	initOnce     sync.Once
	shutdownOnce sync.Once
)

func GetMongo() *mongo.Database {
	initOnce.Do(func() {
		ConnectMongoDB()
		createIndexes()
		// registerShutdownHook()
	})
	return Mongo
}

// 初始化MongoDB连接
func ConnectMongoDB() {
	uri := "mongodb://localhost:27017"
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(10 * time.Second)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatalf("Fail to connect to mongodb: %v\n", err)
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Fail to ping mongodb: %v\n", err)
	}

	mongoClient = client
	Mongo = client.Database("GOIM")
	log.Println("Connect to mongodb successfully")
}

func registerShutdownHook() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c //阻塞等待信号
		shutdownOnce.Do(func() {
			log.Println("接收到关闭信号，正在断开mongodb连接")
			ctx, cancel := context.WithTimeout(
				context.Background(),
				5*time.Second,
			)
			defer cancel()
			if err := mongoClient.Disconnect(ctx); err != nil {
				log.Printf("关闭mongodb连接失败: %v", err)
			} else {
				log.Println("mongodb连接已关闭")
			}
		})
	}()
}

// 由main函数调用，关闭mongodb连接
func CloseMongo() {
	shutdownOnce.Do(func() {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Printf("关闭mongodb连接失败: %v", err)
		} else {
			log.Println("mongodb连接已关闭")
		}
	})
}

func createIndexes() {
	log.Println("开始创建数据库索引...")
	createMessageIndexes()
	log.Println("数据库索引创建完成")
}

func createMessageIndexes() {
	collection := Mongo.Collection(Message{}.CollectionName())
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "room_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().SetName("room_created_desc"),
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().SetName("user_created_desc"),
		},
		{
			Keys: bson.D{
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().SetName("created_at_desc"),
		},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		log.Printf("消息索引创建失败: %v\n", err)
	} else {
		log.Println("消息索引创建成功")
	}
}
