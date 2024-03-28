package main

import (
	"context"
	"fmt"
	"log"

	"github.com/docobro/dimploma_project/cryptopackage"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func main() {
	conn, err := connect()
	if err != nil {
		panic((err))
	}
	// Список криптовалют для запроса цены
	coins := []string{"BTC", "ETH", "BNB", "SOL", "XRP", "DOGE", "ADA", "AVAX"}

	// Получение текущей цены для всех криптовалют
	prices, err := cryptopackage.GetCurrentPrices(coins)
	if err != nil {
		fmt.Printf("Ошибка при получении цен: %v\n", err)
		return
	}

	// Вывод текущих цен для каждой криптовалюты
	fmt.Println("Текущие цены криптовалют:")
	for _, coin := range coins {
		price := prices.Data[coin]
		fmt.Printf("%s: %.2f USD\n", coin, price)
	}
	ctx := context.Background()
	rows, err := conn.Query(ctx, "SELECT name,toString(uuid) as uuid_str FROM system.tables LIMIT 5")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var name, uuid string
		if err := rows.Scan(
			&name,
			&uuid,
		); err != nil {
			log.Fatal(err)
		}
		log.Printf("name: %s, uuid: %s",
			name, uuid)
	}
}

func connect() (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{"localhost:8124"},
			Auth: clickhouse.Auth{
				Database: "CryptoWallet",
				Username: "default",
				Password: "",
			},
			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "an-example-go-client", Version: "0.1"},
				},
			},

			Debugf: func(format string, v ...interface{}) {
				fmt.Printf(format, v)
			},
		})
	)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}
	return conn, nil
}
