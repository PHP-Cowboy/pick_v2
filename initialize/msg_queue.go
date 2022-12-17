package initialize

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"pick_v2/global"
	"pick_v2/handler"
)

func InitMsgQueue(rocketMQ string) (c rocketmq.PushConsumer, mqErr error) {
	rlog.SetLogLevel("error")

	c, mqErr = rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{rocketMQ}),
		consumer.WithGroupName("purchase"),
	)

	if mqErr != nil {
		return
	}

	if err := c.Subscribe("purchase_order", consumer.MessageSelector{}, handler.Order); err != nil {
		global.Logger["err"].Infof("消费topic：purchase_order失败:%s", err.Error())
	}

	if err := c.Subscribe("close_order", consumer.MessageSelector{}, handler.NewCloseOrder); err != nil {
		global.Logger["err"].Infof("消费topic：close_order失败:%s", err.Error())
	}

	return
}
