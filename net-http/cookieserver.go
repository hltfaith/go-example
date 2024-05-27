package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		// 浏览器获取本地存储中是否有cookie
		c, err := r.Cookie("ClientCookieID")
		if err != nil {
			log.Println(err)
		}
		if c != nil {
			w.Write([]byte(c.Value))
			return
		}

		// 服务端响应返回cookie信息
		http.SetCookie(w, &http.Cookie{
			Name:    "ClientCookieID",
			Value:   "12345",
			Expires: time.Now().Add(120 * time.Second),
		})
		w.Write([]byte("This is cookie\n"))
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
