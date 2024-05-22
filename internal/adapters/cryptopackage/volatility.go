package cryptopackage

import (
	"fmt"

	"github.com/docobro/dimploma_project/internal/entity"
)

func (r *Repository) ParseVolumeMinute() error {
	coins := []string{"BTC", "ETH"}
	minuteDataset := make(map[string]entity.Coin, len(coins))
	for i := 0; i < len(coins); i++ {
		res, err := r.GetOneMinuteData(coins[i], "USD", 50)
		if err != nil {
			return err
		}
		minuteDataset[coins[i]] = res[coins[i]]
	}

	minuteReq := []entity.MinuteData{}
	for i := 0; i < len(coins); i++ {
		res := entity.MinuteData{
			CryptoName: coins[i],
		}
		minuteReq = append(minuteReq, res)
	}
	fmt.Println(minuteReq)
	return nil
}
