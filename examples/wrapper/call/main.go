package main

import (
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("/Users/shuxian/Desktop/")))
	http.ListenAndServe(":9999", nil)
}
