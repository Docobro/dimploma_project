package app

import (
	"fmt"
	client "net/http"

	"github.com/docobro/dimploma_project/internal/handler/http"
	v1 "github.com/docobro/dimploma_project/internal/handler/http/api/v1"
	"github.com/hanagantig/gracy"
	"go.uber.org/zap"
)

func (a *App) StartHTTPServer() error {
	go func() {
		a.startHTTPServer()
	}()

	a.l.Info("http server gracefully stopped")
	return nil
}

func (a *App) startHTTPServer() {
	handler := v1.NewHandler(a.c.GetUseCase(), a.l)
	router := http.NewRouter()
	router.WithHandler(handler, a.l)

	srv := http.NewServer(a.cfg.HTTPConfig)
	srv.RegisterRoutes(router)

	gracy.AddCallback(func() error {
		return srv.Stop()
	})

	a.l.Info(fmt.Sprintf("starting HTTP server at %s:%s", a.cfg.HTTPConfig.Host, a.cfg.HTTPConfig.Port))
	err := srv.Start()
	if err != nil {
		a.l.Fatal("Fail to start %s http server:", zap.String("app", a.cfg.Name), zap.Error(err))
	}
}

func (a *App) newHttpClient() *client.Client {
	client := client.Client{}
	return &client
}
