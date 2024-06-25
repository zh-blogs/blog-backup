package config

import (
	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	_ = viper.BindEnv("debug", "DEBUG")
}
