package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 10)
		_, err := io.Copy(ioutil.Discard, r.Body)
		if err != nil {
			panic(err)
		}
		io.WriteString(w, "200\n")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
