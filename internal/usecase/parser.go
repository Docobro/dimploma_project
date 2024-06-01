package usecase

import (
	"fmt"
	"log"
	"time"

	"github.com/docobro/dimploma_project/internal/entity"
)

func (uc *Usecase) CoinList() []string {
	coins := []string{"ETH", "BTC"}
	return coins
}

func (uc *Usecase) ParsePrices() error {
	log.Println("do parse prices")
	coins := uc.CoinList()
	currencies := []string{"USD"}
	pricesReq := []entity.Coin{}

	for i := 0; i < len(coins); i++ {
		cryptoFullInfo, err := uc.cryptoRepo.GetCryptoFullInfo(coins, currencies)
		if err != nil {
			return err
		}
		res := entity.Coin{
			Name:      coins[i],
			Prices:    cryptoFullInfo[coins[i]].Prices,
			MarketCap: cryptoFullInfo[coins[i]].MarketCap,
		}
		pricesReq = append(pricesReq, res)
	}

	log.Println("parse prices done")
	return uc.storage.CreatePrices(pricesReq)
}

func (uc *Usecase) CalculatePriceIndex(coin string, timeAgo time.Duration) float64 {
	return uc.storage.CalculatePriceIndex(coin, timeAgo)
}

func (uc *Usecase) ParsePearson() error {
	log.Println("do parse pearson")
	coins := uc.CoinList()
	pearsonReq := []entity.PearsonPriceTo{}
	for i := 0; i < len(coins); i++ {
		coeffPearsonVolume := uc.storage.PearsonPriceToVolumeCorrelation(coins[i])
		coeffPersonMrktCap := uc.storage.PearsonPriceToMrktCapCorrelation(coins[i])
		coeffPersonVolatility := uc.storage.PearsonPriceToVolatilityCorrelation(coins[i])
		coef := entity.PearsonPriceTo{
			CryptoName: coins[i],
			Volume:     coeffPearsonVolume,
			MrktCap:    coeffPersonMrktCap,
			Volatility: coeffPersonVolatility,
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
	coins := uc.CoinList()

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
	coins := uc.CoinList()
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
	coins := uc.CoinList()
	volumes := make(map[string][]entity.Coin, len(coins))
	for i := 0; i < len(coins); i++ {
		volumeData, err := uc.cryptoRepo.GetOneMinuteData(coins[i], "USD", 1)
		if err != nil {
			return err
		}
		volumes[coins[i]] = volumeData
	}

	volumeReq := []entity.VolumeTo{}
	for i := 0; i < len(coins); i++ {
		if len(volumes[coins[i]]) == 0 {
			return fmt.Errorf("no volume data available for coin: %s", coins[i])
		}
		firstVolume := volumes[coins[i]][0] // получаем первый элемент слайса
		volume := entity.VolumeTo{
			CryptoName: coins[i],
			Value:      firstVolume.VolumeTo,
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
	coins := uc.CoinList()
	volatility := uc.cryptoRepo.ReturnVolatility(coins)

	err := uc.storage.CreateVolatility(volatility)
	if err != nil {
		return err
	}
	log.Println("parse volatility done")
	return nil
}

func (uc *Usecase) ParsePredictions() error {
	log.Println("do parse predictions")
	coins := uc.CoinList()
	res := []entity.Predictions{}
	for i := 0; i < len(coins); i++ {
		predict, mse := uc.storage.ReturnPredictions(coins[i])
		pred := entity.Predictions{
			CryptoName: coins[i],
			Value:      predict,
			Mse:        mse,
		}
		res = append(res, pred)
	}

	err := uc.storage.UpdatePredict(res)
	if err != nil {
		return err
	}
	log.Println("parse predictions done")
	return nil
}
