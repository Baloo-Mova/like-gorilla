package models

import "time"

type PriceInfo struct {
	TimeStamp string `json:"timestamp"`
	Symbol string `json:"symbol"`
	Price float64 `json:"price"`
}

type BitmexRequestModel struct {
	Op string `json:"op"`
	Args []string `json:"args"`
}

type BitmexResponseModel struct {
	Table  string   `json:"table"`
	Action string   `json:"action"`
	Data []struct {
		Symbol                         string      `json:"symbol,omitempty"`
		LastPrice                      float64     `json:"lastPrice,omitempty"`
		Timestamp                      time.Time   `json:"timestamp,omitempty"`
	} `json:"data"`
}

type SubscriptionStatusRequest struct {
	Action string `json:"action"`
	Symbols []string `json:"symbols"`
}