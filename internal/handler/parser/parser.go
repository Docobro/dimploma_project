package parser

import (
	"fmt"
	"sync"
	"time"

	"github.com/docobro/dimploma_project/internal/entity"
	"github.com/docobro/dimploma_project/internal/handler/manager"
)

type UseCase interface {
	GetCurrencies(coins []string) (map[string]*entity.Coin, error)
	GetPriceIndices(coins []string) (map[string]int32, error)
	GetVolumeIndices(coins []string) (map[string]float32, error)
}

// usecase operation with things that we will parse per time
// usecase операции будут являться объектом парсинга сервиса Parser

// mn manager that will be appending items in local store. then bathing into database store
// mn manager оператор который будет собирать данные из parser в локальное хранилище а далее раз в n веремени отправлять эти данные кучей в базу данных

// parser does not use own append function he operates with manager append ops
type Parser struct {
	mn      *manager.Manager
	wg      *sync.WaitGroup
	quit    chan struct{}
	mu      sync.Mutex
	running bool
}

func NewParser(mn *manager.Manager) *Parser {
	return &Parser{
		mn:   mn,
		wg:   &sync.WaitGroup{},
		quit: make(chan struct{}),
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
					fmt.Printf("err: %v\n", err)
				}
			case <-j.quit: // Check if the quit signal is received
				return
			}
		}
	}()
}
