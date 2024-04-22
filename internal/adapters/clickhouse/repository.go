package clickhouse

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/docobro/dimploma_project/internal/entity"
	"github.com/google/uuid"
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

// insertQuery = `INSERT INTO %s (id, currencyName, price)
func (r *Repository) AddCurrencies( /*rows rows*/ ) error {
	currencies := make([]entity.Coin, 10)
	batch, err := r.conn.PrepareBatch(context.Background(), "INSERT INTO")
	if err != nil {
		return err
	}

	for _, v := range currencies {
		err := batch.Append(uuid.New(), v.Name, v.Prices["USD"])
		if err != nil {
			return err
		}
	}

	err = batch.Send()
	if err != nil {
		return err
	}

	return nil
}
