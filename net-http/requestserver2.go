// 服务端:文件上传
package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// 如果超过100字节使用临时文件来存储multipart/form中文件
		// 不超过则存入内存临时存储
		err := r.ParseMultipartForm(100)
		if err != nil {
			panic(err)
		}
		// MultipartForm解析多部分form内容
		m := r.MultipartForm
		for f := range m.File {
			// 取到文件信息
			file, fHeader, err := r.FormFile(f)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			out, err := os.Create("upload/" + fHeader.Filename)
			if err != nil {
				panic(err)
			}
			defer out.Close()
			_, err = io.Copy(out, file)
			if err != nil {
				panic(err)
			}
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
