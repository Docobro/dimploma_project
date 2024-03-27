package cryptopackage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	baseURL = "https://min-api.cryptocompare.com/data"
)

// Структура для хранения цены криптовалюты
type CoinPrice struct {
	Coin  string  `json:"Coin"`
	Price float64 `json:"USD"`
}

// Структура для хранения порядка криптовалют
type OrderedMap struct {
	coins []string
	Data  map[string]float64
}

// Функция для создания OrderedMap
func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		coins: make([]string, 0),
		Data:  make(map[string]float64),
	}
}

// Функция для добавления элемента в OrderedMap
func (m *OrderedMap) Add(coin string, price float64) {
	m.coins = append(m.coins, coin)
	m.Data[coin] = price
}

// Функция для получения значений в порядке добавления
func (m *OrderedMap) Values() []float64 {
	var values []float64
	for _, coin := range m.coins {
		values = append(values, m.Data[coin])
	}
	return values
}

// Функция для получения текущей цены для нескольких криптовалют
func GetCurrentPrices(coins []string) (*OrderedMap, error) {
	coinsParam := strings.Join(coins, ",")
	url := fmt.Sprintf("%s/pricemulti?fsyms=%s&tsyms=USD", baseURL, coinsParam)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var coinPrices map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&coinPrices); err != nil {
		return nil, err
	}

	orderedMap := NewOrderedMap()
	for _, coin := range coins {
		if priceData, ok := coinPrices[coin]; ok {
			orderedMap.Add(coin, priceData["USD"])
		}
	}

	return orderedMap, nil
}
