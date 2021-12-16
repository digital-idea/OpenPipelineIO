package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// handleAPIPartner 함수는 파트너 관련 API다.
func handleAPIPartner(w http.ResponseWriter, r *http.Request) {
	// GET 메소드는 파트너의 id를 받아서 파트너 정보를 반환한다.
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
	} else if r.Method == http.MethodPost {
		// POST 메소드는 새로운 파트너 정보를 DB에 저장한다
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
		// 파트너 생성
		p := Partner{}
		p.ID = primitive.NewObjectID()
		// 파트너 정보 Parsing
		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		name := r.FormValue("name")
		if name == "" {
			http.Error(w, "name을 설정해주세요", http.StatusBadRequest)
			return
		}
		sizeStr := r.FormValue("size")
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s 는 size로 사용할 수 없는 값 입니다", sizeStr), http.StatusBadRequest)
			return
		}
		domain := r.FormValue("domain")
		homepage := r.FormValue("homepage")
		address := r.FormValue("address")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		timezone := r.FormValue("timezone")
		description := r.FormValue("description")
		p.Name = name
		p.Size = size
		p.Domain = domain
		p.Homepage = homepage
		p.Address = address
		p.Phone = phone
		p.Email = email
		p.Timezone = timezone
		p.Description = description
		// 파트너 추가
		err = addPartner(client, p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Response
		item, err := getPartner(client, p.ID.Hex())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	} else if r.Method == http.MethodPut {
		// PUT 메소드는 수정된 파트너 정보를 DB에 저장한다
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
		// 파트너 정보 Parsing
		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := r.FormValue("id")
		p, err := getPartner(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		name := r.FormValue("name")
		if name == "" {
			http.Error(w, "name을 설정해주세요", http.StatusBadRequest)
			return
		}
		sizeStr := r.FormValue("size")
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s 는 size로 사용할 수 없는 값 입니다", sizeStr), http.StatusBadRequest)
			return
		}
		domain := r.FormValue("domain")
		homepage := r.FormValue("homepage")
		address := r.FormValue("address")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		timezone := r.FormValue("timezone")
		description := r.FormValue("description")
		p.Name = name
		p.Size = size
		p.Domain = domain
		p.Homepage = homepage
		p.Address = address
		p.Phone = phone
		p.Email = email
		p.Timezone = timezone
		p.Description = description
		// 수정된 파트너 정보 DB에 저장
		err = setPartner(client, p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Response
		item, err := getPartner(client, p.ID.Hex())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
