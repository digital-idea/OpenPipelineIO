package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func addProjectForPartner(client *mongo.Client, s ProjectForPartner) error {
	collection := client.Database(*flagDBName).Collection("projectforpartner")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, s)
	if err != nil {
		return err
	}
	return nil
}

func getProjectForPartner(client *mongo.Client, id string) (ProjectForPartner, error) {
	collection := client.Database(*flagDBName).Collection("projectforpartner")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s := ProjectForPartner{}
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

func rmProjectForPartner(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("projectforpartner")
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

func setProjectForPartner(client *mongo.Client, s ProjectForPartner) error {
	collection := client.Database(*flagDBName).Collection("projectforpartner")
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

func allProjectForPartners(client *mongo.Client) ([]ProjectForPartner, error) {
	collection := client.Database(*flagDBName).Collection("projectforpartner")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var results []ProjectForPartner
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}
