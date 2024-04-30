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

func (r *Repository) GetCryptoFullInfo(coins []string, currencies []string) interface{} {
	res, err := r.client.GetPriceMultiFull(coins, currencies)
	if err != nil {
		return err
	}
	return res
}
