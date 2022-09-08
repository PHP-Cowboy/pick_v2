package handler

import (
	"pick_v2/global"
	"time"
)

// 打印 写入
func AddPrintJobMap(warehouseCode string, printCh *global.PrintCh) {
	_, ok := global.PrintMapCh[warehouseCode]

	if !ok {
		ch := make(chan *global.PrintCh, 1000)
		global.PrintMapCh[warehouseCode] = ch
	}

	global.PrintMapCh[warehouseCode] <- printCh
}

// 打印 读取消费
func GetPrintJobMap(warehouseCode string) *global.PrintCh {
	job, ok := global.PrintMapCh[warehouseCode]

	if !ok {
		return nil
	}

	select {
	case printCh := <-job:
		return printCh
	default:
		return nil
	}
}

// u8 生成者
func YongYouProducer(id int) {
	global.YongYouCh <- id
}

// u8 消费者
func YongYouConsumer() {
	for {
		select {
		case id := <-global.YongYouCh:
			PushYongYou(id)
			time.Sleep(3 * time.Second)
		case <-time.After(30 * time.Second):
		}
	}
}
