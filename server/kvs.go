package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var ErrorNoSuchKey = errors.New("no such key")

var store = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

func Put(key string, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()

	return nil
}

func Get(key string) (string, error) {
	store.RLock()
	value, ok := store.m[key]
	store.RUnlock()

	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func Delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.Unlock()

	return nil
}

func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	key := vars["key"]
	value, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	err = Put(key, string(value))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	//transact.WritePut(key, string(value))
	w.WriteHeader(http.StatusCreated)
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	key := vars["key"]

	value, err := Get(key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Write([]byte(value))
}

func keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	key := vars["key"]

	err := Delete(key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	//transact.WriteDelete(key)
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", keyValueDeleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
