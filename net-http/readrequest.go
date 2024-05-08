package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ReqData struct {
	Raw     string        // 实际的HTTP GET请求报文
	Req     *http.Request // 展示请求实际GET请求的传参
	Body    string
	Trailer http.Header
	Error   string
}

func main() {
	var noTrailer http.Header = nil
	var noError = ""
	reqdata := &ReqData{
		"GET /get HTTP/1.1\r\n" +
			"Host: httpbin.org\r\n" +
			"User-Agent: Fake\r\n" +
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n" +
			"Accept-Language: en-us,en;q=0.5\r\n" +
			"Accept-Encoding: gzip,deflate\r\n" +
			"Accept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.7\r\n" +
			"Keep-Alive: 300\r\n" +
			"Content-Length: 18\r\n" +
			"Proxy-Connection: keep-alive\r\n\r\n" +
			"帽儿山的枪手\n???",

		&http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme: "http",
				Host:   "httpbin.org",
				Path:   "/",
			},
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: http.Header{
				"Accept":           {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
				"Accept-Language":  {"en-us,en;q=0.5"},
				"Accept-Encoding":  {"gzip,deflate"},
				"Accept-Charset":   {"ISO-8859-1,utf-8;q=0.7,*;q=0.7"},
				"Keep-Alive":       {"300"},
				"Proxy-Connection": {"keep-alive"},
				"Content-Length":   {"18"},
				"User-Agent":       {"Fake"},
			},
			Close:         false,
			ContentLength: 18,
			Host:          "httpbin.org",
			RequestURI:    "http://httpbin.org/get",
		},

		"帽儿山的枪手\n",
		noTrailer,
		noError,
	}

	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(reqdata.Raw)))
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	var bout bytes.Buffer
	_, err = io.Copy(&bout, req.Body)
	if err != nil {
		panic(err)
	}
	body := bout.String()
	fmt.Println(req.Method) // HTTP请求类型为 GET
	fmt.Println(body)       // body内容输出   帽儿山的枪手
}
