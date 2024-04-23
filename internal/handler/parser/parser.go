package parser

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/docobro/dimploma_project/internal/entity"
	"github.com/docobro/dimploma_project/internal/handler/manager"
)

type UseCase interface {
	GetCurrencies(coins []string) (map[string]*entity.Coin, error)
	GetPriceIndices(coins []string) (map[string]float64, error)
	GetVolumeIndices(coins []string) (map[string]float64, error)
}

// usecase operation with things that we will parse per time
// usecase операции будут являться объектом парсинга сервиса Parser

// mn manager that will be appending items in local store. then bathing into database store
// mn manager оператор который будет собирать данные из parser в локальное хранилище а далее раз в n веремени отправлять эти данные кучей в базу данных

// parser does not use own append function he operates with manager append ops
type Parser struct {
	uc            UseCase
	mn            *manager.Manager
	ctx           context.Context
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

func (m *Parser) Start() {
	log.Println("Parser loop started")
	go m.loop()
}

func (m *Parser) Stop() {
	m.cancel()
	log.Println("Parser has stopped!")
}

func (m *Parser) loop() {
	for {
		select {
		case <-time.After(m.flushInterval):
			m.startParsing()

		case <-m.ctx.Done():
			m.startParsing()
			return
		}
	}
}

func (m *Parser) startParsing() {
	log.Println("Parser crypto parsing")

	coins := []string{"BTC", "ETH", "TON"}
	currencies, err := m.uc.GetCurrencies(coins)
	if err != nil {
		log.Printf("Failed to parse Currencies err:%v", err)
	}

	priceIndices, err := m.uc.GetPriceIndices(coins)
	if err != nil {
		log.Printf("Failed to parse Currencies err:%v", err)
	}

	volumeIndices, err := m.uc.GetVolumeIndices(coins)
	if err != nil {
		log.Printf("Failed to parse Currencies err:%v", err)
	}
	for _, v := range currencies {
		statsKey := entity.NewKey(entity.Key{
			Timestamp: time.Now().Unix(),
		})
		statsValue := entity.Value{
			Requests:    1,
			Coin:        v,
			PriceIndex:  priceIndices[v.Name],
			VolumeIndex: volumeIndices[v.Name],
		}
		m.mn.Append(statsKey, statsValue)
	}

	log.Printf("Parser has parsed!")
}
