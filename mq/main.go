package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"os"
	"strconv"
)

func main() {
	p, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:10007"})), // 127.0.0.1:10007 192.168.1.40:9876
		producer.WithRetry(2),
	)

	err := p.Start()

	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	topic := "purchase_order"

	toBeShipped := []int{
		2599, 3056, 3140, 3396, 5648, 5709, 6877, 9869, 11505, 12754, 14387, 20637, 23786, 25181,
		26948, 27855, 27856, 27901, 27923, 28253, 28374, 28387, 28401, 28418, 28439, 28445, 28471,
		28581, 28595, 28788, 28809, 28826, 28837, 28838, 28849, 28853, 28859, 28867, 28874, 28876,
		28886, 28896, 28897,
	}
	for _, v := range toBeShipped {
		msg := &primitive.Message{
			Topic: topic,
			Body:  []byte(strconv.Itoa(v)),
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
