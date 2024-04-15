package cryptopackage

import (
	"errors"

	"github.com/docobro/dimploma_project/internal/adapters/cryptopackage/client"
	"github.com/docobro/dimploma_project/internal/entity"
)

type Repository struct {
	///client
	client *client.Client
}

func New(url string) *Repository {
	return &Repository{
		client: client.New(url),
	}
}

func (r *Repository) GetCurrencies(coins []string) (map[string]*entity.Currency, error) {
	res, err := r.client.GetCurrentPrices(coins)
	if err != nil {
		return nil, errors.New("failed to execute currency prices")
	}
	currencies := make(map[string]*entity.Currency, len(res.Data))
	for i, v := range res.Data {
		currencies[i] = &entity.Currency{
			MaxSupply: int(v),
		}
	}
	return currencies, nil
}
