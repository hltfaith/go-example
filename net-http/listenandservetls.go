package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, TLS!\n")
	})
	log.Fatal(http.ListenAndServeTLS("", "cert.pem", "key.pem", nil))
}
