package app

import (
	"time"

	"github.com/docobro/dimploma_project/internal/handler/parser"
	"github.com/hanagantig/gracy"
)

func (a *App) StartParser() error {
	go func() {
		a.startParser()
	}()

	a.l.Info("parser gracefully stoped")
	return nil
}

func (a *App) startParser() {
	parser := parser.NewParser(a.l)
	gracy.AddCallback(func() error {
		parser.Stop()
		return nil
	})

	uc := a.c.GetUseCase()
	parser.Run(uc.ParsePrices, time.Second*10)
	parser.Run(uc.ParseVolumeMinute, time.Minute*1)
	parser.Run(uc.CreateIndices, time.Second*10)
	parser.Run(uc.ParseTotalSupply, time.Minute*2)
	parser.Run(uc.ParseVolatility, time.Minute*1)
	parser.Run(uc.ParsePearson, time.Minute*1)
	parser.Run(uc.ParsePredictions, time.Minute*1)
}
