package cryptopackage

import (
	"fmt"
	"math"

	"github.com/docobro/dimploma_project/internal/entity"
)

func (r *Repository) GettingCloseOpenValue(coins []string) ([]entity.MinuteData, error) {
	limit := 1440
	minuteDataset := make(map[string][]entity.Coin, len(coins))
	for i := 0; i < len(coins); i++ {
		res, err := r.GetOneMinuteData(coins[i], "USD", limit)
		if err != nil {
			return nil, err
		}
		minuteDataset[coins[i]] = res
	}

	minuteReq := make([]entity.MinuteData, 0)
	for i := 0; i < len(coins); i++ {
		if len(minuteDataset[coins[i]]) == 0 {
			return nil, fmt.Errorf("no data available for coin: %s", coins[i])
		}
		coinData := minuteDataset[coins[i]]
		for j := 0; j < len(coinData); j++ {
			data := coinData[j]
			res := entity.MinuteData{
				CryptoName: coins[i],
				Open:       data.OpenMinute,
				Close:      data.CloseMinute,
				High:       data.HighMinute,
				Low:        data.LowMinute,
			}
			minuteReq = append(minuteReq, res)
		}
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

// Функция для вычисления диапазонов цен
func calculatePriceRanges(data []entity.MinuteData) map[string][]float64 {
	priceRanges := make(map[string][]float64)
	for _, entry := range data {
		rangeValue := entry.High - entry.Low
		priceRanges[entry.CryptoName] = append(priceRanges[entry.CryptoName], rangeValue)
	}
	return priceRanges
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

// Функция для вычисления среднего значения
func mean(data []float64) float64 {
	sum := 0.0
	for _, value := range data {
		sum += value
	}
	return sum / float64(len(data))
}

// Функция для расчета волатильности
func calculateVolatility(data []entity.MinuteData) map[string]float64 {
	percentageChanges := calculatePercentageChanges(data)
	priceRanges := calculatePriceRanges(data)
	volatility := make(map[string]float64)

	for cryptoName, changes := range percentageChanges {
		rangeValues := priceRanges[cryptoName]

		// Расчет стандартного отклонения процентных изменений
		percentageVolatility := standardDeviation(changes)

		// Расчет стандартного отклонения диапазонов цен
		rangeVolatility := standardDeviation(rangeValues)

		// Объединение волатильности с процентными изменениями и диапазонами цен
		combinedVolatility := (percentageVolatility + rangeVolatility) / 2

		volatility[cryptoName] = combinedVolatility
	}
	return volatility
}

// Возвращение волатильности для заданных криптовалют
func (r *Repository) ReturnVolatility(coins []string) map[string]float64 {
	minuteData, err := r.GettingCloseOpenValue(coins)
	if err != nil {
		// Обработка ошибки
		return nil
	}
	volatilityMap := calculateVolatility(minuteData)
	return volatilityMap
}
