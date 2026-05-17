package main

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var urls map[string]string = make(map[string]string)

var mu sync.RWMutex

var reverseURLs map[string]string = make(map[string]string)

type URLStruct struct {
	URL string
}

func main() {
	http.HandleFunc("POST /new", generateSlug)
	http.HandleFunc("/{url}", redirectFromShortened)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func generateSlug(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	decoder := json.NewDecoder(req.Body)
	var response URLStruct
	err := decoder.Decode(&response)

	if err != nil {
		statusCode := http.StatusInternalServerError
		msg := "500 - Interal server error"
		http.Error(w, msg, statusCode)
		return
	}

	ok, normalized := isValidURL(response.URL)

	if !ok {
		statusCode := http.StatusBadRequest
		msg := "400 - Not a valid URL"
		http.Error(w, msg, statusCode)
		return
	}

	response.URL = normalized

	if reverseURLs[response.URL] == "" {
		slug := rand.Text()[:5]

		for urls[slug] != "" {
			slug = rand.Text()[:5]
		}

		urls[slug] = response.URL
		reverseURLs[response.URL] = slug
	}

	url := URLStruct{URL: reverseURLs[response.URL]}

	j, err := json.Marshal(url)

	if err != nil {
		statusCode := http.StatusInternalServerError
		msg := "500 - Interal server error"
		http.Error(w, msg, statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func redirectFromShortened(w http.ResponseWriter, req *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	fullURL := urls[req.PathValue("url")]

	if fullURL == "" {
		statusCode := http.StatusNotFound
		msg := "404 - URL does not exist."
		http.Error(w, msg, statusCode)
		return
	}

	http.Redirect(w, req, fullURL, http.StatusTemporaryRedirect)
}

func isValidURL(urlToParse string) (bool, string) {
	if urlToParse == "" {
		return false, ""
	}

	u, err := url.Parse(urlToParse)

	if err != nil {
		return false, ""
	}

	if !strings.Contains(u.Host, ".") {
		return false, ""
	}

	if u.Scheme == "" {
		urlToParse = "https://" + urlToParse
		u, err = url.Parse(urlToParse)

		if err != nil {
			return false, ""
		}
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false, ""
	}

	if u.Host == "" {
		return false, ""
	}

	return true, urlToParse
}
