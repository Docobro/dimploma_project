package usecase

import (
	"math/rand"

	"github.com/docobro/dimploma_project/internal/entity"
)

// load your adapters
type Usecase struct {
	storage    ClickhouseRepo
	cryptoRepo CryptoRepo
}

func (uc *Usecase) GetCurrencies(coins []string) (map[string]*entity.Coin, error) {
	return uc.cryptoRepo.GetCurrencies(coins)
}

func (uc *Usecase) GetPriceIndices(coins []string) (map[string]float64, error) {
	res := make(map[string]float64, len(coins))
	for _, v := range coins {
		res[v] = rand.Float64()
	}
	return res, nil
}

func (uc *Usecase) GetVolumeIndices(coins []string) (map[string]float64, error) {
	res := make(map[string]float64, len(coins))
	for _, v := range coins {
		res[v] = rand.Float64()
	}
	return res, nil
}

func New(storage ClickhouseRepo) *Usecase {
	return &Usecase{
		storage: storage,
	}
}
