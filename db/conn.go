package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type DatabaseConn struct {
	Conn *mongo.Client
	Name string
}

// Conn 连接数据库
func Conn(applyUrl string, maxPoolSize uint64) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	Client, err := mongo.Connect(ctx,
		options.Client().
			// 连接地址
			ApplyURI(applyUrl).
			// 设置连接数
			SetMaxPoolSize(maxPoolSize),
	)
	if err != nil {
		log.Println(err)
	}

	// 测试连接
	err = Client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println(err)
	}
	log.Println("Database Connect Success")

	return Client
}
