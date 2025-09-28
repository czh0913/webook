package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"time"
)

func main() {
	initViperV1()
	initLogger()
	server := InitWebServer()

	server.Run("0.0.0.0:8080")
}

func initViperV1() {
	cfile := pflag.String("config", "config/dev.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
	})

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	zap.L().Info("zap 日志初始化成功！")
	println("")

}

func initViper() {

	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	//otherviper := viper.New()
	//otherviper.SetConfigName("myjson")
	//otherviper.AddConfigPath("./config")
	//otherviper.SetConfigType("json")

}

func initRemoteViper() {
	err := viper.AddRemoteProvider("etcd3", "127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			err = viper.WatchRemoteConfig()
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(time.Second)

		}

	}()

}
