package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	os.Setenv("HTTP_PROXY", "http://127.0.0.1:12345")
	req, err := http.NewRequest("GET", "http://example.com", nil)

	if err != nil {
		panic(err)
	}
	url, err := http.ProxyFromEnvironment(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(url)
}
