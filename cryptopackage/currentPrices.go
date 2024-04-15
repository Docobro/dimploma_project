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
