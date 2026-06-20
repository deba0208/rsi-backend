package models

import "time"

type Candle struct {
	Date  time.Time `json:"date"`
	Close float64   `json:"close"`
}