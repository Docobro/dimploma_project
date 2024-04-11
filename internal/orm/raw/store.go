package orm

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Store interface {
	Querier
}
type SQLStore struct {
	*Queries
	db driver.Conn
}

func NewStore(db driver.Conn) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
