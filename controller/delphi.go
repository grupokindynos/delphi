package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/coin-factory"
	"github.com/grupokindynos/common/coin-factory/coins"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/delphi/models"
	"io/ioutil"
)

type DelphiController struct{}

var (
	firstVersionCompat = 802010
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
