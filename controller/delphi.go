package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/coin-factory"
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
	responses.GlobalResponseError(coinfactory.Coins, nil, c)
	return
}
