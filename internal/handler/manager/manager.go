package manager

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/docobro/dimploma_project/internal/entity"
)

const (
	defaultCapacity = 1000
)

type (
	Writer interface {
		Insert(rows entity.Rows) error
		CreateIndices(indices []entity.Indices) error
		CreatePrices(prices []entity.Coin) error
	}

	Manager struct {
		writer        Writer
		flushInterval time.Duration

		ctx    context.Context
		cancel context.CancelFunc

		mu   sync.Mutex
		rows entity.Rows
	}
)

func NewManager(w Writer, flushInterval time.Duration) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	return &Manager{
		writer:        w,
		flushInterval: flushInterval,
		ctx:           ctx,
		cancel:        cancel,
		rows:          newRows(),
	}
}

func (m *Manager) Append(k entity.Key, v entity.Value) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.unsafeAppend(k, v)
}

func (m *Manager) AppendRows(rows entity.Rows) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k, v := range rows {
		m.unsafeAppend(k, v)
	}
}

func (m *Manager) unsafeAppend(k entity.Key, v entity.Value) {
	current := m.rows[k]
	current = current.Assign(v)

	m.rows[k] = current
}

func (m *Manager) Start() {
	log.Println("Stats loop started")
	go m.loop()
}

func (m *Manager) Stop() {
	m.cancel()
	log.Println("Manager has stopped!")
}

func (m *Manager) loop() {
	for {
		select {
		case <-time.After(m.flushInterval):
			m.startInserting()

		case <-m.ctx.Done():
			m.startInserting()
			return
		}
	}
}

func (m *Manager) startInserting() {
	log.Println("Start crypto inserting")

	rows := m.withdraw()
	if len(rows) == 0 {
		log.Println("No crypto rows, skipping")
		return
	}

	indices := []entity.Indices{}
	prices := []entity.Coin{}
	for _, v := range rows {
		indices = append(indices, entity.Indices{
			CryptoName: v.CoinName,
			Price:      entity.PriceIndex{Value: v.PriceIndex},
			Volume:     entity.VolumeIndex{Value: v.VolumeIndex},
		})
		prices = append(prices, *v.Coin)
	}
	err := m.writer.CreateIndices(indices)
	if err != nil {
		log.Println("failed to create indices err:" + err.Error())
		m.AppendRows(rows)
		return
	}

	err = m.writer.CreatePrices(prices)
	if err != nil {
		log.Println("failed to create prices err:" + err.Error())
		m.AppendRows(rows)
	}

	err = m.writer.Insert(rows)
	if err != nil {
		log.Printf("Failed to write stats: %v\n", err)
		log.Printf("Return stats rows to map: %d\n", len(rows))

		m.AppendRows(rows)
		return
	}

	log.Printf("Stats rows successfuly written: %d\n", len(rows))
}

func (m *Manager) withdraw() entity.Rows {
	m.mu.Lock()
	defer m.mu.Unlock()

	rows := m.rows
	m.rows = newRows()

	return rows
}

func newRows() entity.Rows {
	return make(entity.Rows, defaultCapacity)
}
