package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

func main() {
	res, err := http.Get("http://httpbin.org/get")
	if err != nil {
		panic(err)
	}
	b, err := httputil.DumpResponse(res, true)
	if err != nil {
		panic(err)
	}
	// os.Stdout.Write(b)

	resp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(string(b))), res.Request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var bout bytes.Buffer
	_, err = io.Copy(&bout, resp.Body)
	fmt.Println(resp.StatusCode) // 响应报文中HTTP状态码
	fmt.Println(bout.String())   // 响应报文中的Body内容
}
