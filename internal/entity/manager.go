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
		PriceIndex  int32
		VolumeIndex float32
		Coin        *Coin
		CoinName    string
		Price       int64
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
	res.Coin = b.Coin
	res.CoinName = b.CoinName
	res.Price = b.Price
	res.PriceIndex = b.PriceIndex
	res.VolumeIndex = b.VolumeIndex
	return res
}
