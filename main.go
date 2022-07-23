package main

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"pick_v2/global"
	"pick_v2/initialize"
	"syscall"
)

func main() {

	initialize.InitLogger()

	initialize.InitConfig()

	initialize.InitMysql()

	g := initialize.InitRouter()

	serverConfig := global.ServerConfig

	zap.S().Info("服务启动中,端口:", serverConfig.Port)

	go func() {
		err := g.Run(fmt.Sprintf(":%d", serverConfig.Port))
		if err != nil {
			zap.S().Panicf("启动失败:", err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
