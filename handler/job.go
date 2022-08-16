package handler

import "pick_v2/global"

// 单仓 写入
func AddPrintJob(ch *global.PrintCh) {
	global.Job <- ch
}

// 单仓 读取消费
func GetPrintJob() *global.PrintCh {

	select {
	case printCh := <-global.Job:
		return printCh
	default:
		return nil
	}
}

// 多仓 写入
func AddPrintJobMap(warehouseCode string, printCh *global.PrintCh) {
	_, ok := global.JobMap[warehouseCode]

	if !ok {
		ch := make(chan *global.PrintCh, 1000)
		global.JobMap[warehouseCode] = ch
	}

	global.JobMap[warehouseCode] <- printCh
}

// 多仓 读取消费
func GetPrintJobMap(warehouseCode string) *global.PrintCh {
	job, ok := global.JobMap[warehouseCode]

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
