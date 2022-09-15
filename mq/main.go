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
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.1.40:9876"})),
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

	toBeShipped := []int{1440, 1608, 1649, 11505, 12165, 13044, 13120, 15523, 15703, 15848, 16022, 16545, 17949, 18046, 18818, 18875, 19724, 19944, 20001, 20003, 20265, 22008, 24409, 24457, 24648, 24977, 26090, 26676, 26960, 26969, 27137, 27151, 27194, 27200, 27249, 27318, 27320, 27324, 27329, 27337, 27340, 27348, 27361, 27367, 27376, 27381, 27383, 27386, 27400, 27407, 27414, 27415, 27420, 27443, 27444, 27459, 27460, 27467, 27474, 27479, 27480, 27487, 27491, 27493, 27495, 27501, 27502, 27504, 27505, 27508, 27511, 27512, 27516, 27517, 27519, 27520, 27521, 27522, 27523, 27530, 27532, 27535, 27536, 27537, 27538, 27541, 27542, 27547, 27548, 27549, 27550}

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
