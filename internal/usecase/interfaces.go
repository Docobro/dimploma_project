package usecase

import (
	"time"

	"github.com/docobro/dimploma_project/internal/entity"
)

type ClickhouseRepo interface {
	// example BTC, time.Hour * 1 (priceIndex1H)
	CreatePrices(prices []entity.Coin) error
	GetPrices(coins []string, start time.Time, end time.Time) interface{}
	CreateIndices(indices []entity.Indices) error
	CalculatePriceIndex(coin string, timeAgo time.Duration) float64
	CreateVolumes1m(volume []entity.VolumeTo) error
	CreateSupplies(supplies []entity.Supplies) error
	CreateVolatility(volatility map[string]float64) error
	PearsonPriceToVolumeCorrelation(coin string) float64
	PearsonPriceToMrktCapCorrelation(coin string) float64
	PearsonPriceToVolatilityCorrelation(coin string) float64
	ReturnPredictions(coin string) ([]float64, float64)
	CreatePearson(coeff []entity.PearsonPriceTo) error
	UpdatePredict(pred []entity.Predictions) error
}

type CryptoRepo interface {
	GetCurrencies(coins []string) (map[string]*entity.Coin, error)
	GetCryptoFullInfo(coins []string, currencies []string) (map[string]entity.Coin, error)
	GetOneMinuteData(coin string, currency string, limit int) ([]entity.Coin, error)
	ReturnVolatility(coins []string) map[string]float64
}
