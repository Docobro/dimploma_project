package clickhouse

import (
	"context"
	"fmt"

	//"math"
	"time"
)

func (r *Repository) GettingData(cryptoID string) (map[time.Time]float64, error) {
	indexQuery := fmt.Sprintf("SELECT created_at, CAST(value AS Float64) FROM cryptowallet.prices WHERE crypto_id = '%s' ORDER BY created_at DESC LIMIT 360", cryptoID)

	rows, err := r.conn.Query(context.Background(), indexQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make(map[time.Time]float64)
	for rows.Next() {
		var createdAt time.Time
		var value float64
		if err := rows.Scan(&createdAt, &value); err != nil {
			return nil, err
		}
		values[createdAt] = value
	}

	return values, nil
}

// func (r *Repository) PearsonCorrelation() float64 {
// 	var x, y map[time.Time]float64
// 	x, err1 := r.GettingData("835a9c62-caaf-4891-91a7-4f9137f3f815")
// 	y, err2 := r.GettingData("43dc5ab9-3ff7-4303-ae06-9aafe0114822")
// 	if len(x) != len(y) {
// 		return 0
// 	}
// 	fmt.Println(err1, err2, x, y)

// 	n := float64(len(x))
// 	if n == 0 {
// 		return 0
// 	}

// 	var sumX, sumY, sumXY, sumXSquare, sumYSquare float64
// 	for key, valX := range x {
// 		valY, exists := y[key]
// 		if !exists {
// 			return 0
// 		}
// 		sumX += valX
// 		sumY += valY
// 		sumXY += valX * valY
// 		sumXSquare += valX * valX
// 		sumYSquare += valY * valY
// 	}

// 	numerator := n*sumXY - sumX*sumY
// 	denominator := math.Sqrt((n*sumXSquare - sumX*sumX) * (n*sumYSquare - sumY*sumY))
// 	if denominator == 0 {
// 		return 0
// 	}

// 	return numerator / denominator
// }
