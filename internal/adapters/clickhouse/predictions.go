package clickhouse

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func (r *Repository) ReturnPredictions(coin string) ([]float64, float64) {
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

	fmt.Println("dsadasdas", len(volumes), len(prices), len(Aprices), len(volatilities))

	// Разделение данных на обучающую и тестовую выборки (70% на обучение, 30% на тестирование)
	trainSize := int(float64(len(volumes)) * 0.70)
	testSize := len(volumes) - trainSize

	trainPrices := make([]float64, trainSize)
	testPrices := make([]float64, testSize)

	trainVolumes := make([]float64, trainSize)
	testVolumes := make([]float64, testSize)

	trainVolatilities := make([]float64, trainSize)
	testVolatilities := make([]float64, testSize)

	// Заполнение обучающих и тестовых данных
	for i := 0; i < trainSize; i++ {
		trainPrices[i] = prices[i]
		trainVolumes[i] = volumes[i]
		trainVolatilities[i] = volatilities[i]
	}

	for i := trainSize; i < len(volumes); i++ {
		testPrices[i-trainSize] = prices[i]
		testVolumes[i-trainSize] = volumes[i]
		testVolatilities[i-trainSize] = volatilities[i]
	}

	// Создание регрессионной модели на обучающих данных
	beta := regressionModel(trainVolumes, trainVolatilities, trainPrices, corrPriceVolume[0], corrPriceVolatility[0])

	// Прогнозирование на тестовых данных
	testFeatures := mat.NewDense(testSize, 2, nil)
	for i := 0; i < testSize; i++ {
		testFeatures.Set(i, 0, testVolumes[i])
		testFeatures.Set(i, 1, testVolatilities[i])
	}

	predictions := predict(testFeatures, beta)

	// Оценка точности модели на тестовых данных
	mse := evaluate(predictions, testPrices)
	return predictions, mse
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

	// // Выведем значения коэффициентов корреляции Пирсона
	// fmt.Println("Коэффициенты корреляции Пирсона:")
	// fmt.Printf("Price-Volume: %f\n", corrPriceVolume)
	// fmt.Printf("Price-Volatility: %f\n", corrPriceVolatility)

	// // Вывод изначальных данных
	// fmt.Println("Данные:")
	// fmt.Printf("Price: %f\n", prices)
	// fmt.Printf("VOlume: %f\n", volumes)
	// fmt.Printf("Volatility: %f\n", volatilities)

	// Вычислим (X^T * X + CorrelationMatrix)
	var xtxPlusCorr mat.Dense
	xtxPlusCorr.Mul(featuresT, features)
	xtxPlusCorr.Add(&xtxPlusCorr, corrMatrix)

	// Вычислим (X^T * y)
	var xty mat.VecDense
	xty.MulVec(featuresT, target)

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

// Функция evaluate используется для оценки точности модели на тестовых данных.
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
