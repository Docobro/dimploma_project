package client

import (
	"encoding/json"
	"fmt"

	"net/http"
	"strconv"
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

func removeLastElement(response *MinuteResponse) {
	dataLen := len(response.Data.Data)
	if dataLen > 0 {
		response.Data.Data = response.Data.Data[:dataLen-1]
	}
}

func (c *Client) GetOneMinuteFull(coin string, currency string, limit int) (MinuteResponse, error) {
	limits := strconv.Itoa(limit)
	url := fmt.Sprintf("%s/v2/histominute?fsym=%s&tsym=%s&limit=%s", c.url, coin, currency, limits)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return MinuteResponse{}, err
	}

	defer resp.Body.Close()
	var minuteResult MinuteResponse
	if err := json.NewDecoder(resp.Body).Decode(&minuteResult); err != nil {
		fmt.Printf("failed to decode json:err:%v", err)
		return minuteResult, err
	}

	removeLastElement(&minuteResult)
	return minuteResult, nil
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
