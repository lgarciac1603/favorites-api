package models

type Favorite struct {
	ID         string `json:"id"`
	UserID     string `json:"userId"`
	CryptoId   string `json:"cryptoId"`
	CryptoName string `json:"cryptoName"`
	CreatedAt  string `json:"createdAt"`
}
