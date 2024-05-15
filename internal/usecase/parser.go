package usecase

import (
	"log"
	"time"

	"github.com/docobro/dimploma_project/internal/entity"
)

func (uc *Usecase) ParsePrices() error {
	log.Println("do parse prices")
	coins := []string{"ETH", "BTC"}
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
func (uc *Usecase) CalculateVolumeIndex(coin string, timeAgo time.Duration) float64 {
	return uc.storage.CalculateVolumeIndex(coin, timeAgo)
}

func (uc *Usecase) CreateIndices() error {
	log.Println("do parse indices")
	coins := []string{"BTC", "ETH"}

	indicesReq := []entity.Indices{}
	for i := 0; i < len(coins); i++ {
		priceIndex := uc.CalculatePriceIndex(coins[i], time.Hour*1)
		volumeIndex := uc.CalculateVolumeIndex(coins[i], time.Hour*1)
		indices := entity.Indices{
			CryptoName: coins[i],
			Price:      entity.PriceIndex{Value: priceIndex},
			Volume:     entity.VolumeIndex{Value: volumeIndex},
		}
		indicesReq = append(indicesReq, indices)
	}
	err := uc.storage.CreateIndices(indicesReq)
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

func (uc *Usecase) ParseTrasactionCount() error {
	if time.Now().UTC().Hour() != 0 {
		return nil
	}

	log.Println("do parse transaction")
	coins := []string{"BTC", "ETH"}

	transaction, err := uc.cryptoRepo.GetCurrencyTransactionCount(coins)
	if err != nil {
		return err
	}
	err = uc.storage.CreateTransaction(transaction)
	if err != nil {
		return err
	}
	log.Println("parse transaction done")
	return nil
}

func (uc *Usecase) ParseTotalSupply() error {
	log.Println("do parse total supply")
	coins := []string{"BTC", "ETH"}
	currencies := []string{"USD"}

	cryptoFullInfo, err := uc.cryptoRepo.GetCryptoFullInfo(coins, currencies)
	if err != nil {
		return err
	}
	supplyReq := []entity.Supplies{}
	for i := 0; i < len(coins); i++ {
		supply := entity.Supplies{
			CryptoName: coins[i],
			Value:      cryptoFullInfo[coins[i]].Supply,
		}
		supplyReq = append(supplyReq, supply)
	}
	err = uc.storage.CreateSupplies(supplyReq)
	if err != nil {
		return err
	}
	log.Println("parse total supply done")
	return nil
}

func (uc *Usecase) ParseVolumeMinute() error {
	log.Println("do parse volume 1 minute")
	coins := []string{"BTC", "ETH"}

	volume, err := uc.cryptoRepo.GetOneMinuteVolume(coins)
	if err != nil {
		return err
	}
	err = uc.storage.CreateVolumes1m(volume)
	if err != nil {
		return err
	}
	log.Println("parse volume 1 minute done")
	return nil
}
