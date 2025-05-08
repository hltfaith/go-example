package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	req := &http.Request{
		// POST请求方法
		Method: "POST",
		// 拼接url路径 http://127.0.0.1/post
		URL: &url.URL{
			Scheme: "http",
			Host:   "127.0.0.1",
			Path:   "/post",
		},
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
	}
	req.SetBasicAuth("changhao", "123456")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	fmt.Println(string(b))
}
