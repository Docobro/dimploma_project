package app

import (
	"log"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/docobro/dimploma_project/internal/adapters/sqlc"
	"github.com/docobro/dimploma_project/internal/usecase"
)

var PostgresRepo string = "postgres_repo"

type Container struct {
	deps map[string]interface{}
}

func NewContainer(postgres *driver.Conn) *Container {
	c := &Container{
		deps: make(map[string]interface{}),
	}
	c.deps[PostgresRepo] = c.GetClickhouseRepo(postgres)
	return c
}

func (c *Container) GetUseCase() *usecase.Usecase {
	if _, ok := c.deps[PostgresRepo]; !ok {
		log.Fatal("container - c.deps[PostgresRepo] failed to get postgres repo. repo were not inited")
	}
	return usecase.New(c.deps[PostgresRepo])
}

func (c *Container) GetClickhouseRepo(driver *driver.Conn) *sqlc.Repository {
	return sqlc.New(driver)
}
