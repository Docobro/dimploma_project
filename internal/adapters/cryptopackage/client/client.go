package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	url string
}

func New(url string) *Client {
	if url == "" {
		url = defualtURl
	}
	return &Client{
		url: url,
	}
}

func (c *Client) GetCurrentPrices(coins []string) (*OrderedMap, error) {
	coinsParam := strings.Join(coins, ",")
	url := fmt.Sprintf("%s/pricemulti?fsyms=%s&tsyms=USD", c.url, coinsParam)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var coinPrices map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&coinPrices); err != nil {
		return nil, err
	}

	orderedMap := newOrderedMap()
	for _, coin := range coins {
		if priceData, ok := coinPrices[coin]; ok {
			orderedMap.Add(coin, priceData["USD"])
		}
	}

	return orderedMap, nil
}

func newOrderedMap() *OrderedMap {
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
