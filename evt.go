package main

const (
	ActionAdd    = "add"
	ActionRemove = "remove"
)

const (
	DefaultRad = 0.1
)

type Evt struct {
	Action string  `json:"action"`
	ID     string  `json:"id"`
	Lng    float64 `json:"lng"`
	Lat    float64 `json:"lat"`
	Rad    float64 `json:"rad"`
}
