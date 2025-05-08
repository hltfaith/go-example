// 服务端: 文件上传
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	City     string `json:"city"`
}

func main() {
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			panic("expected http.ResponseWriter to be an http.Flusher")
		}
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("Connection", "Keep-Alive")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		b, _ := io.ReadAll(r.Body)
		userinfos := []*UserInfo{}
		json.Unmarshal(b, &userinfos)

		for _, info := range userinfos {
			wbody := []byte("Chunk: " + info.Username + "\n")
			_, err := w.Write(wbody)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			flusher.Flush()
		}
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
