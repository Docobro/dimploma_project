package cryptopackage

import (
	//"fmt"
	"math"

	"github.com/docobro/dimploma_project/internal/entity"
)

func (r *Repository) GettingCloseOpenValue(coins []string) ([]entity.MinuteData, error) {
	minuteDataset := make(map[string]entity.Coin, len(coins))
	for i := 0; i < len(coins); i++ {
		res, err := r.GetOneMinuteData(coins[i], "USD", 50)
		if err != nil {
			return nil, err
		}
		minuteDataset[coins[i]] = res[coins[i]]
	}

	minuteReq := []entity.MinuteData{}
	for i := 0; i < len(coins); i++ {
		res := entity.MinuteData{
			CryptoName: coins[i],
			Open:       minuteDataset[coins[i]].OpenMinute,
			Close:      minuteDataset[coins[i]].CloseMinute,
		}
		minuteReq = append(minuteReq, res)
	}

	return minuteReq, nil
}

// Функция для вычисления процентных изменений
func calculatePercentageChanges(data []entity.MinuteData) map[string][]float64 {
	percentageChanges := make(map[string][]float64)
	for _, entry := range data {
		change := ((entry.Close - entry.Open) / entry.Open) * 100
		percentageChanges[entry.CryptoName] = append(percentageChanges[entry.CryptoName], change)
	}
	return percentageChanges
}

// Функция для вычисления среднего значения
func mean(data []float64) float64 {
	sum := 0.0
	for _, value := range data {
		sum += value
	}
	return sum / float64(len(data))
}

// Функция для вычисления стандартного отклонения
func standardDeviation(data []float64) float64 {
	meanValue := mean(data)
	sumOfSquares := 0.0
	for _, value := range data {
		sumOfSquares += math.Pow(value-meanValue, 2)
	}
	variance := sumOfSquares / float64(len(data))
	return math.Sqrt(variance)
}

func calculateVolatility(data []entity.MinuteData) map[string]float64 {
	percentageChanges := calculatePercentageChanges(data)
	volatility := make(map[string]float64)
	for cryptoName, changes := range percentageChanges {
		volatility[cryptoName] = standardDeviation(changes)
	}
	return volatility
}

func (r *Repository) ReturnVolatility(coins []string) map[string]float64 {
	minuteData, err := r.GettingCloseOpenValue(coins)
	if err != nil {
		return nil
	}
	volatilityMap := calculateVolatility(minuteData)
	return volatilityMap
}
