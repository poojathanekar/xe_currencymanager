package model

import (
	"time"
)

// Response data struct
type XEResponse struct {
	Terms     string    `json:"terms"`
	Privacy   string    `json:"privacy"`
	From      string    `json:"from"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
	To        []To      `json:"to"`
}

// To QuoteCurrency model
type To struct {
	Quotecurrency string  `json:"quotecurrency"`
	Mid           float64 `json:"mid"`
}
