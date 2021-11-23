package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// handleAPIPartner 함수는 파트너 관련 API다.
func handleAPIPartner(w http.ResponseWriter, r *http.Request) {
	// GET 메소드는 파트너의 id를 받아서 사용자 정보를 반환한다.
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		//mongoDB client 연결
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
		id := r.FormValue("id")
		if id == "" {
			http.Error(w, "id를 설정해주세요", http.StatusBadRequest)
			return
		}
		partner, err := getPartner(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// json 으로 결과 전송
		data, err := json.Marshal(partner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	} else {
		http.Error(w, "지원하지 않는 메소드입니다", http.StatusMethodNotAllowed)
		return
	}
}
