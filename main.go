package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"pick_v2/global"
	"pick_v2/initialize"
)

func main() {

	initialize.InitLogger()

	initialize.InitConfig()

	initialize.InitMysql()

	initialize.InitRedis()

	g := initialize.InitRouter()

	serverConfig := global.ServerConfig

	fmt.Println("服务启动中,端口:", serverConfig.Port)

	go func() {
		err := g.Run(fmt.Sprintf(":%d", serverConfig.Port))
		if err != nil {
			panic("启动失败:" + err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
