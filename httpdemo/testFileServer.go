package main

import "net/http"

func main() {
	testFileServer()
}

func testFileServer() {
	http.ListenAndServe(":2003", http.FileServer(http.Dir("d:/")))
}
