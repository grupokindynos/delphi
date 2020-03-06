package models

type CoinsResponse struct {
	Coins int `json:"coins"`
}

type CoinsResponseV2 struct {
	CoinsAvailable int      `json:"coins_available"`
	CoinsTickers   []string `json:"coins_tickers"`
}

type CoinsRequestBody struct {
	Version int `json:"version"`
}
