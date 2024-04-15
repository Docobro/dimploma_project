package sqlc

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/docobro/dimploma_project/internal/entity"
	orm "github.com/docobro/dimploma_project/internal/orm/raw"
)

type cryptoAdapterRepo interface {
	GetCurrencies(coins []string) (map[string]*entity.Currency, error)
}
type Repository struct {
	orm orm.Store
	cryptoAdapterRepo
}

func New(clickhouse *driver.Conn, adapterCrypto cryptoAdapterRepo) *Repository {
	return &Repository{
		orm:               orm.NewStore(*clickhouse),
		cryptoAdapterRepo: adapterCrypto,
	}
}

func (r *Repository) GetCurrencies() (*[]entity.Currency, error) {
	return nil, nil
}

func (r *Repository) StoreAdapterCurrencies() error {
	currencies, err := r.cryptoAdapterRepo.GetCurrencies([]string{"USDT", "BTC", "ETH"})
	if err != nil {
		return err
	}
	return r.orm.SaveCurrencies(currencies)
}
