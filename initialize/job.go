package initialize

import (
	"pick_v2/dao"
)

func InitJob() {
	go dao.YongYouConsumer()
}
