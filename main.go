package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", handlePing)

	http.ListenAndServe(":3000", nil)
}

func handlePing(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "pong")
}
