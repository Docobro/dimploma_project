package usecase

import (
	"time"

	"github.com/docobro/dimploma_project/internal/entity"
)

type ClickhouseRepo interface {
	CreatePrices(prices []entity.Coin) error
	GetPrices(coins []string, start time.Time, end time.Time) interface{}
	CreateIndices(indices []entity.Indices) error
	// example BTC, time.Hour * 1 (priceIndex1H)
	// example BTC, time.Hour * 24 (priceIndex24H)
	CalculatePriceIndex(coin string, timeAgo time.Duration) float64
	CalculateVolumeIndex(coin string, timeAgo time.Duration) float64
	CreateTransaction(transaction map[string]uint32) error
	CreateVolumes1m(volume map[string]float32) error
	CreateSupplies(supplies []entity.Supplies) error
}

type CryptoRepo interface {
	GetCurrencies(coins []string) (map[string]*entity.Coin, error)
	GetCryptoFullInfo(coins []string, currencies []string) (map[string]entity.Coin, error)
	GetCurrencyTransactionCount(coins []string) (map[string]uint32, error)
	GetOneMinuteVolume(coins []string) (map[string]float32, error)
}
