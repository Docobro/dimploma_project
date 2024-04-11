package app

import (
	"fmt"

	"github.com/docobro/dimploma_project/pkg/logger"
)

func (a *App) initLogger() error {
	logger, err := logger.NewLogger()
	if err != nil {
		return fmt.Errorf("app - initLogger - logger.NewLogger: %w", err)
	}
	a.l = logger
	return nil
}
