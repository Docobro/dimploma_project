package entity

import "github.com/google/uuid"

type Currency struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
}

type Coin struct {
	Name      string
	Prices    map[string]float64
	MarketCap int32
}

type Indices struct {
	CryptoName string
	Volume     VolumeIndex
	Price      PriceIndex
}

type VolumeIndex struct {
	Value float32
}

type PriceIndex struct {
	Value int32
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
