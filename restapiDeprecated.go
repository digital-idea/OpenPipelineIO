package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"
)

// 검색 RestAPI
//
// Deprecated: handleAPIItems는 더 이상 사용하지 않는다. handleApi2Items를 사용할 것.
func handleAPIItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	defer session.Close()
	q := r.URL.Query()
	op := SearchOption{
		Project:    q.Get("project"),
		Searchword: q.Get("searchword"),
		Sortkey:    q.Get("sortkey"),
		Assign:     str2bool(q.Get("assign")),
		Ready:      str2bool(q.Get("ready")),
		Wip:        str2bool(q.Get("wip")),
		Confirm:    str2bool(q.Get("confirm")),
		Done:       str2bool(q.Get("done")),
		Omit:       str2bool(q.Get("omit")),
		Hold:       str2bool(q.Get("hold")),
		Out:        str2bool(q.Get("out")),
		None:       str2bool(q.Get("none")),
		Shot:       str2bool(q.Get("shot")),
		Assets:     str2bool(q.Get("shot")),
		Type3d:     str2bool(q.Get("type3d")),
		Type2d:     str2bool(q.Get("type2d")),
	}
	result, err := Searchv1(session, op)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	if len(result) == 0 {
		fmt.Fprintln(w, "{\"error\":\"검색결과 0건\"}")
		return
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}
