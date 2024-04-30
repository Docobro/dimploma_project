package usecase

import (
	"time"

	"github.com/docobro/dimploma_project/internal/entity"
)

type ClickhouseRepo interface {
	CreatePrices(prices []entity.Coin) error
	GetPrices(coins []string, start time.Time, end time.Time) interface{}
}

type CryptoRepo interface {
	GetCurrencies(coins []string) (map[string]*entity.Coin, error)
	GetPricesCoock() interface{}
}
