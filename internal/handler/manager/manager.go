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

//type Price struct {
//	ID        int64            `json:"id"`
//	CryptoID  int64            `json:"crypto_id"`
//	Value     pgtype.Numeric   `json:"value"`
//	MarketCap int32            `json:"market_cap"`
//	TimeDiff  interface{}      `json:"time_diff"`
//	CreatedAt pgtype.Timestamp `json:"created_at"`
//}

type (
	Writer interface {
		Insert(rows entity.Rows) error
		CreateIndices(indices []entity.Indices) error
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
	log.Println("Start stats inserting")

	rows := m.withdraw()
	if len(rows) == 0 {
		log.Println("No stats rows, skipping")
		return
	}

	indices := make([]entity.Indices, 0)
	for _, v := range rows {
		indices = append(indices, entity.Indices{
			CryptoName: v.CoinName,
			Price:      entity.PriceIndex{Value: v.Price},
			Volume:     entity.VolumeIndex{Value: v.VolumeIndex},
		})
	}
	if err := m.writer.CreateIndices(indices); err != nil {
		log.Println("failed to create indices err:" + err.Error())
	}

	if err := m.writer.Insert(rows); err != nil {
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
