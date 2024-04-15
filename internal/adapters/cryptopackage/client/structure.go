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
	coins []string
	Data  map[string]float64
}
