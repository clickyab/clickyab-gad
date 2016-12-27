package main

import (
	"time"

	"fastcgi"
	"net/http"
)

func main() {
	handler := fcgi.NewPHPFastCGIHandler("/home/develop/gad/clickyab-server/a/", "/", "127.0.0.1:9999", 30*time.Second, 30*time.Second, 30*time.Second)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})

	http.ListenAndServe(":80", nil)
}
