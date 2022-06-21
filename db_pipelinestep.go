package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func addPipelinestep(client *mongo.Client, s Pipelinestep) error {
	if s.Name == "" {
		return errors.New("빈 문자열입니다. 생성할 수 없습니다")
	}
	collection := client.Database(*flagDBName).Collection("pipelinestep")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	num, err := collection.CountDocuments(ctx, bson.M{"name": s.Name})
	if err != nil {
		return err
	}
	if num != 0 {
		return errors.New("같은 이름이 존재해서 추가할 수 없습니다")
	}
	_, err = collection.InsertOne(ctx, s)
	if err != nil {
		return err
	}
	return nil
}

func getPipelinestep(client *mongo.Client, id string) (Pipelinestep, error) {
	collection := client.Database(*flagDBName).Collection("pipelinestep")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s := Pipelinestep{}
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

func rmPipelinestep(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("pipelinestep")
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

func setPipelinestep(client *mongo.Client, s Pipelinestep) error {
	collection := client.Database(*flagDBName).Collection("pipelinestep")
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

func allPipelinesteps(client *mongo.Client) ([]Pipelinestep, error) {
	collection := client.Database(*flagDBName).Collection("pipelinestep")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var results []Pipelinestep
	opts := options.Find()
	opts.SetSort(bson.M{"name": 1})
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
