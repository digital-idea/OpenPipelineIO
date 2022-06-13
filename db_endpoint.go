package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func addEndpoint(client *mongo.Client, s Endpoint) error {
	if s.Endpoint == "" {
		return errors.New("nead endpoint")
	}
	collection := client.Database(*flagDBName).Collection("endpoint")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	num, err := collection.CountDocuments(ctx, bson.M{"endpoint": s.Endpoint})
	if err != nil {
		return err
	}
	if num != 0 {
		return errors.New("같은 이름을 가진 데이터가 있습니다")
	}
	_, err = collection.InsertOne(ctx, s)
	if err != nil {
		return err
	}
	return nil
}

func getEndpoint(client *mongo.Client, id string) (Endpoint, error) {
	collection := client.Database(*flagDBName).Collection("endpoint")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s := Endpoint{}
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

func rmEndpoint(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("endpoint")
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

func setEndpoint(client *mongo.Client, s Endpoint) error {
	collection := client.Database(*flagDBName).Collection("endpoint")
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

func allEndpoints(client *mongo.Client) ([]Endpoint, error) {
	collection := client.Database(*flagDBName).Collection("endpoint")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var results []Endpoint
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
