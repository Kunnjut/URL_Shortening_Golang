package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

var mapURL = map[string]string{}

func searchKey(val string) string {
	result := ""
	for key, value := range mapURL {
		if value == val {
			result = key
			break
		}
	}
	return result
}

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
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Content-Length", strconv.Itoa(len(mapURL[string(body)])))
		w.WriteHeader(http.StatusCreated)
		_, err = fmt.Fprintln(w, "https://"+r.Host+"/"+mapURL[string(body)])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// curl -i -X POST -H "Content-Type: text/plain" -d "https://dzen.ru/" http://localhost:8080
	}
}

func idPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
	}
	vars := mux.Vars(r)
	id := vars["id"] // TODO: Значение мы получили здесь, вытаскиваем ключ из мапы по этому ключу и делаем редирект
	key := searchKey(id)
	if key == "" {
		http.NotFound(w, r)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Location", key)
	w.WriteHeader(http.StatusTemporaryRedirect)
	//curl -i -X  GET http://localhost:8080/8d26b0
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", mainPage)
	router.HandleFunc("/{id}", idPage).Methods(http.MethodGet)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
