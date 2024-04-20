package cryptoparser

import (
	"context"
	"sync"
	"time"

	"github.com/docobro/dimploma_project/internal/handler/manager"
)

type UseCase interface{}

// usecase operation with things that we will parse per time
// usecase операции будут являться объектом парсинга сервиса Parser

// mn manager that will be appending items in local store. then bathing into database store
// mn manager оператор который будет собирать данные из parser в локальное хранилище а далее раз в n веремени отправлять эти данные кучей в базу данных
type Parser struct {
	uc            UseCase
	ctx           context.Context
	mn            *manager.Manager
	cancel        context.CancelFunc
	flushInterval time.Duration
	mu            sync.Mutex
}

func NewParser(uc UseCase, mn *manager.Manager, flushInterval time.Duration) *Parser {
	ctx, cancel := context.WithCancel(context.Background())

	return &Parser{
		uc:            uc,
		mn:            mn,
		flushInterval: flushInterval,
		ctx:           ctx,
		cancel:        cancel,
	}
}
