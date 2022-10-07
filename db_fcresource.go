package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func addFCResource(client *mongo.Client, e FullCalendarResource) error {
	collection := client.Database("OpenPipelineIO").Collection("fcresource")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, e)
	if err != nil {
		return err
	}
	return nil
}

func getFCResource(client *mongo.Client, id string) (FullCalendarResource, error) {
	collection := client.Database("OpenPipelineIO").Collection("fcresource")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fcResource := FullCalendarResource{}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fcResource, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&fcResource)
	if err != nil {
		return fcResource, err
	}
	return fcResource, nil
}

func rmFCResource(client *mongo.Client, id string) error {
	collection := client.Database("OpenPipelineIO").Collection("fcresource")
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

func setFCResource(client *mongo.Client, r FullCalendarResource) error {
	collection := client.Database("OpenPipelineIO").Collection("fcresource")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": r.ID},
		bson.D{{Key: "$set", Value: r}},
	)
	if err != nil {
		return err
	}
	return nil
}
