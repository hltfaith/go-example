package main

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"
)

func main() {
	url, err := url.Parse("http://google.com")
	if err != nil {
		panic(err)
	}

	os.Setenv("HTTP_PROXY", "http://127.0.0.1:7890")
	client := http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 跳过https
		},
	}

	req := http.Request{
		Method: "GET",
		URL:    url,
		Header: map[string][]string{
			"Proxy-Connection": {"keep-alive"},
		},
	}
	res, err := client.Do(&req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
}
