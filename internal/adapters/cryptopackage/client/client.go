package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	url    string
	apiKey string
}

func New(url string, apiKey string) *Client {
	if url == "" {
		url = defualtURl
	}
	return &Client{
		url:    url,
		apiKey: apiKey,
	}
}

// return currency price in map where
// Key = coin such as BTC,ETH etc..
// Value is map of prices such as USD EUR
func (c *Client) GetCurrentPrices(coins []string) (map[string]Coin, error) {
	coinsParam := strings.Join(coins, ",")
	url := fmt.Sprintf("%s/pricemulti?fsyms=%s&tsyms=USD&api_key=%s", c.url, coinsParam, c.apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var coinPrices map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&coinPrices); err != nil {
		return nil, err
	}

	res := make(map[string]Coin)
	for _, coin := range coins {
		if _, ok := coinPrices[coin]; ok {
			res[coin] = Coin{
				Prices: coinPrices[coin],
			}
		}
	}

	return res, nil
}

// Функция для добавления элемента в OrderedMap
func (m *OrderedMap) Add(coin string, price float64) {
	m.Coins = append(m.Coins, coin)
	m.Data[coin] = price
}

// Функция для получения значений в порядке добавления
func (m *OrderedMap) Values() []float64 {
	var values []float64
	for _, coin := range m.Coins {
		values = append(values, m.Data[coin])
	}
	return values
}
