package handler

import "pick_v2/global"

func AddPrintJob(ch *global.PrintCh) {
	global.Job <- ch
}

func GetPrintJob() *global.PrintCh {

	select {
	case printCh := <-global.Job:
		return printCh
	default:
		return nil
	}
}
