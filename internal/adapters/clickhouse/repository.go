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
	"github.com/shopspring/decimal"
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
		err := batch.Append(uuid.New(), v.Volume.Value, time.Now(), tokens[v.CryptoName].ID, v.Price.Value)
		if err != nil {
			return err
		}
	}

	err = batch.Send()
	if err != nil {
		return err
	}

	if ok := batch.IsSent(); !ok {
		return errors.New("batch is not sent")
	}

	return nil
}

// добавление объемов торгов за 1 час
func (r *Repository) CreateVolumes1m(volume map[string]float32) error {
	tokens, err := r.GetCryptoTokens()
	if err != nil {
		return err
	}

	batch, err := r.conn.PrepareBatch(context.Background(), "INSERT INTO trade_volume_1h")
	if err != nil {
		return err
	}

	for k, v := range volume {
		err := batch.Append(uuid.New(), tokens[k].ID, float32(v), time.Now())
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
func (r *Repository) CreateTransaction(transaction map[string]uint32) error {
	tokens, err := r.GetCryptoTokens()
	if err != nil {
		return err
	}
	batch, err := r.conn.PrepareBatch(context.Background(), "INSERT INTO transaction_per_day")
	if err != nil {
		return err
	}

	for k, v := range transaction {
		err := batch.Append(uuid.New(), int64(v), time.Now(), tokens[k].ID)
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

	batch, err := r.conn.PrepareBatch(context.Background(), "INSERT INTO supplies")
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

type rf struct {
	Value     decimal.Decimal `ch:"value"`
	Name      string          `ch:"name"`
	CreatedAt time.Time       `ch:"created_at"`
}

func (r *Repository) GetPrices(coins []string, start time.Time, end time.Time) interface{} {
	getPricesQuery := `SELECT p.value,c.name,p.created_at  
  FROM cryptowallet.prices p
  INNER JOIN cryptowallet.currencies c on c.id = p.crypto_id
  WHERE 
  p.created_at >= $1
  AND p.created_at <= $2
  AND  c.name IN($3)
  ORDER BY p.created_at DESC LIMIT $4`

	// price index = current_price / price_custom_duration_ago

	// volume_hour = https://min-api.cryptocompare.com/data/pricemultifull?fsyms=BTC&tsyms=USD

	// aggregated_price_index = current_volume / volume_24h_ago

	var cols []rf
	rows, err := r.conn.Query(context.Background(), getPricesQuery, start, end, coins, len(coins))
	if err != nil {
		return err
	}

	for rows.Next() {
		var i rf
		err := rows.Scan(
			&i.Value,
			&i.Name,
			&i.CreatedAt,
		)
		if err != nil {
			return err
		}
		cols = append(cols, i)
	}

	return cols
}
