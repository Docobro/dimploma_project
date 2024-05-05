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

type (
	CoinType       string
	CurrencyType   string
	PriceMultiFull struct {
		RAW map[CoinType]map[CurrencyType]PriceDetails `json:"RAW"`
	}
)

type TransactionResponse struct {
	Data struct {
		Data []struct {
			TransactionCount int `json:"transaction_count"`
		} `json:"data"`
	} `json:"Data"`
}

// USDDetails represents the details for USD
type PriceDetails struct {
	ConversionLastUpdate    uint64  `json:"CONVERSIONLASTUPDATE"`
	LastUpdate              uint64  `json:"LASTUPDATE"`
	Type                    string  `json:"TYPE"`
	Market                  string  `json:"MARKET"`
	FromSymbol              string  `json:"FROMSYMBOL"`
	ToSymbol                string  `json:"TOSYMBOL"`
	Flags                   string  `json:"FLAGS"`
	LastMarket              string  `json:"LASTMARKET"`
	ImageURL                string  `json:"IMAGEURL"`
	ConversionSymbol        string  `json:"CONVERSIONSYMBOL"`
	LastTradeID             string  `json:"LASTTRADEID"`
	ConversionType          string  `json:"CONVERSIONTYPE"`
	LastVolume              float64 `json:"LASTVOLUME"`
	Change24Hour            float64 `json:"CHANGE24HOUR"`
	VolumeHour              float64 `json:"VOLUMEHOUR"`
	VolumeHourTo            float64 `json:"VOLUMEHOURTO"`
	OpenHour                float64 `json:"OPENHOUR"`
	HighHour                float64 `json:"HIGHHOUR"`
	LowHour                 float64 `json:"LOWHOUR"`
	VolumeDay               float64 `json:"VOLUMEDAY"`
	VolumeDayTo             float64 `json:"VOLUMEDAYTO"`
	OpenDay                 float64 `json:"OPENDAY"`
	HighDay                 float64 `json:"HIGHDAY"`
	LowDay                  float64 `json:"LOWDAY"`
	Volume24Hour            float64 `json:"VOLUME24HOUR"`
	Volume24HourTo          float64 `json:"VOLUME24HOURTO"`
	Open24Hour              float64 `json:"OPEN24HOUR"`
	High24Hour              float64 `json:"HIGH24HOUR"`
	Low24Hour               float64 `json:"LOW24HOUR"`
	LastVolumeTo            float64 `json:"LASTVOLUMETO"`
	ChangePct24Hour         float64 `json:"CHANGEPCT24HOUR"`
	ChangeDay               float64 `json:"CHANGEDAY"`
	ChangePctDay            float64 `json:"CHANGEPCTDAY"`
	ChangeHour              float64 `json:"CHANGEHOUR"`
	ChangePctHour           float64 `json:"CHANGEPCTHOUR"`
	Price                   float64 `json:"PRICE"`
	TopTierVolume24HourTo   float64 `json:"TOPTIERVOLUME24HOURTO"`
	TopTierVolume24Hour     float64 `json:"TOPTIERVOLUME24HOUR"`
	Supply                  float64 `json:"SUPPLY"`
	MktCap                  float64 `json:"MKTCAP"`
	MktCapPenalty           float64 `json:"MKTCAPPENALTY"`
	CirculatingSupply       float64 `json:"CIRCULATINGSUPPLY"`
	CirculatingSupplyMktCap float64 `json:"CIRCULATINGSUPPLYMKTCAP"`
	TotalVolume24H          float64 `json:"TOTALVOLUME24H"`
	TotalVolume24HTo        float64 `json:"TOTALVOLUME24HTO"`
	TotalTopTierVolume24H   float64 `json:"TOTALTOPTIERVOLUME24H"`
	TotalTopTierVolume24HTo float64 `json:"TOTALTOPTIERVOLUME24HTO"`
	Median                  float64 `json:"MEDIAN"`
}
