package clickhouse

import (
	"context"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/docobro/dimploma_project/internal/entity"
	orm "github.com/docobro/dimploma_project/internal/orm/raw"
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
	log.Printf("mock inserting rows:%v", rows)
	return nil
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

// добавление индексов
func (r *Repository) CreateIndices(indices []entity.Indices) error {
	tokens, err := r.GetCryptoTokens()
	if err != nil {
		return err
	}

	batch, err := r.conn.PrepareBatch(context.Background(), "INSERT INTO indices(id, crypto_id, price_index, volume_index, created_at)")
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

// добавление текущих цен и маркет капов
func (r *Repository) CreatePrices(prices []entity.Coin) error {
	tokens, err := r.GetCryptoTokens()
	if err != nil {
		return err
	}

	batch, err := r.conn.PrepareBatch(context.Background(), "INSERT INTO prices(id, crypto_id, value, market_cap, time_diff , created_at)")
	if err != nil {
		return err
	}

	for _, v := range prices {
		err := batch.Append(uuid.New(), tokens[v.Name].ID, v.Prices, v.MarketCap, time.Now(), time.Now())
		// формат для time_diff нужно поменять на что-то другое !!!!!!
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

// добавление объемов торгов за 1 час
func (r *Repository) CreateVolumes1h(volume []entity.Volume) error {
	tokens, err := r.GetCryptoTokens()
	if err != nil {
		return err
	}

	batch, err := r.conn.PrepareBatch(context.Background(), "INSERT INTO trade_volume_1h(id, crypto_id, trade_volume, time_diff)")
	if err != nil {
		return err
	}

	for _, v := range volume {
		err := batch.Append(uuid.New(), tokens[v.CryptoName].ID, v.Value, time.Now())
		// формат для time_diff нужно поменять на что-то другое !!!!!!
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

// добавление транзакций за день
func (r *Repository) CreateTransaction(transaction []entity.Transaction) error {
	tokens, err := r.GetCryptoTokens()
	if err != nil {
		return err
	}

	batch, err := r.conn.PrepareBatch(context.Background(), "INSERT INTO transaction_per_day(id, crypto_id, trans_volume, created_at)")
	if err != nil {
		return err
	}

	for _, v := range transaction {
		err := batch.Append(uuid.New(), tokens[v.CryptoName].ID, v.Value, time.Now())
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

// добавление текущий саплаев
func (r *Repository) CreateSupplies(supplies []entity.Supplies) error {
	tokens, err := r.GetCryptoTokens()
	if err != nil {
		return err
	}

	batch, err := r.conn.PrepareBatch(context.Background(), "INSERT INTO supplies(id, crypto_id, total_supply, created_at)")
	if err != nil {
		return err
	}

	for _, v := range supplies {
		err := batch.Append(uuid.New(), tokens[v.CryptoName].ID, v.Value, time.Now())
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
