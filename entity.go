package main

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type LocationResponse struct {
	OrderID string     `json:"order_id"`
	History []Location `json:"history"`
}
