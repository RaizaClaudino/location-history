package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize(os.Getenv("HISTORY_SERVER_LISTEN_ADDR"))
	code := m.Run()
	os.Exit(code)
}

func TestGetLocationHistory(t *testing.T) {
	a.Initialize("")
	location := Location{
		Lat: 11.2,
		Lng: 12.2,
	}
	a.Locations["abc123"] = []Location{location}
	go a.Run()

	res, err := http.Get("http://localhost:8080/location/abc123?max=10")
	if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
	}

	bodyBytes, _ := ioutil.ReadAll(res.Body)
	var resp LocationResponse
	err = json.Unmarshal(bodyBytes, &resp)

	if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
	}

	if resp.History[0].Lat != location.Lat || resp.History[0].Lng != location.Lng {
		t.Errorf("Expected to see abc123 location history, received %s", string(bodyBytes))
	}
}

func TestCreateLocation(t *testing.T) {
	a.Initialize("")
	go a.Run()

	body := &Location{
		Lat: 11.2,
		Lng: 12.2,
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	req, _ := http.NewRequest("PUT", "http://localhost:8080/location/abc123", payloadBuf)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
	}

	history := a.Locations["abc123"]
	if len(history) == 0 {
		t.Errorf("Expected locations history to have a new location, but got none")
	}
}
