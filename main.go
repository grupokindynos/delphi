package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/delphi/controller"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	App := GetApp()
	err := App.Run(":" + port)
	if err != nil {
		panic(err)
	}
}

// GetApp is used to wrap all the additions to the GIN API.
func GetApp() *gin.Engine {
	App := gin.Default()
	App.Use(cors.Default())
	ApplyRoutes(App)
	return App
}

func ApplyRoutes(r *gin.Engine) {
	api := r.Group("/")
	{
		rate := limiter.Rate{
			Period: 1 * time.Hour,
			Limit:  1000,
		}
		store := memory.NewStore()
		limiterMiddleware := mgin.NewMiddleware(limiter.New(store, rate))
		api.Use(limiterMiddleware)
		delphiCtrl := controller.DelphiController{}
		api.GET("version", delphiCtrl.GetVersions)

		// v1 Routes
		api.POST("coins", delphiCtrl.GetCoins)
		api.POST("list", delphiCtrl.GetCoinsList)
		api.POST("dev/coins", delphiCtrl.GetCoinsDev)
		api.POST("dev/list", delphiCtrl.GetCoinsListDev)

		// v2 Route for PP >= v8.4.0
		api.POST("v2/coins", delphiCtrl.GetCoinsV2)
		api.POST("v2/dev/coins", delphiCtrl.GetDevCoinsV2)
		api.GET("v2/coin/:tag", delphiCtrl.GetCoinInfoV2)

		// v3 Route for PP >= v8.5.0 (Upgrade coin struct)
		api.GET("v3/coin/:tag", delphiCtrl.GetCoinInfoV3)

		// returns the latest commit hash of delphi on git
		api.GET("coins/version", delphiCtrl.GetCoinsVersion)
	}
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})
}
