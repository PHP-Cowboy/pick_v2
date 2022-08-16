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
}

var (
	DB           *gorm.DB
	Redis        *redis.Client
	ServerConfig = &config.ServerConfig{}
	SugarLogger  *zap.SugaredLogger
	Job          = make(chan *PrintCh, 1000)
	JobMap       = make(map[string]chan *PrintCh, 0)
)
