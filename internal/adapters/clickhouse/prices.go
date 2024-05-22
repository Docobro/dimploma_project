package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/docobro/dimploma_project/internal/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// добавление текущих цен и маркет капов
func (r *Repository) CreatePrices(prices []entity.Coin) error {
	tokens, err := r.GetCryptoTokens()
	if err != nil {
		return err
	}

	batch, err := r.conn.PrepareBatch(context.Background(), "INSERT INTO prices")
	if err != nil {
		return err
	}

	for _, v := range prices {
		for _, p := range v.Prices {
			err := batch.Append(
				uuid.New(), decimal.NewFromFloat(p), time.Now(), time.Now(), tokens[v.Name].ID, v.MarketCap)
			if err != nil {
				return err
			}
		}
	}

	err = batch.Send()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CalculatePriceIndex(coin string, timeAgo time.Duration) float64 {
	tokens, err := r.GetCryptoTokens()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return 0
	}
	id := tokens[coin].ID
	indexQuery := `SELECT 
((latest_price - lag_price)/lag_price)*100 AS price_index
FROM 
(SELECT 
    (
        SELECT 
            CAST(value AS Float64)
        FROM 
            cryptowallet.prices
        WHERE 
            crypto_id = $1
        ORDER BY 
            created_at DESC
        LIMIT 1
    ) AS latest_price,
    (
        SELECT 
            CAST(value AS Float64)
        FROM 
            cryptowallet.prices
        WHERE 
            crypto_id = $1 AND toStartOfHour(created_at) = toStartOfHour(NOW() - INTERVAL $2 HOUR)      
        LIMIT 1 
    ) AS lag_price,
    (
        SELECT 
            toStartOfHour(created_at)
        FROM 
            cryptowallet.prices
        WHERE 
            crypto_id = $1
        ORDER BY 
            created_at DESC
        LIMIT 1
    ) AS current_time
FROM 
    system.one)
LIMIT 1;
`
	var res float64
	err = r.conn.QueryRow(context.Background(), indexQuery, id, timeAgo.Hours()).Scan(&res)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return res
}
