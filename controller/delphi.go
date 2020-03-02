package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/coin-factory"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/delphi/models"
	"io/ioutil"
)

type DelphiController struct{}

var (
	// Versions for different system status
	systemVersion    = 101000
	latestVersion    = 804000
	minVersionCompat = 804000
)

func (d *DelphiController) GetVersions(c *gin.Context) {
	version := models.VersionResponse{
		LatestVersion: latestVersion,
		MinVersion:    minVersionCompat,
		SystemVersion: systemVersion,
	}
	responses.GlobalResponseError(version, nil, c)
	return
}

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
	coinsResp := models.CoinsResponse{CoinsAvailable: len(availableCoinsTags), CoinsTickers: availableCoinsTags}
	responses.GlobalResponseError(coinsResp, nil, c)
	return
}

func (d *DelphiController) GetDevCoins(c *gin.Context) {
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
	for _, coins := range availableCoins {
		availableCoinsTags = append(availableCoinsTags, coins.Info.Tag)
	}
	coinsResp := models.CoinsResponse{CoinsAvailable: len(availableCoinsTags), CoinsTickers: availableCoinsTags}
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
