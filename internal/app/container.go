package app

import (
	"log"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/docobro/dimploma_project/internal/adapters/clickhouse"
	"github.com/docobro/dimploma_project/internal/adapters/cryptopackage"
	"github.com/docobro/dimploma_project/internal/config"
	"github.com/docobro/dimploma_project/internal/usecase"
)

var (
	ClickhouseRepo string = "clickhouse_repo"
	CryptoRepo     string = "crypto_repo"
)

type Container struct {
	deps map[string]interface{}
}

func NewContainer(clickhouse *driver.Conn, cfg config.CryptoConfig) *Container {
	c := &Container{
		deps: make(map[string]interface{}),
	}
	c.deps[ClickhouseRepo] = c.GetClickhouseRepo(clickhouse)
	c.deps[CryptoRepo] = c.GetCryptoRepo(cfg)
	return c
}

func (c *Container) GetUseCase() *usecase.Usecase {
	clickhouseRepo, ok := c.deps[ClickhouseRepo].(*clickhouse.Repository)
	if !ok {
		log.Fatal("container - c.deps[ClickhouseRepo] failed to get clickhouse repo. repo were not inited")
	}
	cryptoRepo, ok := c.deps[CryptoRepo].(*cryptopackage.Repository)
	if !ok {
		log.Fatal("container - c.deps[CryptoRepo] failed to get crypto repo. repo were not inited")
	}
	return usecase.New(clickhouseRepo, cryptoRepo)
}

func (c *Container) GetClickhouseRepo(driver *driver.Conn) *clickhouse.Repository {
	repo := clickhouse.New(driver)
	c.deps[ClickhouseRepo] = repo
	return repo
}

func (c *Container) GetCryptoRepo(config config.CryptoConfig) *cryptopackage.Repository {
	return cryptopackage.New(config)
}
