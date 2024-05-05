package cryptopackage

import (
	"errors"

	"github.com/docobro/dimploma_project/internal/adapters/cryptopackage/client"
	"github.com/docobro/dimploma_project/internal/config"
	"github.com/docobro/dimploma_project/internal/entity"
)

type Repository struct {
	///client
	client *client.Client
}

func New(cfg config.CryptoConfig) *Repository {
	return &Repository{
		client: client.New(cfg.Url, cfg.Key),
	}
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
		}

		for currency_key, currency := range coinName {
			coinsRes[string(k)].Prices[string(currency_key)] = currency.Price
		}
	}
	return coinsRes, nil
}

func (r *Repository) GetCurrencyTransactionCount(coins []string) (map[string]uint32, error) {
	raw, err := r.client.GetTransactionData(coins)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
