// 客户端:文件上传
package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

func createBody(filenames []string) (string, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	// 生成multipart和随机的分界线
	w := multipart.NewWriter(buf)
	for _, file := range filenames {
		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		// 创建文件并携带MIME的头部信息
		fd, _ := w.CreateFormFile(file, file)
		io.Copy(fd, f)
	}
	w.Close()
	// 这里的 FormDataContentType 则是随机生成分界符用来区分文件消息
	// 例如 multipart/form-data; boundary=7cfa31806e7151431dffe1d1d086eaaefbc2dbe5a61ced7c2bd8f51db01c
	return w.FormDataContentType(), buf
}

func main() {
	// 准备两个二进制文件
	// dd if=/dev/zero of=file1.bin bs=10240 count=1024
	// dd if=/dev/zero of=file2.bin bs=10240 count=1024
	filesname := []string{"file1.bin", "file2.bin"}
	contenetType, buf := createBody(filesname)
	req := &http.Request{
		// POST请求方法
		Method: "POST",
		// 拼接url路径 http://127.0.0.1:8080/upload
		URL: &url.URL{
			Scheme: "http",
			Host:   "127.0.0.1:8080",
			Path:   "/upload",
		},
		// http协议版本
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		// 请求头
		Header: map[string][]string{
			"Content-Type": {contenetType},
		},
		// body内容
		Body: io.NopCloser(buf),
		// 默认不声明ContentLength默认机制会使用Content-Type: chunked分块
		// ContentLength: 0,
		// 服务地址
		Host: "127.0.0.1",
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
}
