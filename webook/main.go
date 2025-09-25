package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	initViperV1()
	server := InitWebServer()

	server.Run("0.0.0.0:8080")
}

func initViperV1() {
	cfile := pflag.String("config", "config/config.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

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
