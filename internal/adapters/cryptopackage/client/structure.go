package client

const (
	defualtURl = "https://min-api.cryptocompare.com/data"
)

// Структура для хранения цены криптовалюты
type CoinPrice struct {
	Coin  string  `json:"Coin"`
	Price float64 `json:"USD"`
}

// Структура для хранения порядка криптовалют
type OrderedMap struct {
	Coins []string
	Data  map[string]float64
}

type (
	Coin struct {
		Prices map[string]float64
	}
)
