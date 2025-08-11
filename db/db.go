package db

import (
	"context"
	"crypto-api/config"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect(cfg config.Config) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	var err error
	DB, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	createTable()
}

func createTable() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
    CREATE TABLE IF NOT EXISTS cryptocurrencies (
        id SERIAL PRIMARY KEY,
        coin_id INT UNIQUE,
        name TEXT,
        symbol TEXT,
        price_usd DECIMAL,
        market_cap_usd BIGINT,
        volume_24h BIGINT,
        percent_change_24h DECIMAL,
        last_updated TIMESTAMP
    );
`

	_, err := DB.Exec(ctx, query)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}
