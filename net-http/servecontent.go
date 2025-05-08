package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		file := "servecontent.go"
		fileBytes, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		mime := http.DetectContentType(fileBytes)
		fileSize := len(string(fileBytes))
		w.Header().Set("Content-Type", mime)
		w.Header().Set("Content-Disposition", "attachment; filename="+file)
		w.Header().Set("Content-Length", strconv.Itoa(fileSize))

		http.ServeContent(w, r, file, time.Now(), bytes.NewReader(fileBytes))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
