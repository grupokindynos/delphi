package controller

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	coinfactory "github.com/grupokindynos/common/coin-factory"
	"github.com/grupokindynos/common/coin-factory/coins"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/delphi/models"
	"gopkg.in/src-d/go-git.v3"
)

type DelphiController struct{}

var (
	// Versions for different system status
	firstVersionCompat = 802010

	systemVersion    = 100000
	latestVersion    = 802010
	minVersionCompat = 802010
)

func (d *DelphiController) GetCoins(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	var BodyRequest models.CoinsRequestBody
	err = json.Unmarshal(body, &BodyRequest)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	allCoins := coinfactory.Coins
	var matchCoins []*coins.Coin
	if BodyRequest.Version >= firstVersionCompat {
		for _, coin := range allCoins {
			// Filter tokens by network
			if coin.Info.Token {
				if coin.Info.TokenNetwork == "ethereum" {
					// TODO remove this to enable all ERC20 tokens
					if coin.Info.Tag == "ETH" {
						matchCoins = append(matchCoins, coin)
					}
					// 	matchCoins = append(matchCoins, coin)
				}
			}
			// Filter by builders
			if coin.Info.TxBuilder == "bitcoinjs" ||
				coin.Info.TxBuilder == "groestljs" {
				matchCoins = append(matchCoins, coin)
			}
		}
	}
	coinsResp := models.CoinsResponse{Coins: len(matchCoins)}
	responses.GlobalResponseError(coinsResp, nil, c)
	return
}

func (d *DelphiController) GetCoinsList(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	var BodyRequest models.CoinsRequestBody
	err = json.Unmarshal(body, &BodyRequest)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	allCoins := coinfactory.Coins
	var matchCoins []coins.CoinInfo
	if BodyRequest.Version >= firstVersionCompat {
		for _, coin := range allCoins {
			// Filter tokens by network
			if coin.Info.Token {
				if coin.Info.TokenNetwork == "ethereum" {
					if coin.Info.Tag == "ETH" {
						matchCoins = append(matchCoins, coin.Info)
					}
				}
			}
			// Filter by builders
			if coin.Info.TxBuilder == "bitcoinjs" ||
				coin.Info.TxBuilder == "groestljs" {
				matchCoins = append(matchCoins, coin.Info)
			}
		}
	}
	responses.GlobalResponseError(matchCoins, nil, c)
	return
}

func (d *DelphiController) GetVersions(c *gin.Context) {
	version := models.VersionResponse{
		LatestVersion: latestVersion,
		MinVersion:    minVersionCompat,
		SystemVersion: systemVersion,
	}
	responses.GlobalResponseError(version, nil, c)
	return
}

func (d *DelphiController) GetCoinsDev(c *gin.Context) {
	allCoins := coinfactory.Coins
	var matchCoins []*coins.Coin
	for _, coin := range allCoins {
		matchCoins = append(matchCoins, coin)
	}
	coinsResp := models.CoinsResponse{Coins: len(matchCoins)}
	responses.GlobalResponseError(coinsResp, nil, c)
	return
}

func (d *DelphiController) GetCoinsListDev(c *gin.Context) {
	allCoins := coinfactory.Coins
	var matchCoins []*coins.CoinInfo
	for _, coin := range allCoins {
		matchCoins = append(matchCoins, &coin.Info)
	}
	responses.GlobalResponseError(matchCoins, nil, c)
	return
}

func (d *DelphiController) GetCoinsV2(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	var BodyRequest models.CoinsRequestBody
	err = json.Unmarshal(body, &BodyRequest)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	availableCoins := coinfactory.Coins
	var availableCoinsTags []string
	for _, coin := range availableCoins {
		if coin.Info.Tag == "BAT" || coin.Info.Tag == "LINK" || coin.Info.Tag == "MANA" || coin.Info.Tag == "DAPS" || coin.Info.Tag == "DAI" || coin.Info.Tag == "GTH" {
			continue
		}

		if BodyRequest.Version >= 805000 {
			availableCoinsTags = append(availableCoinsTags, coin.Info.Tag)
		}

		// Version 804000 is the minimum version for this new API system, includes all coins. ERC20 are experimental but probable compatible.
		if BodyRequest.Version >= 804000 && BodyRequest.Version < 805000 {
			// This coins are never available for version below 805000 since it requires an upgrade on the Coin model
			if coin.Info.Tag == "CRW" {
				continue
			}
			availableCoinsTags = append(availableCoinsTags, coin.Info.Tag)
		}

		// All version lower than 804000 must enforce update.
		if BodyRequest.Version < minVersionCompat {
			responses.GlobalResponseError(nil, errors.New("version is not compatible, need to update"), c)
			return
		}

	}
	coinsResp := models.CoinsResponseV2{CoinsAvailable: len(availableCoinsTags), CoinsTickers: availableCoinsTags}
	responses.GlobalResponseError(coinsResp, nil, c)
	return
}

func (d *DelphiController) GetDevCoinsV2(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	var BodyRequest models.CoinsRequestBody
	err = json.Unmarshal(body, &BodyRequest)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	availableCoins := coinfactory.Coins
	var availableCoinsTags []string
	for _, c := range availableCoins {
		availableCoinsTags = append(availableCoinsTags, c.Info.Tag)
	}
	coinsResp := models.CoinsResponseV2{CoinsAvailable: len(availableCoinsTags), CoinsTickers: availableCoinsTags}
	responses.GlobalResponseError(coinsResp, nil, c)
	return
}

func (d *DelphiController) GetCoinInfoV2(c *gin.Context) {
	tag := c.Param("tag")
	if tag == "" {
		responses.GlobalResponseError(nil, errors.New("no coin tag defined"), c)
		return
	}
	coin, err := coinfactory.GetCoin(tag)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	coinLegacyFormat := models.CoinLegacy{
		Info: models.CoinInfoLegacy{
			Icon:         coin.Info.Icon,
			Tag:          coin.Info.Tag,
			Name:         coin.Info.Name,
			Trezor:       coin.Info.Trezor,
			Ledger:       coin.Info.Ledger,
			Segwit:       coin.Info.Segwit,
			Masternodes:  coin.Info.Masternodes,
			Token:        coin.Info.Token,
			StableCoin:   coin.Info.StableCoin,
			TokenNetwork: coin.Info.TokenNetwork,
			Contract:     coin.Info.Contract,
			Decimals:     coin.Info.Decimals,
			Blockbook:    coin.Info.Blockbook,
			Protocol:     coin.Info.Protocol,
			TxBuilder:    coin.Info.TxBuilder,
			TxVersion:    coin.Info.TxVersion,
			HDIndex:      coin.Info.HDIndex,
			Networks:     make(map[string]models.CoinNetworkInfoLegacy),
		},
	}

	for name, params := range coin.Info.Networks {
		coinLegacyFormat.Info.Networks[name] = models.CoinNetworkInfoLegacy{
			MessagePrefix: params.MessagePrefix,
			Bech32:        params.Bech32,
			Bip32:         params.Bip32,
			PubKeyHash:    params.PubKeyHash[0],
			ScriptHash:    params.ScriptHash[0],
			Wif:           params.Wif[0],
		}
	}

	responses.GlobalResponseError(coinLegacyFormat.Info, nil, c)
	return
}

func (d *DelphiController) GetCoinInfoV3(c *gin.Context) {
	tag := c.Param("tag")
	if tag == "" {
		responses.GlobalResponseError(nil, errors.New("no coin tag defined"), c)
		return
	}
	coin, err := coinfactory.GetCoin(tag)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	responses.GlobalResponseError(coin.Info, nil, c)
	return
}

func (d *DelphiController) GetCoinsVersion(c *gin.Context) {
	res, err := git.NewRepository("https://github.com/grupokindynos/delphi", nil)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}

	if err := res.Pull("origin", "refs/heads/master"); err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}

	head, err := res.Head("origin")
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}

	responses.GlobalResponse(head.String(), c)
}
