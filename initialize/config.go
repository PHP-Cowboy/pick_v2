package initialize

import (

	"fmt"
	"pick_v2/global"

	"github.com/spf13/viper"
	"go.uber.org/zap"

)

func InitConfig() {
	v := viper.New()

	v.SetConfigFile("config.yaml")

	err := v.ReadInConfig()
	if err != nil {
		zap.S().Panicf("读取配置文件失败:", err.Error())
	}

	//fmt.Println(content) //字符串 - yaml
	//想要将一个json字符串转换成struct，需要去设置这个struct的tag
	err =  v.Unmarshal(global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取配置失败： %s", err.Error())
	}
	fmt.Println(&global.ServerConfig)
}
