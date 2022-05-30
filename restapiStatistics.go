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
	q := r.URL.Query()
	project := q.Get("project")
	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

func handleAPI1StatisticsTag(w http.ResponseWriter, r *http.Request) {
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
	q := r.URL.Query()
	project := q.Get("project")
	tagName := q.Get("name")
	typ := q.Get("type")
	if typ == "" {
		typ = "shot"
	}
	if !(typ == "shot" || typ == "asset") {
		http.Error(w, "The type value must be either 'shot' or 'asset' value.", http.StatusBadRequest)
		return
	}
	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
	typeFilter := bson.A{bson.D{{"type", "org"}}, bson.D{{"type", "left"}}}
	if typ == "asset" {
		typeFilter = bson.A{bson.D{{"type", "asset"}}}
	}
	noneFilter := bson.D{{"status", NONE}, {"tag", tagName}, {"$or", typeFilter}}
	holdFilter := bson.D{{"status", HOLD}, {"tag", tagName}, {"$or", typeFilter}}
	doneFilter := bson.D{{"status", DONE}, {"tag", tagName}, {"$or", typeFilter}}
	outFilter := bson.D{{"status", OUT}, {"tag", tagName}, {"$or", typeFilter}}
	assignFilter := bson.D{{"status", ASSIGN}, {"tag", tagName}, {"$or", typeFilter}}
	readyFilter := bson.D{{"status", READY}, {"tag", tagName}, {"$or", typeFilter}}
	wipFilter := bson.D{{"status", WIP}, {"tag", tagName}, {"$or", typeFilter}}
	confirmFilter := bson.D{{"status", CONFIRM}, {"tag", tagName}, {"$or", typeFilter}}
	omitFilter := bson.D{{"status", OMIT}, {"tag", tagName}, {"$or", typeFilter}}
	clientFilter := bson.D{{"status", CLIENT}, {"tag", tagName}, {"$or", typeFilter}}

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

func handleAPI1StatisticsTask(w http.ResponseWriter, r *http.Request) {
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
	q := r.URL.Query()
	project := q.Get("project")
	task := q.Get("task")
	typ := q.Get("type")
	if typ == "" {
		typ = "shot"
	}
	if !(typ == "shot" || typ == "asset") {
		http.Error(w, "The type value must be either 'shot' or 'asset' value.", http.StatusBadRequest)
		return
	}
	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
	typeFilter := bson.A{bson.D{{"type", "org"}}, bson.D{{"type", "left"}}}
	if typ == "asset" {
		typeFilter = bson.A{bson.D{{"type", "asset"}}}
	}
	noneFilter := bson.D{{"tasks." + task + ".status", NONE}, {"$or", typeFilter}}
	holdFilter := bson.D{{"tasks." + task + ".status", HOLD}, {"$or", typeFilter}}
	doneFilter := bson.D{{"tasks." + task + ".status", DONE}, {"$or", typeFilter}}
	outFilter := bson.D{{"tasks." + task + ".status", OUT}, {"$or", typeFilter}}
	assignFilter := bson.D{{"tasks." + task + ".status", ASSIGN}, {"$or", typeFilter}}
	readyFilter := bson.D{{"tasks." + task + ".status", READY}, {"$or", typeFilter}}
	wipFilter := bson.D{{"tasks." + task + ".status", WIP}, {"$or", typeFilter}}
	confirmFilter := bson.D{{"tasks." + task + ".status", CONFIRM}, {"$or", typeFilter}}
	omitFilter := bson.D{{"tasks." + task + ".status", OMIT}, {"$or", typeFilter}}
	clientFilter := bson.D{{"tasks." + task + ".status", CLIENT}, {"$or", typeFilter}}

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

func handleAPI2StatisticsShot(w http.ResponseWriter, r *http.Request) {
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
	q := r.URL.Query()
	project := q.Get("project")
	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	// 모든 상태를 가지고 옵니다.
	status, err := AllStatusV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	type Recipe struct {
		Status  map[string]int64  `json:"status"`
		Filters map[string]bson.D `json:"-"` // 통계를 위한 필터저장에만 사용한다. 반환하지 않는다.
	}
	rcp := Recipe{}
	rcp.Status = make(map[string]int64)
	rcp.Filters = make(map[string]bson.D)
	// filter를 생성합니다.
	shotFilter := bson.A{bson.D{{"type", "org"}}, bson.D{{"type", "left"}}}
	for _, s := range status {
		rcp.Filters[s.ID] = bson.D{{"statusv2", s.ID}, {"$or", shotFilter}}
	}

	for _, project := range projects {
		collection := client.Database("project").Collection(project)
		// filter를 for 문 돌면서 나오는 카운트를 검색하고 상태에 넣는다.
		for status, filter := range rcp.Filters {
			count, err := collection.CountDocuments(ctx, filter)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			rcp.Status[status] += count
		}
	}

	data, err := json.Marshal(rcp.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPI2StatisticsAsset(w http.ResponseWriter, r *http.Request) {
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
	q := r.URL.Query()
	project := q.Get("project")
	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	// 모든 상태를 가지고 옵니다.
	status, err := AllStatusV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	type Recipe struct {
		Status  map[string]int64  `json:"status"`
		Filters map[string]bson.D `json:"-"` // 통계를 위한 필터저장에만 사용한다. 반환하지 않는다.
	}
	rcp := Recipe{}
	rcp.Status = make(map[string]int64)
	rcp.Filters = make(map[string]bson.D)
	// filter를 생성합니다.
	for _, s := range status {
		rcp.Filters[s.ID] = bson.D{{"statusv2", s.ID}, {"type", "asset"}}
	}

	for _, project := range projects {
		collection := client.Database("project").Collection(project)
		// filter를 for 문 돌면서 나오는 카운트를 검색하고 상태에 넣는다.
		for status, filter := range rcp.Filters {
			count, err := collection.CountDocuments(ctx, filter)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			rcp.Status[status] += count
		}
	}

	data, err := json.Marshal(rcp.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPI1StatisticsAsset(w http.ResponseWriter, r *http.Request) {
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
	q := r.URL.Query()
	project := q.Get("project")
	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
	noneFilter := bson.D{{"status", NONE}, {"type", "asset"}}
	holdFilter := bson.D{{"status", HOLD}, {"type", "asset"}}
	doneFilter := bson.D{{"status", DONE}, {"type", "asset"}}
	outFilter := bson.D{{"status", OUT}, {"type", "asset"}}
	assignFilter := bson.D{{"status", ASSIGN}, {"type", "asset"}}
	readyFilter := bson.D{{"status", READY}, {"type", "asset"}}
	wipFilter := bson.D{{"status", WIP}, {"type", "asset"}}
	confirmFilter := bson.D{{"status", CONFIRM}, {"type", "asset"}}
	omitFilter := bson.D{{"status", OMIT}, {"type", "asset"}}
	clientFilter := bson.D{{"status", CLIENT}, {"type", "asset"}}

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
