package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Add 在数据库中添加内容
func Add(conn *DatabaseConn, Collection string, addData bson.D) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	actionCollection := conn.Conn.Database(conn.Name).Collection(Collection)
	_, err := actionCollection.InsertOne(ctx, addData)
	if err != nil {
		log.Println(err)
	}
}
