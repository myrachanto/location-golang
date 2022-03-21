package model

import "time"

type Order struct {
	OrderId  string     `json:"order_id"`
	Location []Location `json:"history"`
}
type Location struct {
	Lat     float64   `json:"lat"`
	Lng     float64   `json:"lng"`
	Thetime time.Time `json:"thetime"`
}
