package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Find 数据库内寻找
func Find(Conn *mongo.Client, Database string, Collection string, Condition bson.M) ([]bson.M, error) {
	collection := Conn.Database(Database).Collection(Collection)

	find, err := collection.Find(context.Background(), Condition)
	if err != nil {
		//log.Println(err)
		return nil, err
	}

	var results []bson.M
	if err = find.All(context.TODO(), &results); err != nil {
		//log.Println(err)
		return nil, err
	}

	return results, nil
}
