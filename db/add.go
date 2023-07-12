package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

// Add 在数据库中添加内容
func Add(c *DatabaseConn, Collection string, addData bson.D) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	actionCollection := c.Conn.
		Database(c.Name).
		Collection(Collection)

	_, err := actionCollection.InsertOne(ctx, addData)
	if err != nil {
		log.Println(err)
	}
}
