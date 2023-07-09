package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

// Add 在数据库中添加内容
func Add(ConnDB *mongo.Client, Database string, Collection string, addData bson.D) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	actionCollection := ConnDB.
		Database(Database).
		Collection(Collection)

	_, err := actionCollection.InsertOne(ctx, addData)
	if err != nil {
		log.Println(err)
	}
}
