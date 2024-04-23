package app

import (
	"fmt"
	"time"

	"github.com/docobro/dimploma_project/internal/config"
	"github.com/docobro/dimploma_project/internal/handler/manager"
	"github.com/docobro/dimploma_project/internal/handler/parser"
	"github.com/docobro/dimploma_project/pkg/clickhouse"
	"github.com/docobro/dimploma_project/pkg/logger"
	"github.com/hanagantig/gracy"
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
	connStrf := "clickhouse://localhost:9000?username=default&x-multi-statement=true&password=&database=cryptowallet;"
	_ = fmt.Sprintf("clickhouse://%v:%v?username=%v&password=%v&database=%v", config.SQLConfig.Host, config.SQLConfig.Port, config.User, config.Password, config.DBName)
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
	runMigrations(connStrf)
	// init usecases
	app.c = NewContainer(&pg.Conn, app.cfg.CryptoConfig)

	mn := manager.NewManager(app.c.GetClickhouseRepo(&pg.Conn), time.Second*10)
	gracy.AddCallback(func() error {
		mn.Stop()
		return nil
	})

	mn.Start()

	parser := parser.NewParser(app.c.GetUseCase(), mn, time.Second*5)
	gracy.AddCallback(func() error {
		parser.Stop()
		return nil
	})

	parser.Start()

	return app, nil
}
