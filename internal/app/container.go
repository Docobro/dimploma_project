package app

import (
	"log"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/docobro/dimploma_project/internal/adapters/clickhouse"
	"github.com/docobro/dimploma_project/internal/adapters/cryptopackage"
	"github.com/docobro/dimploma_project/internal/usecase"
)

var ClickhouseRepo string = "clickhouse_repo"

type Container struct {
	deps map[string]interface{}
}

func NewContainer(clickhouse *driver.Conn) *Container {
	c := &Container{
		deps: make(map[string]interface{}),
	}
	c.deps[ClickhouseRepo] = c.GetClickhouseRepo(clickhouse)
	return c
}

func (c *Container) GetUseCase() *usecase.Usecase {
	if _, ok := c.deps[ClickhouseRepo]; !ok {
		log.Fatal("container - c.deps[ClickhouseRepo] failed to get clickhouse repo. repo were not inited")
	}
	return usecase.New(c.deps[ClickhouseRepo])
}

func (c *Container) GetClickhouseRepo(driver *driver.Conn) *clickhouse.Repository {
	return clickhouse.New(driver, c.GetCryptoRepo(""))
}

func (c *Container) GetCryptoRepo(url string) *cryptopackage.Repository {
	return cryptopackage.New(url)
}
