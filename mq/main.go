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
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:10007"})),
		producer.WithRetry(2),
	)

	err := p.Start()

	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	topic := "purchase_order"

	toBeShipped := []int{
		27595, 28352, 28176, 28283, 28239, 28268, 28139, 28150, 28270, 28292,
		28192, 14902, 23055, 28351, 28353, 28354, 28362, 28123, 28125, 28128,
		27647, 28242, 28360, 27579, 28306, 28146, 28279, 25493, 28300, 28314,
		18418, 23261, 22123, 28074, 26584, 27631, 28309,
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
