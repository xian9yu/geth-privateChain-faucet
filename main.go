package main

import (
	"faucet/conn"
	"faucet/eth"
	"log"

	"github.com/gin-gonic/gin"
)

func initRouter(router *gin.Engine) {
	router.GET("/", eth.Index)
	router.POST("/", eth.SendTransaction)

}
func main() {
	// 初始化redis
	conn.InitRedis()
	// 初始化router
	r := gin.Default()
	initRouter(r)

	// 加载views目录下的所有html模板文件
	r.LoadHTMLGlob("./views/*")

	// 使用gin自带的异常恢复中间件，避免出现异常时程序退出
	r.Use(gin.Recovery())

	err := r.Run(":3003")
	if err != nil {
		log.Fatal(err)
	}
}
