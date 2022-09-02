package initialize

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"pick_v2/global"
	"time"
)

func InitSqlServer() {

	//dsn := fmt.Sprintf("sqlserver://ding2:%s@8.136.191.24:9930?database=ETL_U8", "%&#c68%vM5")

	dsn := "odbc:server=8.136.191.24;user id=ding2;password=%&#c68%vM5;database=ETL_U8"

	logger := logger2.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger2.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			Colorful:      true,        //禁用彩色打印
			LogLevel:      logger2.Info,
		},
	)

	var err error

	global.SqlServer, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: logger,
	})

	if err != nil {
		return
	}

	if err != nil {
		panic(err)
	}
}
