package usecase

import (
	"log"
	"time"

	"github.com/docobro/dimploma_project/internal/entity"
)

func (uc *Usecase) ParsePrices() error {
	log.Println("do parse prices")
	coins := []string{"ETH", "BTC", "TON"}
	currencies := []string{"USD"}
	cryptoFullInfo, err := uc.cryptoRepo.GetCryptoFullInfo(coins, currencies)
	if err != nil {
		return err
	}
	pricesReq := []entity.Coin{}
	for _, v := range cryptoFullInfo {
		pricesReq = append(pricesReq, v)
	}
	return uc.storage.CreatePrices(pricesReq)
}

func (uc *Usecase) CalculatePriceIndex(coin string, timeAgo time.Duration) float64 {
	return uc.storage.CalculatePriceIndex(coin, timeAgo)
}

func (uc *Usecase) GetVolumeIndex(coins []string, start time.Time, end time.Time) float32 {
	return 0
}

func (uc *Usecase) CreateIndices() error {
	log.Println("do parse prices")
	coins := []string{"BTC"}
	currencies := []string{"USD"}
	cryptoFullInfo, err := uc.cryptoRepo.GetCryptoFullInfo(coins, currencies)
	if err != nil {
		return err
	}
	priceIndex := uc.CalculatePriceIndex(coins[0], time.Hour*1)
	indices := entity.Indices{
		CryptoName: coins[0],
		Price:      entity.PriceIndex{Value: priceIndex},
		Volume:     entity.VolumeIndex{Value: float32(cryptoFullInfo[coins[0]].VolumeHour)},
	}
	err = uc.storage.CreateIndices([]entity.Indices{indices})
	if err != nil {
		return err
	}
	log.Println("Indices created!")
	return nil
}

func (uc *Usecase) GetPrices(coins []string, start time.Time, end time.Time) {
	res := uc.storage.GetPrices(coins, start, end)
	log.Println(res)
}
