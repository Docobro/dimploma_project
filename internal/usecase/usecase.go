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

func (uc *Usecase) GetPriceIndices(coins []string) (map[string]int32, error) {
	res := make(map[string]int32, len(coins))
	for _, v := range coins {
		res[v] = rand.Int31()
	}
	return res, nil
}

func (uc *Usecase) GetVolumeIndices(coins []string) (map[string]float32, error) {
	res := make(map[string]float32, len(coins))
	for _, v := range coins {
		res[v] = rand.Float32()
	}
	return res, nil
}

func New(storage ClickhouseRepo, cryptoRepo CryptoRepo) *Usecase {
	return &Usecase{
		storage:    storage,
		cryptoRepo: cryptoRepo,
	}
}
