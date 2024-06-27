package app

import (
	"fmt"

	"github.com/docobro/dimploma_project/internal/config"
	"github.com/docobro/dimploma_project/pkg/clickhouse"
	"github.com/docobro/dimploma_project/pkg/logger"
)

type App struct {
	l   logger.Logger
	cfg *config.Config
	c   *Container
}

func New(configpath string) (*App, error) {
	config, err := config.New(configpath)
	if err != nil {
		return nil, fmt.Errorf("app - New - config.New: %w", err)
	}

	app := &App{
		cfg: config,
	}

	err = app.initLogger()
	if err != nil {
		return nil, err
	}

	// create connection to database
	connStr := fmt.Sprintf("clickhouse://%v:%v?username=%v&x-multi-statement=true&password=%v&database=%v;", config.SQLConfig.Host, config.SQLConfig.Port, config.User, config.Password, config.DBName)
	pg, err := clickhouse.New(&clickhouse.Config{
		Host:     config.SQLConfig.Host,
		Port:     config.SQLConfig.Port,
		Username: config.User,
		Password: config.Password,
		Database: config.DBName,
	})
	if err != nil {
		return nil, fmt.Errorf("app - clickhouse.New - failed to connect to clickhouse with error:%v", err)
	}
	runMigrations(connStr)
	// init usecases
	app.c = NewContainer(&pg.Conn, app.cfg.CryptoConfig)

	return app, nil
}
