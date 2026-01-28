package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)

	fmt.Println("Server started at http://localhost:8080/hello")
	http.ListenAndServe(":8080", mux)
}
