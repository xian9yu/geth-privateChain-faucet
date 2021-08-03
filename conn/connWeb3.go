package conn

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

var ethAddr = "ws://127.0.0.1:8546"

func Web3Conn() *ethclient.Client {
	client, err := ethclient.Dial(ethAddr)
	if err != nil {
		log.Fatalln("connect error", err)
	}
	return client
}
