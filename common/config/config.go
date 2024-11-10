package config

import (
	"github.com/spf13/viper"
)

func Init(path string) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func IsDebug() bool {
	env := viper.GetString("global.env")
	return env == "debug"
}
