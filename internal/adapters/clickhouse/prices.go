package clickhouse

import (
	"context"
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
				uuid.New(), decimal.NewFromFloat(p), v.MarketCap, time.Now(), time.Now(), tokens[v.Name].ID)
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
