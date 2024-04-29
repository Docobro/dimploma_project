package usecase

import "github.com/docobro/dimploma_project/internal/entity"

type ClickhouseRepo interface {
	CreatePrices(prices []entity.Coin) error
}

type CryptoRepo interface {
	GetCurrencies(coins []string) (map[string]*entity.Coin, error)
}
