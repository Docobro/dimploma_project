package usecase

import (
	"log"

	"github.com/docobro/dimploma_project/internal/entity"
)

func (uc *Usecase) ParsePrices() error {
	log.Println("do parse prices")
	prices := []string{"ETH", "BTC", "TON"}
	currencies, err := uc.GetCurrencies(prices)
	if err != nil {
		return err
	}
	pricesReq := []entity.Coin{}
	for _, v := range currencies {
		pricesReq = append(pricesReq, *v)
	}
	return uc.storage.CreatePrices(pricesReq)
}
