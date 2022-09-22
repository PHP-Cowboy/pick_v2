package global

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"pick_v2/config"
)

type PrintCh struct {
	DeliveryOrderNo string
	ShopId          int
	Type            int // 1-全部打印 2-打印箱单 3-打印出库单 第一次全打，后边的前段选
}

var (
	DB           *gorm.DB
	SqlServer    *gorm.DB
	Redis        *redis.Client
	ServerConfig = &config.ServerConfig{}
	SugarLogger  *zap.SugaredLogger
	Logger       = make(map[string]*zap.SugaredLogger, 0)
	PrintMapCh   = make(map[string]chan *PrintCh, 0)
	YongYouCh    = make(chan int, 1000)
)
