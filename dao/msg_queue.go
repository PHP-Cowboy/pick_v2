package dao

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"pick_v2/global"
)

// 推送批次信息到消息队列
func SendMsgQueue(topic string, messages []string) error {
	p, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{global.ServerConfig.RocketMQ})),
		producer.WithRetry(2),
	)

	err := p.Start()

	if err != nil {
		global.Logger["err"].Infof("start producer error: %s", err.Error())
		return err
	}

	mq := make([]*primitive.Message, 0, len(messages))

	for _, m := range messages {

		msg := &primitive.Message{
			Topic: topic,
			Body:  []byte(m),
		}

		mq = append(mq, msg)
	}

	res, err := p.SendSync(context.Background(), mq...)

	if err != nil {
		global.Logger["err"].Infof("send message error: %s", err.Error())
		return err
	} else {
		global.Logger["info"].Infof("send message success: result=%s", res.String())
	}

	err = p.Shutdown()

	if err != nil {
		global.Logger["err"].Infof("shutdown producer error: %s", err.Error())
		return err
	}

	return nil
}

func CloseOrderResult(occId, typ int) (err error) {

	type Ret struct {
		OccId int `json:"occ_id"`
		Typ   int `json:"typ"`
	}

	ret := Ret{
		OccId: occId,
		Typ:   typ,
	}

	var retJson []byte

	retJson, err = json.Marshal(ret)

	if err != nil {
		return
	}

	err = SendMsgQueue("close_order_result", []string{string(retJson)})

	if err != nil {
		return
	}

	return
}
