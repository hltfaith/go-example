package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			log.Panic("BasicAuth is none")
		}
		fmt.Fprintf(w, "username=%s, password=%s", username, password)
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
