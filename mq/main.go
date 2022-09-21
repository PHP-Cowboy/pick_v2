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
		18418, 22123, 23261, 25493, 26584, 27579, 27595, 28074, 28146, 28239, 28268, 28270,
		28300, 28309, 28314, 28725, 28843, 28857, 28860, 28862, 28895, 28902, 28905, 28906,
		28918, 28925, 28926, 28959, 28966, 28973, 28975, 28983, 29000, 29011, 29021, 29036,
		29041, 29054, 29060, 29065,
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
