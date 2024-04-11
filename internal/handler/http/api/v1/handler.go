package v1

import "github.com/docobro/dimploma_project/pkg/logger"

type UseCase interface{}

type Handler struct {
	uc     UseCase
	logger logger.Logger
}

func NewHandler(uc UseCase, logs logger.Logger) *Handler {
	return &Handler{uc: uc, logger: logs}
}
