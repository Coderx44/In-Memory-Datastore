package main

import (
	"github.com/Coderx44/gg/server"
)

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"sync"
// 	"time"

// 	"github.com/gorilla/mux"
// )

// type datastore struct {
// 	sync.RWMutex
// 	data map[string]*dataItem
// }

// type dataItem struct {
// 	value      string
// 	expiryTime time.Time
// }

// func newDatastore() *datastore {
// 	return &datastore{
// 		data: make(map[string]*dataItem),
// 	}
// }

// func (ds *datastore) set(key string, value string, expiryTime time.Time, condition string) error {
// 	ds.Lock()
// 	defer ds.Unlock()

// 	_, ok := ds.data[key]
// 	if condition == "NX" && ok {
// 		return fmt.Errorf("key already exists")
// 	}
// 	if condition == "XX" && !ok {
// 		return fmt.Errorf("key does not exist")
// 	}

// 	ds.data[key] = &dataItem{
// 		value:      value,
// 		expiryTime: expiryTime,
// 	}

// 	return nil
// }

// func (ds *datastore) get(key string) (string, bool) {
// 	ds.RLock()
// 	defer ds.RUnlock()

// 	item, ok := ds.data[key]
// 	if !ok {
// 		return "", false
// 	}

// 	if !item.expiryTime.IsZero() && item.expiryTime.Before(time.Now()) {
// 		// If the item has expired, delete it and return false
// 		delete(ds.data, key)
// 		return "", false
// 	}

// 	return item.value, true
// }

// func main() {
// 	ds := newDatastore()

// 	r := mux.NewRouter()
// 	r.HandleFunc("/set/{key}", func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		key := vars["key"]

// 		value := r.FormValue("value")
// 		if value == "" {
// 			http.Error(w, "missing value parameter", http.StatusBadRequest)
// 			return
// 		}

// 		var expiryTime time.Time
// 		expiryStr := r.FormValue("expiry")
// 		if expiryStr != "" {
// 			d, err := time.ParseDuration(expiryStr)
// 			if err != nil {
// 				http.Error(w, "invalid expiry parameter", http.StatusBadRequest)
// 				return
// 			}
// 			expiryTime = time.Now().Add(d)
// 		}

// 		condition := r.FormValue("condition")

// 		if err := ds.set(key, value, expiryTime, condition); err != nil {
// 			http.Error(w, err.Error(), http.StatusConflict)
// 			return
// 		}

// 		w.WriteHeader(http.StatusOK)
// 	}).Methods("POST")

// 	r.HandleFunc("/get/{key}", func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		key := vars["key"]

// 		value, ok := ds.get(key)
// 		if !ok {
// 			http.NotFound(w, r)
// 			return
// 		}

// 		json.NewEncoder(w).Encode(map[string]string{
// 			"value": value,
// 		})
// 	}).Methods("GET")

// 	http.ListenAndServe(":8080", r)
// }

const PORT = ":3000"

func main() {
	server.StartApiServer()

}
