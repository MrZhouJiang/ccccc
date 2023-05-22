package conf

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

//数据库配置
type DbConfig struct {
	Host   string
	User   string
	Passwd string
	DBName string
}
type Struct struct {
	Mysql DbConfig
}

var Config Struct

//初始化，优先加载
func init() {
	var appName = "app"
	//路径
	viper.AddConfigPath("./resource")
	//如果 环境信息，或者命令行有 这个变量，则追加环境信息，如果没有，默认取 app.yml
	//加载了环境则是 app-dev.yml |app-test.yml
	configEnv := os.Getenv("GO_ENV")

	if configEnv != "" {
		appName += "-" + configEnv
	}

	viper.SetConfigName(appName)
	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Error reading config file  %s\n", err)
	}
	// viper读取了配置，并且将 json字符 反序列化 到对象
	err := viper.Unmarshal(&Config)
	if err != nil {
		log.Panicf("unable to decode into struct, %v", err)
	}
}
