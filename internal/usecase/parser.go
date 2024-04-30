package usecase

import (
	"log"
	"time"

	"github.com/docobro/dimploma_project/internal/entity"
)

func (uc *Usecase) ParsePrices() error {
	log.Println("do parse prices")
	prices := []string{"ETH", "BTC", "TON"}
	currencies, err := uc.GetCurrencies(prices)
	if err != nil {
		return err
	}
	pricesReq := []entity.Coin{}
	for _, v := range currencies {
		pricesReq = append(pricesReq, *v)
	}
	return uc.storage.CreatePrices(pricesReq)
}

func (uc *Usecase) GetPriceIndex(coins []string, period time.Duration) int32 {
	return 0
}

func (uc *Usecase) GetVolumeIndex(coins []string, period time.Duration) int32 {
	return 0
}

func (uc *Usecase) CreateIndices(entity.Indices) error {
	return nil
}

func (uc *Usecase) GetPrices(coins []string, start time.Time, end time.Time) {
	res := uc.storage.GetPrices(coins, start, end)
	log.Println(res)
}

func (uc *Usecase) Aboba() {
	res := uc.cryptoRepo.GetPricesCoock()
	log.Println(res)
}
