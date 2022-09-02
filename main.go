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

	initialize.InitSqlServer()

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

	//c, mqErr := rocketmq.NewPushConsumer(
	//	consumer.WithNameServer([]string{"192.168.1.40:9876"}),
	//	consumer.WithGroupName("purchase"),
	//)
	//
	//if mqErr != nil {
	//	panic("MQ失败:" + mqErr.Error())
	//}
	//
	//if err := c.Subscribe("purchase_order", consumer.MessageSelector{}, handler.Order); err != nil {
	//	global.SugarLogger.Errorf("消费topic：purchase_order失败:%s", err.Error())
	//}
	//
	//_ = c.Start()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	//_ = c.Shutdown()
}
