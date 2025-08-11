package services

import (
	"context"
	"crypto-api/config"
	"crypto-api/db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type cmcResponse struct {
	Data []cmcCoin `json:"data"`
}

type cmcCoin struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Symbol      string            `json:"symbol"`
	LastUpdated string            `json:"last_updated"`
	Quote       map[string]cmcUSD `json:"quote"`
}

type cmcUSD struct {
	Price            float64 `json:"price"`
	MarketCap        float64 `json:"market_cap"`
	Volume24h        float64 `json:"volume_24h"`
	PercentChange24h float64 `json:"percent_change_24h"`
}

func FetchAndStoreCryptoData(cfg config.Config) {
	url := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?limit=100&convert=USD"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Failed to create request:", err)
	}

	req.Header.Add("X-CMC_PRO_API_KEY", cfg.CMCAPIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to fetch data:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Non-200 response: %d\nBody: %s", resp.StatusCode, body)
	}

	var result cmcResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal("Failed to decode JSON:", err)
	}

	ctx := context.Background()
	for _, coin := range result.Data {
		usd, ok := coin.Quote["USD"]
		if !ok {
			log.Printf("Skipping %s: no USD quote", coin.Symbol)
			continue
		}

		updated, err := time.Parse(time.RFC3339, coin.LastUpdated)
		if err != nil {
			updated = time.Now() // fallback
		}

		query := `
            INSERT INTO cryptocurrencies (coin_id, name, symbol, price_usd, market_cap_usd, volume_24h, percent_change_24h, last_updated)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
            ON CONFLICT (coin_id) DO UPDATE
            SET price_usd=$4, market_cap_usd=$5, volume_24h=$6, percent_change_24h=$7, last_updated=$8;
        `
		_, err = db.DB.Exec(ctx, query,
			coin.ID,
			coin.Name,
			coin.Symbol,
			usd.Price,
			int64(usd.MarketCap),
			int64(usd.Volume24h),
			usd.PercentChange24h,
			updated,
		)
		if err != nil {
			log.Printf("Failed to insert/update %s: %v", coin.Symbol, err)
		}
	}

	fmt.Println("Crypto data updated")
}
