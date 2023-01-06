package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/spf13/viper"
)

func main() {

	topicName := "close_order"
	//topicName := "purchase_order"

	v := viper.New()

	v.SetConfigFile("mq.json")

	err := v.ReadInConfig()
	if err != nil {
		panic("读取配置文件失败:" + err.Error())
	}

	type OrderIds struct {
		Ids []int `json:"ids"`
	}

	var orderIds OrderIds

	//fmt.Println(content) //字符串 - yaml
	//想要将一个json字符串转换成struct，需要去设置这个struct的tag
	err = v.Unmarshal(&orderIds)
	if err != nil {
		panic("解析配置失败:" + err.Error())
	}

	p, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.1.40:9876"})), // 127.0.0.1:10007 192.168.1.40:9876
		producer.WithRetry(2),
	)

	err = p.Start()

	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	for _, id := range orderIds.Ids {
		msg := &primitive.Message{
			Topic: topicName,
			Body:  []byte(strconv.Itoa(id)),
		}

		res, err := p.SendSync(context.Background(), msg)

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
	}

	err = p.Shutdown()

	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}
