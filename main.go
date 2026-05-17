package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var urls map[string]string = make(map[string]string)

var reverse_urls map[string]string = make(map[string]string)

type UrlStruct struct {
	Url string
}

func main() {
	http.HandleFunc("/new", generateShortened)
	http.HandleFunc("/{url}", redirectFromShortened)

	http.ListenAndServe(":3000", nil)
}

func generateShortened(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var response UrlStruct
	err := decoder.Decode(&response)

	if err != nil {
		log.Fatal(err)
	}

	slug := rand.Text()[:5]

	if reverse_urls[response.Url] == "" {
		urls[slug] = response.Url
		reverse_urls[response.Url] = slug
	}

	fmt.Fprint(w, slug)
}

func redirectFromShortened(w http.ResponseWriter, req *http.Request) {
	fullUrl := urls[req.URL.Path[1:]]

	if fullUrl == "" {
		fmt.Fprint(w, "404 - URL does not exists")
	}

	fmt.Fprint(w, "Redirecting to ", fullUrl)
}
