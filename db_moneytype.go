package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func addMoneytype(client *mongo.Client, s Moneytype) error {
	collection := client.Database(*flagDBName).Collection("moneytype")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, s)
	if err != nil {
		return err
	}
	return nil
}

func getMoneytype(client *mongo.Client, id string) (Moneytype, error) {
	collection := client.Database(*flagDBName).Collection("moneytype")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s := Moneytype{}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return s, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&s)
	if err != nil {
		return s, err
	}
	return s, nil
}

func rmMoneytype(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("moneytype")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	return nil
}

func setMoneytype(client *mongo.Client, s Moneytype) error {
	collection := client.Database(*flagDBName).Collection("moneytype")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": s.ID},
		bson.D{{Key: "$set", Value: s}},
	)
	if err != nil {
		return err
	}
	return nil
}
