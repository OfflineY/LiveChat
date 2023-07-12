package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

// Find 数据库内寻找
func Find(c *DatabaseConn, Collection string, Condition bson.M) ([]bson.M, error) {
	collection := c.Conn.Database(c.Name).Collection(Collection)

	find, err := collection.Find(context.Background(), Condition)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = find.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
