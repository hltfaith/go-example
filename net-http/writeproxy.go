package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Scheme = "http"
		r.WriteProxy(w)
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
