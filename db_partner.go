// 이 코드는 파트너사와 관련된 DBapi가 모여있는 파일입니다.
package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func addPartner(client *mongo.Client, p Partner) error {
	p.CreateTime = time.Now().Format(time.RFC3339)
	p.UpdateTime = p.CreateTime
	collection := client.Database("csi").Collection("partners")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// getPartner 함수는 사용자를 가지고오는 함수이다.
func getPartner(client *mongo.Client, id string) (Partner, error) {
	collection := client.Database("csi").Collection("partners")
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

// rmPartner 함수는 파트너사 정보를 삭제하는 함수다.
func rmPartner(client *mongo.Client, id string) error {
	collection := client.Database("csi").Collection("partners")
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

// setPartner 함수는 파트너사 정보를 수정하는 함수다.
func setPartner(client *mongo.Client, p Partner) error {
	p.UpdateTime = time.Now().Format(time.RFC3339)

	collection := client.Database("csi").Collection("partners")
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

// allPartners 함수는 모든 파트너사 정보를 가져오는 함수다.
func allPartners(client *mongo.Client) ([]Partner, error) {
	collection := client.Database("csi").Collection("partners")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result []Partner

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return result, err
	}
	err = cursor.All(ctx, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}