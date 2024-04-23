package entity

type Currency struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type Coin struct {
	Prices    map[string]float64
	MarketCap uint64
	Name      string
}

type Indices struct {
	CryptoName string
	Volume     VolumeIndex
	Price      PriceIndex
}

type VolumeIndex struct {
	Value float64
}

type PriceIndex struct {
	Value float64
}

type Volume struct {
	CryptoName string
	Value      float64
}

type Transaction struct {
	CryptoName string
	Value      uint64
}

type Supplies struct {
	CryptoName string
	Value      uint64
}
