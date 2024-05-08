package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type RespData struct {
	Raw  string        // 实际的HTTP GET请求的响应报文
	Resp http.Response // 展示实际响应的传参
	Body string
}

func main() {
	respdata := &RespData{
		"HTTP/1.0 200 OK\r\n" +
			"Connection: close\r\n" +
			"\r\n" +
			"帽儿山的枪手\n",

		http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.0",
			ProtoMajor: 1,
			ProtoMinor: 0,
			Request:    &http.Request{Method: "GET"},
			Header: http.Header{
				"Connection": {"close"},
			},
			Close:         true,
			ContentLength: -1,
		},

		"帽儿山的枪手\n",
	}

	resp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(respdata.Raw)), respdata.Resp.Request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var bout bytes.Buffer
	_, err = io.Copy(&bout, resp.Body)
	fmt.Println(resp.StatusCode) // 响应报文中HTTP状态码  200
	fmt.Println(bout.String())   // 响应报文中的Body内容  帽儿山的枪手
}
