package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AllStatusV2(client *mongo.Client) ([]Status, error) {
	collection := client.Database("setting").Collection("status")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"order", -1}})
	var status []Status
	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return status, err
	}
	err = cursor.All(ctx, &status)
	if err != nil {
		return status, err
	}
	return status, nil
}
