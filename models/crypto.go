package models

import "time"

type Cryptocurrency struct {
	CoinID          int       `json:"coin_id"`
	Name            string    `json:"name"`
	Symbol          string    `json:"symbol"`
	PriceUSD        float64   `json:"price_usd"`
	MarketCapUSD    int64     `json:"market_cap_usd"`
	Volume24h       int64     `json:"volume_24h"`
	PercentChange24 float64   `json:"percent_change_24h"`
	LastUpdated     time.Time `json:"last_updated"`
}
