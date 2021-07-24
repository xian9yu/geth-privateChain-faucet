package eth

import (
	"context"
	"faucet/conn"
)

var (
	client = conn.Web3Conn()
	ctx    = context.Background()

	// 水龙头钱包地址
	faucetFrom = "0x3687868fd18edb417d434a5e7157e428a716a54d"

	// 水龙头私钥
	faucetPrivateKey = "43be8992c040d5150c8dca5f46283765c20cab15b19390814c749d3040404cdd"

	// 每次领取数量(单位: ether)
	onceCoin int64 = 2
)
