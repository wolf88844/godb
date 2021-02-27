package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/index", indexHandler)

	err := http.ListenAndServe(":3013", nil)
	fmt.Println(err)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/index==============")
	w.Write([]byte("这是默认首页"))
}
