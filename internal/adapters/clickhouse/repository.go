package clickhouse

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/docobro/dimploma_project/internal/entity"
	orm "github.com/docobro/dimploma_project/internal/orm/raw"
	"github.com/google/uuid"
)

type Repository struct {
	conn driver.Conn
}

func New(clickhouse *driver.Conn) *Repository {
	if clickhouse == nil {
		return nil
	}
	return &Repository{
		conn: *clickhouse,
	}
}

func (repository *Repository) Insert(rows entity.Rows) error {
	log.Printf("mock inserting rows:%v", rows)
	return nil
}

func (r *Repository) GetCryptoTokens() (map[string]entity.Currency, error) {
	var cols []orm.Currency
	getTokenQuery := "SELECT * FROM cryptowallet.currencies"
	err := r.conn.Select(context.Background(), &cols, getTokenQuery)
	if err != nil {
		log.Printf("clickhouse_repo Query err:%v\n", err)
		return nil, err
	}
	res := make(map[string]entity.Currency)
	for _, v := range cols {
		res[v.Name] = entity.Currency{
			ID:   v.ID,
			Name: v.Name,
			Code: v.Code,
		}
	}
	return res, nil
}

// insertQuery = `INSERT INTO %s (id, currencyName, price)

// добавление индексов
func (r *Repository) CreateIndices(indices []entity.Indices) error {
	if len(indices) == 0 {
		return errors.New("nothing to insert. abort")
	}
	tokens, err := r.GetCryptoTokens()
	if err != nil {
		return err
	}

	batch, err := r.conn.PrepareBatch(context.Background(), "INSERT INTO indices")
	if err != nil {
		return err
	}

	for _, v := range indices {
		log.Printf("crypto_name:%v price:%v", v.CryptoName, v.Price.Value)
		err := batch.Append(uuid.New(), v.Price.Value, v.Volume.Value, time.Now(), tokens[v.CryptoName].ID)
		if err != nil {
			return err
		}
	}

	err = batch.Send()
	if err != nil {
		return err
	}

	if ok := batch.IsSent(); !ok {
		return errors.New("batch is not sended")
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
