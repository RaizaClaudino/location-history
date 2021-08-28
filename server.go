package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type LocationResponse struct {
	OrderID string     `json:"order_id"`
	History []Location `json:"history"`
}

var locations map[string][]Location

func GetLocationHistory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	index := params["order_id"]
	history := locations[index]
	max := len(history)

	if params["max"] != "" {
		num, err := strconv.Atoi(params["max"])
		if err == nil && max > num {
			max = num
		}
	}
	response := LocationResponse{OrderID: index, History: history[:max]}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateLocation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var location Location
	err := json.NewDecoder(r.Body).Decode(&location)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	index := params["order_id"]
	history := locations[index]
	history = append(history, location)
	locations[index] = history

	w.WriteHeader(http.StatusOK)
}

func DeleteLocationHistory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	index := params["order_id"]
	locations[index] = []Location{}
}

func main() {
	router := mux.NewRouter()
	locations = map[string][]Location{}

	router.HandleFunc("/location/{order_id}", CreateLocation).Methods("PUT")
	router.HandleFunc("/location/{order_id}", GetLocationHistory).Queries("max", "{[0-9]*?}").Methods("GET")
	router.HandleFunc("/location/{order_id}", DeleteLocationHistory).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
