package global

import (
	"gorm.io/gorm"
	"pick_v2/config"
)

var (
	DB           *gorm.DB
	ServerConfig = &config.ServerConfig{}
)
