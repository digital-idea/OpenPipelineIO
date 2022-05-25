package main

import (
	"net/http"
)

func handleAPIStatusInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Help", "legacy status option")
	if r.Method == http.MethodOptions {
		return
	}
	q := r.URL.Query()
	reverse := str2bool(q.Get("reverse"))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if reverse {
		reverseStatusInfo := `{"client":"9","omit":"8","confirm":"7","wip":"6","ready":"5","assign":"4","out":"3","done":"2","hold":"1","none":"0"}`
		w.Write([]byte(reverseStatusInfo))
	} else {
		statusInfo := `{"9":"client","8":"omit","7":"confirm","6":"wip","5":"ready","4":"assign","3":"out","2":"done","1":"hold","0":"none"}`
		w.Write([]byte(statusInfo))
	}
}
