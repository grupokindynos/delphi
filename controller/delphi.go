package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/coin-factory"
	"github.com/grupokindynos/common/coin-factory/coins"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/delphi/models"
)

type DelphiController struct{}

func (d *DelphiController) GetCoins(c *gin.Context) {
	coinsResp := models.CoinsResponse{Coins:len(coinfactory.Coins)}
	responses.GlobalResponseError(coinsResp, nil, c)
	return
}

func (d *DelphiController) GetCoinsList(c *gin.Context) {
	var coinArray []coins.CoinInfo
	for _, coin := range coinfactory.Coins {
		coinArray = append(coinArray, coin.Info)
	}
	responses.GlobalResponseError(coinArray, nil, c)
	return
}
