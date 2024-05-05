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

func (c *Client) GetTransactionData(coins []string) (map[string]int, error) {
	symbolsParam := strings.Join(coins, ",")
	url := fmt.Sprintf("%s/data/blockchain/histo/day?fsym=%s&limit=1", c.url, symbolsParam)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("authorization", "Apikey "+c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	var response map[string]map[string][]struct {
		TransactionCount int `json:"transaction_count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	transactionCounts := make(map[string]int)
	for symbol, data := range response {
		if len(data["Data"]) > 0 {
			transactionCounts[symbol] = data["Data"][0].TransactionCount
		}
	}

	return transactionCounts, nil
}

// return currency price in map where
// Key = coin such as BTC,ETH etc..
// Value is map of prices such as USD EUR
func (c *Client) GetPriceMultiFull(coins []string, currencies []string) (PriceMultiFull, error) {
	coinsParam := strings.Join(coins, ",")
	currencyParam := strings.Join(currencies, ",")
	url := fmt.Sprintf("%s/pricemultifull?fsyms=%v&tsyms=%v", c.url, coinsParam, currencyParam)
	resp, err := http.Get(url)
	if err != nil {
		return PriceMultiFull{}, err
	}
	// Time.UnmarshalJSON: input is not a JSON string
	defer resp.Body.Close()
	var coinPrices PriceMultiFull
	if err := json.NewDecoder(resp.Body).Decode(&coinPrices); err != nil {
		return PriceMultiFull{}, err
	}

	return coinPrices, nil
}

func (c *Client) GetCurrentPrices(coins []string) (map[string]Coin, error) {
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
