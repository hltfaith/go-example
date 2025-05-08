package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
)

type Info struct {
	ID string `json:"id"`
}

func main() {
	info := &Info{ID: "12345"}
	data, _ := json.Marshal(info)
	req, _ := http.NewRequest("POST", "http://httpbin.org/post", strings.NewReader(string(data)))
	req.Header.Add("Content-Type", "application/json")

	// 下面通过 Clone() 函数克隆request请求体
	req2 := req.Clone(req.Context())

	// 查看req原生的请求HTTP报文
	reqRaw, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(reqRaw))
	fmt.Println("---------------------------")

	// 查看克隆后的req请求HTTP报文
	req2Raw, _ := httputil.DumpRequest(req2, true)
	fmt.Println(string(req2Raw))
	fmt.Println("---------------------------")

	// 下面通过 Clone() 函数克隆request请求体, 补充Body内容后, 发送HTTP请求
	req2.Body = io.NopCloser(strings.NewReader(string(data)))
	client := http.Client{}
	res, err := client.Do(req2)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	b, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(b))
}
