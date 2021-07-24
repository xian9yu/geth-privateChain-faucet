package eth

import (
	"faucet/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	account := common.HexToAddress(faucetFrom)

	// 获取水龙头账号余额
	balance, err := client.BalanceAt(ctx, account, nil)
	if err != nil {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"message": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"address": account,
		"balance": utils.BalanceAtFromWei(balance),
	})

}
