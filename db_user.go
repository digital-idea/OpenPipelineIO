package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// getUser 함수는 사용자를 가지고오는 함수이다.
func getUserV2(client *mongo.Client, id string) (User, error) {
	collection := client.Database("user").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u := User{}
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&u)
	if err != nil {
		return u, err
	}
	return u, nil
}
