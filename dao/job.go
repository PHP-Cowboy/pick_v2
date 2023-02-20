package dao

import (
	"pick_v2/global"
	"time"
)

var BaseNum = 3

// 打印 写入
func AddPrintJobMap(warehouseCode string, typ int, printCh *global.PrintCh) {
	_, ok := global.PrintMapCh[warehouseCode][typ]

	if !ok {
		printMapCh, printMapChOk := global.PrintMapCh[warehouseCode]

		if !printMapChOk {
			ch := make(chan *global.PrintCh, 1000)
			ch <- printCh
			global.PrintMapCh[warehouseCode][typ] = ch
			return
		}

		printMapCh[typ] <- printCh

		global.PrintMapCh[warehouseCode] = printMapCh
		return
	}

	global.PrintMapCh[warehouseCode][typ] <- printCh
}

// 打印 读取消费
func GetPrintJobMap(warehouseCode string, typ int) *global.PrintCh {
	job, ok := global.PrintMapCh[warehouseCode][typ]

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

// u8 生产者
func YongYouProducer(id int) {
	global.YongYouCh <- id
}

// u8 消费者
func YongYouConsumer() {
	for {
		select {
		case id := <-global.YongYouCh:
			time.Sleep(time.Duration(BaseNum) * time.Second)
			PushYongYou(id)
		case <-time.After(time.Duration(BaseNum) * 10 * time.Second):
		}
	}
}
