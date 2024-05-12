package main

import (
	"log"
	"net/http"
	"testing/fstest"
)

func main() {
	filename := "index.html"
	contents := []byte("<h1>帽儿山的枪手</h1>")
	fsys := fstest.MapFS{
		filename: {Data: contents},
	}
	http.Handle("/", http.FileServerFS(fsys))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
