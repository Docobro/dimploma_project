package clickhouse

import (
	"context"
	"fmt"
	"math"
)

func (r *Repository) GettingPriceData(coin string) map[int]float64 {
	tokens, err := r.GetCryptoTokens()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil
	}
	id := tokens[coin].ID
	indexQuery := fmt.Sprintf(`SELECT CAST(value AS Float64) FROM cryptowallet.prices WHERE crypto_id = '%s' ORDER BY created_at DESC LIMIT 300`, id)
	rows, err := r.conn.Query(context.Background(), indexQuery)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil
	}

	values := make(map[int]float64)
	index := 0
	for rows.Next() {
		var value float64
		if err = rows.Scan(&value); err != nil {
			return nil
		}
		values[index] = value
		index++
	}

	averagedValues := make(map[int]float64, 50)
	for i := 0; i < 50; i++ {
		var sum float64
		for j := 0; j < 6; j++ {
			sum += values[i*6+j]
		}
		average := sum / 6
		averagedValues[i] = average
	}

	return averagedValues
}

func (r *Repository) GettingVolumeData(coin string) map[int]float64 {
	tokens, err := r.GetCryptoTokens()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil
	}
	id := tokens[coin].ID
	indexQuery := fmt.Sprintf(`SELECT CAST(trade_volume AS Float64) FROM cryptowallet.trade_volume_1m WHERE crypto_id = '%s' ORDER BY created_at DESC LIMIT 50`, id)
	rows, err := r.conn.Query(context.Background(), indexQuery)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	values := make(map[int]float64)
	index := 0
	for rows.Next() {
		var value float64
		if err = rows.Scan(&value); err != nil {
			return nil
		}
		values[index] = value
		index++
	}

	return values
}

func (r *Repository) PearsonPriceToVolumeCorrelation(coin string) float64 {
	var x, y map[int]float64
	x = r.GettingVolumeData(coin)
	y = r.GettingPriceData(coin)
	if len(x) != len(y) {
		return 0
	}
	n := float64(len(x))
	if n == 0 {
		return 0
	}

	var sumX, sumY, sumXY, sumXSquare, sumYSquare float64
	for key, valX := range x {
		valY := y[key]
		sumX += valX
		sumY += valY
		sumXY += valX * valY
		sumXSquare += valX * valX
		sumYSquare += valY * valY
	}

	numerator := n*sumXY - sumX*sumY
	denominator := math.Sqrt((n*sumXSquare - sumX*sumX) * (n*sumYSquare - sumY*sumY))
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}
