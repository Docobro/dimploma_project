package parser

import (
	"sync"
	"time"

	"github.com/docobro/dimploma_project/internal/entity"
	"github.com/docobro/dimploma_project/pkg/logger"
	"go.uber.org/zap"
)

type UseCase interface {
	GetCurrencies(coins []string) (map[string]*entity.Coin, error)
	GetPriceIndices(coins []string) (map[string]int32, error)
	GetVolumeIndices(coins []string) (map[string]float32, error)
}

// usecase operation with things that we will parse per time
// usecase операции будут являться объектом парсинга сервиса Parser

// parser does not use own append function he operates with manager append ops
type Parser struct {
	wg      *sync.WaitGroup
	quit    chan struct{}
	mu      sync.Mutex
	running bool
	l       logger.Logger
}

func NewParser(l logger.Logger) *Parser {
	return &Parser{
		wg:   &sync.WaitGroup{},
		quit: make(chan struct{}),
		l:    l,
	}
}

func (j *Parser) Stop() {
	close(j.quit)
}

func (j *Parser) Wait() {
	j.wg.Wait()
}

func (j *Parser) Run(fn func() error, flushInterval time.Duration) {
	j.wg.Add(1)
	go func() {
		defer j.wg.Done()
		for {
			select {
			case <-time.After(flushInterval):
				err := fn()
				if err != nil {
					j.l.Info("error", zap.Error(err))
				}
			case <-j.quit: // Check if the quit signal is received
				return
			}
		}
	}()
}
