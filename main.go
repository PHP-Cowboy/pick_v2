package main

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"os"
	"os/signal"
	"pick_v2/handler"
	"syscall"

	"pick_v2/global"
	"pick_v2/initialize"
)

func main() {

	initialize.InitNewLogger()

	initialize.InitConfig()

	initialize.InitMysql()

	initialize.InitJob()

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

	rlog.SetLogLevel("error")

	c, mqErr := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{serverConfig.RocketMQ}),
		consumer.WithGroupName("purchase"),
	)

	if mqErr != nil {
		panic("MQ失败:" + mqErr.Error())
	}

	if err := c.Subscribe("purchase_order", consumer.MessageSelector{}, handler.Order); err != nil {
		global.Logger["err"].Infof("消费topic：purchase_order失败:%s", err.Error())
	}

	if err := c.Subscribe("close_order", consumer.MessageSelector{}, handler.NewCloseOrder); err != nil {
		global.Logger["err"].Infof("消费topic：close_order失败:%s", err.Error())
	}

	_ = c.Start()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	_ = c.Shutdown()
}
