package entity

type Currency struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	MaxSupply  int    `json:"max_supply"`
	Desciption string `json:"desciption"`
}
