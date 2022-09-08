package initialize

import "pick_v2/handler"

func InitJob() {
	go handler.YongYouConsumer()
}
