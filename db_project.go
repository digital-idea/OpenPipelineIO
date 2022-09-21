package main

import (
	"context"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProjectlistV2(client *mongo.Client) ([]string, error) {
	var projects []string
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	projects, err := client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
	if err != nil {
		return projects, err
	}
	sort.Strings(projects)
	return projects, nil
}
