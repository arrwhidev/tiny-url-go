package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/arrwhidev/base-converter"
)

const offset = uint64(300)

var db = make(map[uint64]string)

type jsonInterface struct {
	URL string `json:"url"`
}

func shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Expecting POST request", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	// Decode JSON.
	var jsonRequest jsonInterface
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&jsonRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusUnprocessableEntity)
		return
	}

	// Insert new row into database.
	nextID := offset + uint64(len(db)) + 1
	db[nextID] = jsonRequest.URL

	encoded := radix.Encode(uint64(nextID), radix.Base62)

	// Write JSON response object
	encoder := json.NewEncoder(w)
	encoder.Encode(&jsonInterface{
		URL: "http://localhost:8080/" + encoded,
	})
}

func expand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Expecting GET request", http.StatusBadRequest)
		return
	}

	shortCode := strings.Split(r.RequestURI, "/")[1]
	id := radix.Decode(shortCode, radix.Base62)

	if url, ok := db[id]; ok {
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	} else {
		http.Error(w, "Unecognised short code: "+shortCode, http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/" {
			shorten(w, r)
		} else {
			expand(w, r)
		}
	})

	fmt.Println("Listening on :8080...")
	srv := &http.Server{Addr: ":8080"}
	srv.ListenAndServe()
}
