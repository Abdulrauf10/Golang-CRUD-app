package models

type Harga struct {
	ID        int64 `json:"id"`
	Liter     int64 `json:"liter"`
	Premium   int64 `json:"premium"`
	Pertalite int64 `json:"pertalite"`
}