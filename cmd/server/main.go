package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
)

var mapURL = map[string]string{}

func mainPage(w http.ResponseWriter, r *http.Request) {
	//TODO: handler "/"
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
	} else {
		//TODO: Обработать body и сократить
		var body []byte
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		hasher := sha1.New()
		hasher.Write(body)
		hash := hex.EncodeToString(hasher.Sum(nil))
		mapURL[string(body)] = hash[:6]
		//fmt.Print(mapURL)
		w.WriteHeader(http.StatusCreated)

	}
}

func idPage(w http.ResponseWriter, r *http.Request) {
	//TODO: handler "id"
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainPage)
	mux.HandleFunc("/{id}", idPage)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
