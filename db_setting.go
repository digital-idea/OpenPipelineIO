package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AllTaskSettingsV2(client *mongo.Client) ([]Tasksetting, error) {
	results := []Tasksetting{}
	collection := client.Database("setting").Collection("tasksetting")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.Find()
	opts.SetSort(bson.M{"id": 1})
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}

func TaskSettingNamesV2(client *mongo.Client) ([]string, error) {
	var results []string
	collection := client.Database("setting").Collection("tasksetting")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	values, err := collection.Distinct(ctx, "name", bson.D{})
	if err != nil {
		return results, err
	}
	for _, value := range values {
		results = append(results, fmt.Sprintf("%v", value))
	}
	return results, nil
}
