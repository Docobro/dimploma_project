package clickhouse

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/docobro/dimploma_project/internal/entity"
	orm "github.com/docobro/dimploma_project/internal/orm/raw"
	"github.com/google/uuid"
)

type cryptoAdapterRepo interface {
	GetCurrencies(coins []string) (map[string]*entity.Coin, error)
	GetPriceIndex() (uint, error)
	GetVolumeIndex() (uint, error)
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

func (r *Repository) GetCryptoTokens() (map[string]entity.Currency, error) {
	getTokenQuery := "SELECT * FROM currencies"
	rows, err := r.conn.Query(context.Background(), getTokenQuery)
	if err != nil {
		return nil, err
	}
	res := make(map[string]entity.Currency)
	for rows.Next() {
		var col1 orm.Currency
		if err := rows.Scan(&col1); err != nil {
			return nil, err
		}
		res[col1.Name] = entity.Currency{
			ID:   uint64(col1.ID),
			Name: col1.Name,
			Code: col1.Code,
		}
	}
	rows.Close()
	return res, nil
}

// insertQuery = `INSERT INTO %s (id, currencyName, price)
func (r *Repository) CreateIndices(indices []entity.Indices) error {
	tokens, err := r.GetCryptoTokens()
	if err != nil {
		return err
	}

	batch, err := r.conn.PrepareBatch(context.Background(), "INSERT INTO indices(id,crypto_id,price_index,volume_index,created_at)")
	if err != nil {
		return err
	}

	for _, v := range indices {
		err := batch.Append(uuid.New(), tokens[v.CryptoName].ID, v.Price, v.Volume, time.Now())
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
