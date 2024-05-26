package cryptopackage

import (
	"errors"
	"fmt"
	"log"

	"github.com/docobro/dimploma_project/internal/adapters/cryptopackage/client"
	"github.com/docobro/dimploma_project/internal/config"
	"github.com/docobro/dimploma_project/internal/entity"
)

type Repository struct {
	///client
	client *client.Client
}

func New(cfg config.CryptoConfig) *Repository {
	r := &Repository{
		client: client.New(cfg.Url, cfg.Key),
	}

	// prefligh ping
	_, err := r.client.GetPrice([]string{"BTC"}, []string{"USD"})
	if err != nil {
		log.Fatalf("failed to init crypto package err:%v", err)
	}

	return r
}

// return coins where string is coin_name
func (r *Repository) GetCurrencies(coins []string) (map[string]*entity.Coin, error) {
	currencies := make(map[string]*entity.Coin, len(coins))
	if len(coins) == 0 {
		return currencies, errors.New("coins lengt must be at least 1")
	}

	// get currencies
	res, err := r.client.GetCurrentPrices(coins)
	if err != nil {
		return nil, errors.New("adapters - cryptopackage - GetCurrencies err:%v" + err.Error())
	}
	// map over and parse into entity value
	for _, coin := range coins {
		currencies[coin] = nil
		if _, ok := res[coin]; ok {
			currencies[coin] = &entity.Coin{
				Name:      coin,
				Prices:    res[coin].Prices,
				MarketCap: 30487,
			}
		}
	}
	return currencies, nil
}

func (r *Repository) GetCryptoFullInfo(coins []string, currencies []string) (map[string]entity.Coin, error) {
	res, err := r.client.GetPriceMultiFull(coins, currencies)
	if err != nil {
		return nil, err
	}
	coinsRes := make(map[string]entity.Coin)
	for k, coinName := range res.RAW {
		coinsRes[string(k)] = entity.Coin{
			Name:         string(k),
			MarketCap:    res.RAW[client.CoinType(k)]["USD"].MktCap,
			VolumeHour:   res.RAW[client.CoinType(k)]["USD"].VolumeHour,
			Volume24Hour: res.RAW[client.CoinType(k)]["USD"].Volume24Hour,
			Prices:       map[string]float64{},
			Supply:       res.RAW[client.CoinType(k)]["USD"].Supply,
		}

		for currency_key, currency := range coinName {
			coinsRes[string(k)].Prices[string(currency_key)] = currency.Price
		}
	}
	return coinsRes, nil
}

func (r *Repository) GetOneMinuteData(coin string, currency string, limit int) ([]entity.Coin, error) {
	res, err := r.client.GetOneMinuteFull(coin, currency, limit)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}
	minuteRes := make([]entity.Coin, limit)
	for i := 0; i < limit; i++ {
		minuteRes[i] = entity.Coin{
			Name:        coin,
			VolumeTo:    res.Data.Data[i].VolumeTo,
			CloseMinute: res.Data.Data[i].Close,
			OpenMinute:  res.Data.Data[i].Open,
			HighMinute:  res.Data.Data[i].High,
			LowMinute:   res.Data.Data[i].Low,
		}
	}
	fmt.Printf("minuteRes: %v\n", minuteRes)
	return minuteRes, nil
}
