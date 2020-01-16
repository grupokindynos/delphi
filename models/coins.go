package models

type CoinsResponse struct {
	Coins int `json:"coins"`
}

type CoinsRequestBody struct {
	Version int `json:"version"`
}
