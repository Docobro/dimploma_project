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
	log.Println("parse prices done")
	return uc.storage.CreatePrices(pricesReq)
}

func (uc *Usecase) CalculatePriceIndex(coin string, timeAgo time.Duration) float64 {
	return uc.storage.CalculatePriceIndex(coin, timeAgo)
}

func (uc *Usecase) ParsePearson() error {
	log.Println("do parse pearson")
	coins := []string{"BTC", "ETH"}
	pearsonReq := []entity.PearsonPriceVolMrkt{}
	for i := 0; i < len(coins); i++ {
		coeffPearsonVolume := uc.storage.PearsonPriceToVolumeCorrelation(coins[i])
		coeffPersonMrktCap := uc.storage.PearsonPriceToMrktCapCorrelation(coins[i])
		coef := entity.PearsonPriceVolMrkt{
			CryptoName: coins[i],
			Volume:     coeffPearsonVolume,
			MrktCap:    coeffPersonMrktCap,
		}
		pearsonReq = append(pearsonReq, coef)
	}
	err := uc.storage.CreatePearson(pearsonReq)
	if err != nil {
		return err
	}
	log.Println("parse pearson done")
	return nil
}

func (uc *Usecase) CreateIndices() error {
	log.Println("do parse indices")
	coins := []string{"BTC", "ETH"}

	indicesReq := []entity.Indices{}
	for i := 0; i < len(coins); i++ {
		priceIndex := uc.CalculatePriceIndex(coins[i], time.Hour*1)
		indices := entity.Indices{
			CryptoName: coins[i],
			Price:      entity.PriceIndex{Value: priceIndex},
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
	volumes := make(map[string]entity.Coin, len(coins))
	for i := 0; i < len(coins); i++ {
		volume, err := uc.cryptoRepo.GetOneMinuteData(coins[i], "USD", 1)
		if err != nil {
			return err
		}
		volumes[coins[i]] = volume[coins[i]]
	}

	volumeReq := []entity.VolumeTo{}
	for i := 0; i < len(coins); i++ {
		volume := entity.VolumeTo{
			CryptoName: coins[i],
			Value:      volumes[coins[i]].VolumeTo,
		}
		volumeReq = append(volumeReq, volume)
	}
	err := uc.storage.CreateVolumes1m(volumeReq)
	if err != nil {
		return err
	}
	log.Println("parse volume 1 minute done")
	return nil
}

func (uc *Usecase) ParseVolatility() error {
	log.Println("do parse volatility minute")
	coins := []string{"BTC", "ETH"}
	volatility := uc.cryptoRepo.ReturnVolatility(coins)

	err := uc.storage.CreateVolatility(volatility)
	if err != nil {
		return err
	}
	log.Println("parse volatility done")
	return nil
}
