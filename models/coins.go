package models

import (
	"github.com/grupokindynos/common/coin-factory/coins"
)

type CoinsResponse struct {
	Coins int `json:"coins"`
}

type CoinLegacy struct {
	Info CoinInfoLegacy `json:"info"`
}

type CoinInfoLegacy struct {
	Icon         string                           `json:"icon"`
	Tag          string                           `json:"tag"`
	Name         string                           `json:"name"`
	Trezor       bool                             `json:"trezor"`
	Ledger       bool                             `json:"ledger"`
	Segwit       bool                             `json:"segwit"`
	Masternodes  bool                             `json:"masternodes"`
	Token        bool                             `json:"token"`
	StableCoin   bool                             `json:"stable_coin"`
	TokenNetwork string                           `json:"token_network,omitempty"`
	Contract     string                           `json:"contract,omitempty"`
	Decimals     int                              `json:"decimals,omitempty"`
	Blockbook    string                           `json:"blockbook"`
	Protocol     string                           `json:"protocol"`
	TxVersion    int                              `json:"tx_version"`
	TxBuilder    string                           `json:"tx_builder"`
	HDIndex      int                              `json:"hd_index"`
	Networks     map[string]CoinNetworkInfoLegacy `json:"networks"`
}

type CoinNetworkInfoLegacy struct {
	MessagePrefix string                     `json:"messagePrefix"`
	Bech32        string                     `json:"bech32"`
	Bip32         coins.CoinNetWorkBip32Info `json:"bip32"`
	PubKeyHash    int                        `json:"pubKeyHash"`
	ScriptHash    int                        `json:"scriptHash"`
	Wif           int                        `json:"wif"`
}

type CoinsResponseV2 struct {
	CoinsAvailable int      `json:"coins_available"`
	CoinsTickers   []string `json:"coins_tickers"`
}

type CoinsRequestBody struct {
	Version int `json:"version"`
}
