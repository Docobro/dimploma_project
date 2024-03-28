package app

import "github.com/docobro/dimploma_project/internal/db"

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() error {
	db.Generate()
	return nil
}
