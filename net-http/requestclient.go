// 客户端: 文件上传
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	City     string `json:"city"`
}

func main() {
	userinfos := make([]*UserInfo, 0)
	userinfos = append(userinfos,
		&UserInfo{
			ID:       "001",
			Username: "帽儿山的枪手",
			City:     "北京",
		},
		&UserInfo{
			ID:       "002",
			Username: "changhao",
			City:     "北京",
		})
	b, err := json.Marshal(&userinfos)
	if err != nil {
		panic(err)
	}
	// 封装请求体
	req := &http.Request{
		// POST请求方法
		Method: "POST",
		// 拼接url路径 http://127.0.0.1/post
		URL: &url.URL{
			Scheme: "http",
			Host:   "127.0.0.1",
			Path:   "/post",
		},
		// http协议版本
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		// 请求头
		Header: map[string][]string{
			"Content-Type": {"application/json"}, // 使用json格式载荷体
		},
		// body内容
		Body: io.NopCloser(strings.NewReader(string(b))),
		// 内容长度
		ContentLength: int64(len(b)),
		// 传输编码, 采用分块
		// TransferEncoding: []string{"chunked"},
		// 服务地址
		Host: "127.0.0.1",
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	reader := bufio.NewReader(res.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			fmt.Print(string(line))
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
	}

}
