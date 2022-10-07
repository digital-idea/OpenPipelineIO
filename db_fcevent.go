package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func addFCEvent(client *mongo.Client, e FullCalendarEvent) error {
	collection := client.Database("OpenPipelineIO").Collection("fcevent")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, e)
	if err != nil {
		return err
	}
	return nil
}

func getFCEvent(client *mongo.Client, id string) (FullCalendarEvent, error) {
	collection := client.Database("OpenPipelineIO").Collection("fcevent")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fc := FullCalendarEvent{}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fc, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&fc)
	if err != nil {
		return fc, err
	}
	return fc, nil
}

func rmFCEvent(client *mongo.Client, id string) error {
	collection := client.Database("OpenPipelineIO").Collection("fcevent")
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

func setFCEvent(client *mongo.Client, e FullCalendarEvent) error {
	collection := client.Database("OpenPipelineIO").Collection("fcevent")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": e.ID},
		bson.D{{Key: "$set", Value: e}},
	)
	if err != nil {
		return err
	}
	return nil
}
