package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

var Op *viper.Viper

func init() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	Op = viper.New()
	Op.AddConfigPath(path)   //设置读取的文件路径
	Op.SetConfigName("app")  //设置读取的文件名
	Op.SetConfigType("yaml") //设置文件的类型
	Op.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println(e)
	})
	if err := Op.ReadInConfig(); err != nil {
		fmt.Println(err.Error(), "获取配置文件失败")
	}
}
