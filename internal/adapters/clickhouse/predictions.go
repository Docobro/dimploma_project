package clickhouse

import (
	"fmt"
	"sort"

	"gonum.org/v1/gonum/mat"
)

// Function to convert a map[int]float64 to a sorted []float64 based on the keys
func mapToSortedSlice(data map[int]float64) []float64 {
	keys := make([]int, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	sortedData := make([]float64, len(data))
	for i, k := range keys {
		sortedData[i] = data[k]
	}
	return sortedData
}

func (r *Repository) ReturnPredictions(coin string) ([]float64, float64) {
	var Aprices, volumes, volatilities, corrPriceVolume, corrPriceVolatility map[int]float64
	name_value_1, name_value_2, name_value_4 := "price_index", "trade_volume", "volatility"
	nameBase_1, nameBase_2, nameBase_4 := "indices", "trade_volume_1m", "volatilities"
	limit_1, limit_2, limit_4 := 36, 6, 6

	Aprices = r.GettingDataFromDB(coin, name_value_1, nameBase_1, limit_1)
	volumes = r.GettingDataFromDB(coin, name_value_2, nameBase_2, limit_2)
	volatilities = r.GettingDataFromDB(coin, name_value_4, nameBase_4, limit_4)

	valuePearson_1, valuePearson_3 := "priceToVolume", "priceToVolatility"
	basePearson := "pearson_correlation"
	limit_5 := 1

	corrPriceVolume = r.GettingDataFromDB(coin, valuePearson_1, basePearson, limit_5)
	corrPriceVolatility = r.GettingDataFromDB(coin, valuePearson_3, basePearson, limit_5)

	prices := make(map[int]float64, limit_2)
	for i := 0; i < limit_2; i++ {
		var sum float64
		for j := 0; j < 6; j++ {
			sum += Aprices[i*6+j]
		}
		average := sum / 6
		prices[i] = average
	}
	if len(volumes) != limit_2 {
		fmt.Println("Недостаточно данных для подсчета предиктов криптовалюты ", coin, ", необходимо дождаться ", limit_2, " записей для цены, объема и волатильности, текущее значение: ", len(volumes))
		return nil, 0
	}

	// Преобразование map в отсортированные срезы
	priceSlice := mapToSortedSlice(prices)
	volumeSlice := mapToSortedSlice(volumes)
	volatilitySlice := mapToSortedSlice(volatilities)

	// Присваивание весов для последних данных
	weights := make([]float64, limit_2)
	for i := 0; i < limit_2; i++ {
		weights[i] = 1.0 + float64(i)/float64(limit_2)
	}

	// Создание регрессионной модели с учетом весов
	beta := weightedRegressionModel(volumeSlice, volatilitySlice, priceSlice, corrPriceVolume[0], corrPriceVolatility[0], weights)

	// Прогнозирование на тех же данных, которые использовались для обучения
	features := mat.NewDense(limit_2, 2, nil)
	for i := 0; i < limit_2; i++ {
		features.Set(i, 0, volumeSlice[i])
		features.Set(i, 1, volatilitySlice[i])
	}

	predictions := predict(features, beta)

	// Оценка точности модели
	mse := evaluate(predictions, priceSlice)
	return predictions, mse
}

// Функция weightedRegressionModel строит регрессионную модель с учетом весов на основе переданных признаков и цен,
// включая коэффициенты корреляции Пирсона. Возвращает вектор коэффициентов регрессии.
func weightedRegressionModel(volumes, volatilities, prices []float64, corrPriceVolume, corrPriceVolatility float64, weights []float64) *mat.VecDense {
	// Соберем обучающие данные в матрицу признаков
	features := mat.NewDense(len(prices), 2, nil)
	for i := 0; i < len(prices); i++ {
		features.Set(i, 0, volumes[i])
		features.Set(i, 1, volatilities[i])
	}

	// Применение весов к признакам и целевым значениям
	weightedFeatures := mat.NewDense(len(prices), 2, nil)
	weightedTarget := mat.NewVecDense(len(prices), nil)
	for i := 0; i < len(prices); i++ {
		weightedFeatures.Set(i, 0, volumes[i]*weights[i])
		weightedFeatures.Set(i, 1, volatilities[i]*weights[i])
		weightedTarget.SetVec(i, prices[i]*weights[i])
	}

	// Создадим транспонированную матрицу признаков
	rows, cols := weightedFeatures.Dims()
	featuresT := mat.NewDense(cols, rows, nil)
	featuresT.Copy(weightedFeatures.T())

	// Создание матрицы-диагонали с коэффициентами корреляции
	corrMatrix := mat.NewDiagDense(2, []float64{corrPriceVolume, corrPriceVolatility})

	// Выведем значения коэффициентов корреляции Пирсона
	fmt.Println("Коэффициенты корреляции Пирсона:")
	fmt.Printf("Price-Volume: %f\n", corrPriceVolume)
	fmt.Printf("Price-Volatility: %f\n", corrPriceVolatility)

	// Вычислим (X^T * X + CorrelationMatrix)
	var xtxPlusCorr mat.Dense
	xtxPlusCorr.Mul(featuresT, weightedFeatures)
	xtxPlusCorr.Add(&xtxPlusCorr, corrMatrix)

	// Вычислим (X^T * y)
	var xty mat.VecDense
	xty.MulVec(featuresT, weightedTarget)

	// Решим систему уравнений (X^T * X + CorrelationMatrix) * beta = (X^T * y)
	var beta mat.VecDense

	err := beta.SolveVec(&xtxPlusCorr, &xty)
	if err != nil {
		fmt.Println("Ошибка при решении системы уравнений:", err)
		return nil
	}

	return &beta
}

// Функция predict используется для прогнозирования цен на основе переданных признаков и коэффициентов регрессии.
func predict(features *mat.Dense, beta *mat.VecDense) []float64 {
	rows, _ := features.Dims()
	predictions := make([]float64, rows)
	for i := 0; i < rows; i++ {
		row := features.RowView(i)
		predictions[i] = mat.Dot(row, beta)
	}
	fmt.Println("Прогнозируемые цены:", predictions)
	return predictions
}

// Функция evaluate используется для оценки точности модели.
// Рассчитывается среднеквадратичная ошибка (MSE).
func evaluate(predictions, actual []float64) float64 {
	var sumError float64
	for i := 0; i < len(predictions); i++ {
		error := predictions[i] - actual[i]
		sumError += error * error
	}
	mse := sumError / float64(len(predictions))
	fmt.Printf("Среднеквадратичная ошибка (MSE): %f\n", mse)
	return mse
}
