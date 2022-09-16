package initialize

import (
	"fmt"
	"io"
	"log"
	"os"
	"pick_v2/utils/timeutil"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"pick_v2/global"
)

func InitMysql() {
	info := global.ServerConfig.MysqlInfo

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		info.User,
		info.Password,
		info.Host,
		info.Port,
		info.Name,
	)

	fileName := fmt.Sprintf("logs/%s.sql", time.Now().Format(timeutil.DateNumberFormat))

	file, err := os.Create(fileName)
	if err != nil {
		// Handle error
		panic(err)
	}
	// Make sure file is closed before your app shuts down.

	multiOutput := io.MultiWriter(os.Stdout, file)

	multiLogger := log.New(multiOutput, "", log.LstdFlags)

	logger := logger2.New(
		multiLogger,
		logger2.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			Colorful:      true,        //禁用彩色打印
			LogLevel:      logger2.Info,
		},
	)

	//var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}
}
