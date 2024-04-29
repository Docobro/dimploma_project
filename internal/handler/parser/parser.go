package parser

import (
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
				fn()
			case <-j.quit: // Check if the quit signal is received
				return
			}
		}
	}()
}

//func (m *Parser) startParsing() {
//	log.Println("Start parsing crypto remote!")
//
//	coins := []string{"BTC", "ETH", "USDT", "TON"}
//	currencies, err := m.uc.GetCurrencies(coins)
//	if err != nil {
//		log.Printf("Failed to parse Currencies err:%v", err)
//	}
//
//	priceIndices, err := m.uc.GetPriceIndices(coins)
//	if err != nil {
//		log.Printf("Failed to parse Currencies err:%v", err)
//	}
//
//	volumeIndices, err := m.uc.GetVolumeIndices(coins)
//	if err != nil {
//		log.Printf("Failed to parse Currencies err:%v", err)
//	}
//	for _, v := range currencies {
//		statsKey := entity.NewKey(entity.Key{
//			Timestamp: time.Now().Unix(),
//		})
//		statsValue := entity.Value{
//			Coin:        v,
//			CoinName:    v.Name,
//			PriceIndex:  priceIndices[v.Name],
//			VolumeIndex: volumeIndices[v.Name],
//		}
//		m.mn.Append(statsKey, statsValue)
//	}
//
//	log.Printf("Parser has parsed!")
//}
