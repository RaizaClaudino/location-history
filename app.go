package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type App struct {
	Router    *mux.Router
	Addr      string
	Locations map[string][]Location
}

func (a *App) Initialize(addr string) {
	router := mux.NewRouter()
	locations := map[string][]Location{}

	router.HandleFunc("/location/{order_id}", a.CreateLocation).Methods("PUT")
	router.HandleFunc("/location/{order_id}", a.GetLocationHistory).Queries("max", "{[0-9]*?}").Methods("GET")
	router.HandleFunc("/location/{order_id}", a.DeleteLocationHistory).Methods("DELETE")

	if addr == "" {
		addr = ":8080"
	}

	a.Addr = addr
	a.Locations = locations
	a.Router = router
}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(a.Addr, a.Router))
}

func (a *App) GetLocationHistory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	index := params["order_id"]
	history := a.Locations[index]
	max := len(history)
	maxparam := params["max"]

	num, err := strconv.Atoi(maxparam)
	if err == nil && max > num {
		max = num
	}

	response := LocationResponse{OrderID: index, History: history[:max]}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (a *App) CreateLocation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var location Location
	err := json.NewDecoder(r.Body).Decode(&location)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	index := params["order_id"]
	history := a.Locations[index]
	history = append(history, location)
	a.Locations[index] = history

	w.WriteHeader(http.StatusOK)
}

func (a *App) DeleteLocationHistory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	index := params["order_id"]
	a.Locations[index] = []Location{}
}
