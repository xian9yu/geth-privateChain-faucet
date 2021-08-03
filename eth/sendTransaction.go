package eth

import (
	"crypto/ecdsa"
	"errors"
	"faucet/conn"
	"faucet/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

func SendTransaction(c *gin.Context) {

	if c.Request.Method == "GET" {
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
		//c.HTML(http.StatusOK, "index.html", nil)

	} else if c.Request.Method == "POST" {
		// 获取收款地址
		toAddress := c.PostForm("wallet")
		gasFeeStr := c.PostForm("gasFee")
		gasFee, err := strconv.Atoi(gasFeeStr)
		if err != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"message": "gas 费不正确",
			})
			return
		}
		// 检查接水地址是否有效
		if !utils.CheckAddress(toAddress) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"message": "接水地址不合法",
			})
			return
		}
		// 判断当天是否已接水
		exists, err := conn.StrExists(toAddress)
		if err != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"message": err.Error(),
			})
			return
		}
		if exists {
			account := common.HexToAddress(faucetFrom)
			balance, _ := client.BalanceAt(ctx, account, nil)
			c.HTML(http.StatusOK, "index.html", gin.H{
				"message": "今天已经领过了, 当前余额: " + utils.BalanceAtFromWei(balance).String() + " ether",
			})
		} else {

			txHEX, err := faucet(toAddress, uint64(gasFee)) // 发送交易
			if err != nil {
				c.HTML(http.StatusOK, "index.html", gin.H{
					"message": err.Error(),
				})
				return
			}

			// 交易成功把地址写入缓存
			err = conn.StrSet(toAddress, "faucet", time.Second*time.Duration(utils.GetTodaySurplusSecond()))
			if err != nil {
				log.Println(err, "===============")
			}

			c.HTML(http.StatusOK, "index.html", gin.H{
				"message": "交易成功 hash: " + txHEX,
			})
		}
	}
}

// 发送交易
func faucet(toAddress string, gasFee uint64) (string, error) {
	// 水龙头私钥
	privateKey, err := crypto.HexToECDSA(faucetPrivateKey)
	if err != nil {
		return "", errors.New("获取 faucet 账户失败: " + err.Error())
	}

	// 验证加密方式
	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	// 读取账户交易次数
	nonce, err := client.PendingNonceAt(ctx, crypto.PubkeyToAddress(*publicKeyECDSA))
	if err != nil {
		return "", err
	}

	// go-ethereum客户端提供 SuggestGasPrice 函数，用于根据'x'个先前块来获得平均燃气价格
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	// 接收者
	to := common.HexToAddress(toAddress)
	// 生成未签名以太坊事务
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    big.NewInt(0).Mul(big.NewInt(onceCoin), big.NewInt(1e18)), // 设置将要转移的 ether 数量
		Gas:      gasFee,                                                    // 设置 gas费上限
		GasPrice: gasPrice,
		Data:     []byte("faucet"),
	})

	// 使用发件人的私钥对事务进行签名
	// SignTx方法，接受一个未签名的事务和上面构造的私钥
	// SignTx方法需要EIP155签名者，需要先从客户端拿到链ID
	networkId, err := client.NetworkID(ctx) // get NetworkID
	if err != nil {
		return "", err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(networkId), privateKey)
	if err != nil {
		return "", err
	}

	// 调用“SendTransaction”来将已签名的事务广播到整个网络
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}
