package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/coin-factory"
	"github.com/grupokindynos/common/coin-factory/coins"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/delphi/models"
	"io/ioutil"
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
			// TODO enable Onion
			if coin.Info.Tag == "ONION" {
				continue
			}
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
			// TODO enable Onion
			if coin.Info.Tag == "ONION" {
				continue
			}
			// Filter tokens by network
			if coin.Info.Token {
				if coin.Info.TokenNetwork == "ethereum" {
					// TODO remove this to enable all ERC20 tokens
					if coin.Info.Tag == "ETH" {
						matchCoins = append(matchCoins, coin.Info)
					}
					// 	matchCoins = append(matchCoins, coin.Info)
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
		// Here we do the filtering

		// All version lower than 804000 must enforce update.
		if BodyRequest.Version < minVersionCompat {
			responses.GlobalResponseError(nil, errors.New("version is not compatible, need to update"), c)
			return
		}

		// Version 804000 is the minimum version for this new API system, includes all coins expect ONION. ERC20 are experimental but probable compatible.
		if BodyRequest.Version >= 804000 {
			if coin.Info.Token && coin.Info.Tag == "ETH" {
				availableCoinsTags = append(availableCoinsTags, coin.Info.Tag)
			}
			if !coin.Info.Token && coin.Info.Tag != "ONION" {
				availableCoinsTags = append(availableCoinsTags, coin.Info.Tag)
			}
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

func (d *DelphiController) GetCoinInfo(c *gin.Context) {
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
