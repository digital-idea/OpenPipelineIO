package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func addItemV2(client *mongo.Client, i Item) error {
	err := i.CheckError()
	if err != nil {
		return err
	}
	collection := client.Database("project").Collection(i.Project)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	num, err := collection.CountDocuments(ctx, bson.M{"name": i.ID})
	if err != nil {
		return err
	}
	if num != 0 {
		return errors.New("같은 이름을 가진 데이터가 있습니다")
	}
	_, err = collection.InsertOne(ctx, i)
	if err != nil {
		return err
	}
	return nil
}
