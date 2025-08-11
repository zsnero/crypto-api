package handlers

import (
	"context"
	"crypto-api/db"
	"crypto-api/models"
	"encoding/json"
	"net/http"
)

func GetCryptocurrencies(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	rows, err := db.DB.Query(ctx, `
        SELECT coin_id, name, symbol, price_usd, market_cap_usd, volume_24h, percent_change_24h, last_updated
        FROM cryptocurrencies;
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var coins []models.Cryptocurrency
	for rows.Next() {
		var c models.Cryptocurrency
		err := rows.Scan(&c.CoinID, &c.Name, &c.Symbol, &c.PriceUSD, &c.MarketCapUSD, &c.Volume24h, &c.PercentChange24, &c.LastUpdated)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		coins = append(coins, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coins)
}
