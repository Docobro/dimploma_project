package usecase

import "github.com/docobro/dimploma_project/internal/entity"

type ClickhouseRepo interface{}

type CryptoRepo interface {
	GetCurrencies(coins []string) (map[string]*entity.Coin, error)
}
type Cr interface {
	GetCurrencies(coins []string) (map[string]*entity.Coin, error)
	GetPriceIndices(coins []string) (map[string]float64, error)
	GetVolumeIndices(coins []string) (map[string]float64, error)
}
