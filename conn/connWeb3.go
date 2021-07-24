package conn

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

var (
	addr = "ws://127.0.0.1:8546"
)

func Web3Conn() *ethclient.Client {
	client, err := ethclient.Dial(addr)
	if err != nil {
		log.Fatalln("connect error", err)
	}
	return client
}
