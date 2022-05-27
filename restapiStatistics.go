package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func handleAPI1StatisticsShot(w http.ResponseWriter, r *http.Request) {
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// token 체크
	_, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	projects, err := client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	type Recipe struct {
		None    int64 `json:"none"`
		Hold    int64 `json:"hold"`
		Done    int64 `json:"done"`
		Out     int64 `json:"out"`
		Assign  int64 `json:"assign"`
		Ready   int64 `json:"ready"`
		Wip     int64 `json:"wip"`
		Confirm int64 `json:"confirm"`
		Omit    int64 `json:"omit"`
		Client  int64 `json:"client"`
	}
	rcp := Recipe{}
	shotFilter := bson.A{bson.D{{"type", "org"}}, bson.D{{"type", "left"}}}
	noneFilter := bson.D{{"status", NONE}, {"$or", shotFilter}}
	holdFilter := bson.D{{"status", HOLD}, {"$or", shotFilter}}
	doneFilter := bson.D{{"status", DONE}, {"$or", shotFilter}}
	outFilter := bson.D{{"status", OUT}, {"$or", shotFilter}}
	assignFilter := bson.D{{"status", ASSIGN}, {"$or", shotFilter}}
	readyFilter := bson.D{{"status", READY}, {"$or", shotFilter}}
	wipFilter := bson.D{{"status", WIP}, {"$or", shotFilter}}
	confirmFilter := bson.D{{"status", CONFIRM}, {"$or", shotFilter}}
	omitFilter := bson.D{{"status", OMIT}, {"$or", shotFilter}}
	clientFilter := bson.D{{"status", CLIENT}, {"$or", shotFilter}}

	for _, project := range projects {
		collection := client.Database("project").Collection(project)
		noneCount, err := collection.CountDocuments(ctx, noneFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		holdCount, err := collection.CountDocuments(ctx, holdFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		doneCount, err := collection.CountDocuments(ctx, doneFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		outCount, err := collection.CountDocuments(ctx, outFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		assignCount, err := collection.CountDocuments(ctx, assignFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		readyCount, err := collection.CountDocuments(ctx, readyFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		wipCount, err := collection.CountDocuments(ctx, wipFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		confirmCount, err := collection.CountDocuments(ctx, confirmFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		omitCount, err := collection.CountDocuments(ctx, omitFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		clientCount, err := collection.CountDocuments(ctx, clientFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		rcp.None += noneCount
		rcp.Hold += holdCount
		rcp.Done += doneCount
		rcp.Out += outCount
		rcp.Assign += assignCount
		rcp.Ready += readyCount
		rcp.Wip += wipCount
		rcp.Confirm += confirmCount
		rcp.Omit += omitCount
		rcp.Client += clientCount
	}

	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPIStatisticsProjectnum(w http.ResponseWriter, r *http.Request) {
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// token 체크
	_, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	projects, err := client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	type Recipe struct {
		Count    int      `json:"count"`
		Projects []string `json:"projects"`
	}
	rcp := Recipe{}
	rcp.Projects = projects
	rcp.Count = len(projects)

	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
