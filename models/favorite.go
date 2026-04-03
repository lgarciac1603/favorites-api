package models

type Favorite struct {
	ID         int    `json:"id"`
	UserID     int    `json:"userId"`
	CryptoID   string `json:"cryptoId"`
	CryptoName string `json:"cryptoName"`
	CreatedAt  string `json:"createdAt"`
}