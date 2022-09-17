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
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"47.114.60.210:10007"})),
		producer.WithRetry(2),
	)

	err := p.Start()

	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	topic := "purchase_order"

	//for i := 276; i < 291; i++ {
	//
	//	msg := &primitive.Message{
	//		Topic: topic,
	//		Body:  []byte(strconv.Itoa(i)),
	//	}
	//
	//	res, err := p.SendSync(context.Background(), msg)
	//
	//	if err != nil {
	//		fmt.Printf("send message error: %s\n", err)
	//	} else {
	//		fmt.Printf("send message success: result=%s\n", res.String())
	//	}
	//}

	toBeShipped := []int{27932, 28045, 28091, 28134, 26969, 27375, 28034, 7749, 18223, 27936, 28158, 14548, 15353, 27829, 27673, 28016, 28096, 27642, 28127, 28059, 28060, 27216, 27942, 28180, 27892, 28005, 28135, 26407, 28109}

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
