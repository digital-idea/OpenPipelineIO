package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func addPartner(client *mongo.Client, p Partner) error {
	if p.Name == "" {
		return errors.New("빈 문자열입니다. 파트너를 생성할 수 없습니다")
	}
	collection := client.Database(*flagDBName).Collection("partner")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	num, err := collection.CountDocuments(ctx, bson.M{"name": p.Name})
	if err != nil {
		return err
	}
	if num != 0 {
		return errors.New("같은 이름의 파트너사가 존재해서 파트너사를 추가할 수 없습니다")
	}
	_, err = collection.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

func getPartner(client *mongo.Client, id string) (Partner, error) {
	collection := client.Database(*flagDBName).Collection("partner")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	p := Partner{}
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

func rmPartner(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("partner")
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

func setPartner(client *mongo.Client, p Partner) error {
	collection := client.Database(*flagDBName).Collection("partner")
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

func allPartners(client *mongo.Client) ([]Partner, error) {
	collection := client.Database(*flagDBName).Collection("partner")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var partners []Partner
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return partners, err
	}
	err = cursor.All(ctx, &partners)
	if err != nil {
		return partners, err
	}
	return partners, nil
}
