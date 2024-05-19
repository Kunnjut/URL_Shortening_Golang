package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var mapURL = map[string]string{}

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
	} else {
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
		//TODO: Сделать ответ "Ссылка запроса - полный ответ с localhost" (мб)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(len(mapURL[string(body)])))
		fmt.Fprintln(w, "http://localhost:8080/"+mapURL[string(body)])
		// curl -i -X POST -H "Content-Type: text/plain" -d "Hello, W23or" http://localhost:8080
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
