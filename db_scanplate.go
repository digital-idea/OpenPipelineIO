package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func addScanPlate(client *mongo.Client, s ScanPlate) error {
	s.ProcessStatus = "wait"
	s.CreateTime = time.Now().Format(time.RFC3339)
	collection := client.Database("OpenPipelineIO").Collection("scanplate")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, s)
	if err != nil {
		return err
	}
	return nil
}

func getScanPlate(client *mongo.Client, id string) (ScanPlate, error) {
	collection := client.Database("OpenPipelineIO").Collection("scanplate")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s := ScanPlate{}
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

func rmScanPlate(client *mongo.Client, id string) error {
	collection := client.Database("OpenPipelineIO").Collection("scanplate")
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

func setScanPlate(client *mongo.Client, s ScanPlate) error {
	collection := client.Database("OpenPipelineIO").Collection("scanplate")
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

func SetScanPlateProcessStatus(client *mongo.Client, id string, status string) error {
	collection := client.Database("OpenPipelineIO").Collection("scanplate")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{"processstatus": status},
	}
	err = collection.FindOneAndUpdate(ctx, filter, update).Err()
	if err != nil {
		return err
	}
	return nil
}

func allScanPlate(client *mongo.Client) ([]ScanPlate, error) {
	collection := client.Database("OpenPipelineIO").Collection("scanplate")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var results []ScanPlate
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

// SetScanPlateErrStatus 함수는 인수로 받은 Status를 error status로 바꾼다
func SetScanPlateErrStatus(client *mongo.Client, id, errmsg string) error {
	collection := client.Database("OpenPipelineIO").Collection("scanplate")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// item의 Status를 업데이트 한다.
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}

	update := bson.M{
		"$set": bson.M{
			"processstatus": "error",
		},
		"$push": bson.M{"logs": errmsg},
	}
	err = collection.FindOneAndUpdate(ctx, filter, update).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetWaitProcessStatusScanPlate() (ScanPlate, error) {
	var result ScanPlate
	//mongoDB client 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		return result, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return result, err
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return result, err
	}
	collection := client.Database("OpenPipelineIO").Collection("scanplate")
	filter := bson.M{"processstatus": "wait"}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func GetUnDoneScanPlate(client *mongo.Client) ([]ScanPlate, error) {
	var results []ScanPlate
	collection := client.Database("OpenPipelineIO").Collection("scanplate")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"status": bson.M{"$ne": "done"}} // done이 아닌 리스트를 가지고 온다.
	opts := options.Find()
	opts.SetSort(bson.M{"name": 1})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}
