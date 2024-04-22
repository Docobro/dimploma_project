package clickhouse

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/docobro/dimploma_project/internal/entity"
)

type cryptoAdapterRepo interface {
	GetCurrencies(coins []string) (map[string]*entity.Coin, error)
}
type Repository struct {
	conn driver.Conn
	cryptoAdapterRepo
}

func New(clickhouse *driver.Conn, adapterCrypto cryptoAdapterRepo) *Repository {
	return &Repository{
		cryptoAdapterRepo: adapterCrypto,
	}
}

func (repository *Repository) Insert(rows entity.Rows) error {
	panic("not implemented") // TODO: Implement
}

func (r *Repository) GetCurrencies() (*[]entity.Currency, error) {
	return nil, nil
}
