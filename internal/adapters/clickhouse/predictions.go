package clickhouse

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func (r *Repository) ReturnNextPrediction(coin string) float64 {
	var Aprices, volumes, volatilities, corrPriceVolume, corrPriceVolatility map[int]float64
	name_value_1, name_value_2, name_value_4 := "value", "trade_volume", "volatility"
	nameBase_1, nameBase_2, nameBase_4 := "prices", "trade_volume_1m", "volatilities"
	limit_1, limit_2, limit_4 := 360, 60, 60

	Aprices = r.GettingDataFromDB(coin, name_value_1, nameBase_1, limit_1)
	volumes = r.GettingDataFromDB(coin, name_value_2, nameBase_2, limit_2)
	volatilities = r.GettingDataFromDB(coin, name_value_4, nameBase_4, limit_4)

	valuePearson_1, valuePearson_3 := "priceToVolume", "priceToVolatility"
	basePearson := "pearson_correlation"
	limit_5 := 1

	corrPriceVolume = r.GettingDataFromDB(coin, valuePearson_1, basePearson, limit_5)
	corrPriceVolatility = r.GettingDataFromDB(coin, valuePearson_3, basePearson, limit_5)

	prices := make(map[int]float64, 60)
	for i := 0; i < 60; i++ {
		var sum float64
		for j := 0; j < 6; j++ {
			sum += Aprices[i*6+j]
		}
		average := sum / 6
		prices[i] = average
	}

	// Используем все данные для обучения модели
	trainPrices := make([]float64, len(volumes))
	trainVolumes := make([]float64, len(volumes))
	trainVolatilities := make([]float64, len(volumes))

	for i := 0; i < len(volumes); i++ {
		trainPrices[i] = prices[i]
		trainVolumes[i] = volumes[i]
		trainVolatilities[i] = volatilities[i]
	}

	// Создание регрессионной модели на всех данных
	beta := regressionModel(trainVolumes, trainVolatilities, trainPrices, corrPriceVolume[0], corrPriceVolatility[0])

	// Прогнозирование следующей цены
	nextVolume := volumes[len(volumes)-1]
	nextVolatility := volatilities[len(volatilities)-1]

	nextFeatures := mat.NewVecDense(2, []float64{nextVolume, nextVolatility})
	nextPrediction := mat.Dot(nextFeatures, beta)

	return nextPrediction
}

// Функция regressionModel строит регрессионную модель на основе переданных признаков и цен,
// включая коэффициенты корреляции Пирсона. Возвращает вектор коэффициентов регрессии.
func regressionModel(volumes, volatilities, prices []float64, corrPriceVolume, corrPriceVolatility float64) *mat.VecDense {
	// Соберем обучающие данные в матрицу признаков
	features := mat.NewDense(len(prices), 2, nil)
	for i := 0; i < len(prices); i++ {
		features.Set(i, 0, volumes[i])
		features.Set(i, 1, volatilities[i])
	}

	// Вектор обучающих цен
	target := mat.NewVecDense(len(prices), prices)

	// Создадим транспонированную матрицу признаков
	rows, cols := features.Dims()
	featuresT := mat.NewDense(cols, rows, nil)
	featuresT.Copy(features.T())

	// Создание матрицы-диагонали с коэффициентами корреляции
	corrMatrix := mat.NewDiagDense(2, []float64{corrPriceVolume, corrPriceVolatility})

	// Выведем значения коэффициентов корреляции Пирсона
	fmt.Println("Коэффициенты корреляции Пирсона:")
	fmt.Printf("Price-Volume: %f\n", corrPriceVolume)
	fmt.Printf("Price-Volatility: %f\n", corrPriceVolatility)

	// Вывод изначальных данных
	fmt.Println("Данные:")
	fmt.Printf("Price: %f\n", prices)
	fmt.Printf("Volume: %f\n", volumes)
	fmt.Printf("Volatility: %f\n", volatilities)

	// Выведем содержимое матрицы-диагонали
	fmt.Println("Матрица-диагональ с коэффициентами корреляции:")
	fmt.Println(mat.Formatted(corrMatrix))

	// Вычислим (X^T * X + CorrelationMatrix)
	var xtxPlusCorr mat.Dense
	xtxPlusCorr.Mul(featuresT, features)
	xtxPlusCorr.Add(&xtxPlusCorr, corrMatrix)

	// Выведем содержимое матрицы X^T * X + CorrelationMatrix
	fmt.Println("Матрица (X^T * X + CorrelationMatrix):")
	fmt.Println(mat.Formatted(&xtxPlusCorr, mat.Prefix(""), mat.Squeeze()))

	// Вычислим (X^T * y)
	var xty mat.VecDense
	xty.MulVec(featuresT, target)

	// Выведем содержимое матрицы (X^T * y)
	fmt.Println("Матрица (X^T * y):")
	fmt.Println(mat.Formatted(&xty, mat.Prefix(""), mat.Squeeze()))

	// Решим систему уравнений (X^T * X + CorrelationMatrix) * beta = (X^T * y)
	var beta mat.VecDense

	fmt.Println("beta: ", beta)

	err := beta.SolveVec(&xtxPlusCorr, &xty)
	if err != nil {
		fmt.Println("Ошибка при решении системы уравнений:", err)
		return nil
	}

	// Выведем коэффициенты регрессии
	fmt.Println("Коэффициенты регрессии:")
	fmt.Println(mat.Formatted(&beta, mat.Prefix(""), mat.Squeeze()))

	return &beta
}
