package model

import (
	"time"

	httperors "github.com/myrachanto/custom-http-error"
)

type Order struct {
	OrderId  string     `json:"order_id"`
	Location []Location `json:"history"`
}
type Location struct {
	Lat     float64   `json:"lat"`
	Lng     float64   `json:"lng"`
	Thetime time.Time `json:"thetime"`
}

//Validate ...
func (l Location) Validate() httperors.HttpErr {
	if l.Lat == 0 {
		return httperors.NewNotFoundError("Invalid lat")
	}
	if l.Lng == 0 {
		return httperors.NewNotFoundError("Invalid lng")
	}
	return nil
}
