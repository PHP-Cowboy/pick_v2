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

	//日志初始化
	initialize.InitNewLogger()
	//配置初始化
	initialize.InitConfig()
	//数据库初始化
	initialize.InitMysql()
	//任务初始化
	initialize.InitJob()
	//redis初始化
	initialize.InitRedis()
	//缓存初始化
	initialize.InitGoCache()
	//路由初始化
	g := initialize.InitRouter()

	serverConfig := global.ServerConfig

	fmt.Println("服务启动中,端口:", serverConfig.Port)

	go func() {
		err := g.Run(fmt.Sprintf(":%d", serverConfig.Port))
		if err != nil {
			panic("启动失败:" + err.Error())
		}
	}()

	queue, err := initialize.InitMsgQueue(serverConfig.RocketMQ)

	if err != nil {
		panic("MQ失败:" + err.Error())
	}

	_ = queue.Start()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	_ = queue.Shutdown()
}
