package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func addPipelinestep(client *mongo.Client, p Pipelinestep) error {
	if p.Name == "" {
		return errors.New("빈 문자열입니다. Pipelinestep을 생성할 수 없습니다")
	}
	collection := client.Database(*flagDBName).Collection("pipelinestep")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	num, err := collection.CountDocuments(ctx, bson.M{"name": p.Name})
	if err != nil {
		return err
	}
	if num != 0 {
		return errors.New("같은 이름의 Pipelinestep이 존재해서 추가할 수 없습니다")
	}
	_, err = collection.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

func getPipelinestep(client *mongo.Client, id string) (Pipelinestep, error) {
	collection := client.Database(*flagDBName).Collection("pipelinestep")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	p := Pipelinestep{}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return p, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&p)
	if err != nil {
		return p, err
	}
	return p, nil
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

func setPipelinestep(client *mongo.Client, p Pipelinestep) error {
	collection := client.Database(*flagDBName).Collection("pipelinestep")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": p.ID},
		bson.D{{Key: "$set", Value: p}},
	)
	if err != nil {
		return err
	}
	return nil
}

func allPipelinestep(client *mongo.Client) ([]Pipelinestep, error) {
	collection := client.Database(*flagDBName).Collection("pipelinestep")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var pipelinesteps []Pipelinestep
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return pipelinesteps, err
	}
	err = cursor.All(ctx, &pipelinesteps)
	if err != nil {
		return pipelinesteps, err
	}
	return pipelinesteps, nil
}
