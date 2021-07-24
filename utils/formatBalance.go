package utils

import (
	"math"
	"math/big"
)

// BalanceAtFromWei 格式化余额,返回格式: ether
func BalanceAtFromWei(balance *big.Int) *big.Float {
	fBalance := new(big.Float)
	fBalance.SetString(balance.String())
	return new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(18)))
}
