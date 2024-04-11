package sqlc

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	orm "github.com/docobro/dimploma_project/internal/orm/raw"
)

type Repository struct {
	orm orm.Store
}

func New(postgres *driver.Conn) *Repository {
	return &Repository{orm: orm.NewStore(*postgres)}
}

func (r *Repository) Aboa() {
}
