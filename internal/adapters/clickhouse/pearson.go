package clickhouse

import (
	"context"
	"fmt"
	"math"
	"strconv"
)

func (r *Repository) GettingDataFromDB(coin string, nameValue string, nameBase string, limit int) map[int]float64 {
	limits := strconv.Itoa(limit)
	tokens, err := r.GetCryptoTokens()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil
	}
	id := tokens[coin].ID
	indexQuery := fmt.Sprintf(`SELECT CAST(%s AS Float64) FROM cryptowallet.%s WHERE crypto_id = '%s' ORDER BY created_at DESC LIMIT %s`, nameValue, nameBase, id, limits)
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
	name_value_1, name_value_2 := "trade_volume", "value"
	nameBase_1, nameBase_2 := "trade_volume_1m", "prices"
	limit_1, limit_2 := 60, 360
	x = r.GettingDataFromDB(coin, name_value_1, nameBase_1, limit_1)
	y = r.GettingDataFromDB(coin, name_value_2, nameBase_2, limit_2)

	averagedY := make(map[int]float64, 60)
	for i := 0; i < 60; i++ {
		var sum float64
		for j := 0; j < 6; j++ {
			sum += y[i*6+j]
		}
		average := sum / 6
		averagedY[i] = average
	}

	if len(x) != len(averagedY) {
		return 0
	}
	n := float64(len(x))
	if n == 0 {
		return 0
	}

	var sumX, sumY, sumXY, sumXSquare, sumYSquare float64
	for key, valX := range x {
		valY := averagedY[key]
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

func (r *Repository) PearsonPriceToVolatilityCorrelation(coin string) float64 {
	var x, y map[int]float64
	name_value_1, name_value_2 := "volatility", "value"
	nameBase_1, nameBase_2 := "volatilities", "prices"
	limit_1, limit_2 := 60, 360
	x = r.GettingDataFromDB(coin, name_value_1, nameBase_1, limit_1)
	y = r.GettingDataFromDB(coin, name_value_2, nameBase_2, limit_2)

	averagedY := make(map[int]float64, 60)
	for i := 0; i < 60; i++ {
		var sum float64
		for j := 0; j < 6; j++ {
			sum += y[i*6+j]
		}
		average := sum / 6
		averagedY[i] = average
	}

	if len(x) != len(averagedY) {
		return 0
	}
	n := float64(len(x))
	if n == 0 {
		return 0
	}

	var sumX, sumY, sumXY, sumXSquare, sumYSquare float64
	for key, valX := range x {
		valY := averagedY[key]
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

func (r *Repository) PearsonPriceToMrktCapCorrelation(coin string) float64 {
	var x, y map[int]float64
	name_value_1, name_value_2 := "market_cap", "value"
	nameBase_1, nameBase_2 := "prices", "prices"
	limit_1, limit_2 := 360, 360
	x = r.GettingDataFromDB(coin, name_value_1, nameBase_1, limit_1)
	y = r.GettingDataFromDB(coin, name_value_2, nameBase_2, limit_2)

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
