package main

import (
	"fmt"
	"net/http"
)

func main() {
	// image/png
	fmt.Println(http.DetectContentType([]byte("\x89PNG\x0D\x0A\x1A\x0A")))
	// image/jpeg
	fmt.Println(http.DetectContentType([]byte("\xFF\xD8\xFF")))
}
