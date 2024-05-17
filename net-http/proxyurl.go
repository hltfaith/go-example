package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	url, err := url.Parse("http://188.68.176.2:8080")
	if err != nil {
		panic(err)
	}
	client := http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(url),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	res, err := client.Get("http://baidu.com")
	if err != nil {
		panic(err)
	}
	b, _ := httputil.DumpRequest(res.Request, false)
	fmt.Println(string(b))
}
