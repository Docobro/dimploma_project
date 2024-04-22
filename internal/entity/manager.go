package entity

/*
 There is data types that used for operation with manager
*/
import "time"

type (
	Key struct {
		Timestamp int64
	}

	Value struct {
		Requests    int64
		Impressions int64

		PriceIndex  float64
		VolumeIndex float64
		Coin        *Coin
		CoinName    string
		Price       float64
	}

	Rows map[Key]Value
)

func NewKey(k Key) Key {
	k.Timestamp = time.Now().Unix()
	k.Timestamp -= k.Timestamp % 60
	return k
}

func (a Value) Assign(b Value) Value {
	res := a
	res.Requests += b.Requests
	res.Impressions += b.Impressions

	return res
}
