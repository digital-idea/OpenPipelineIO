package main

import (
	"net/http"
)

// handleErrorChaptcha 핸들러는 Chaptcha 코드가 다를 때 출력하는 웹페이지이다.
func handleErrorChaptcha(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := TEMPLATES.ExecuteTemplate(w, "error-chaptcha", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
