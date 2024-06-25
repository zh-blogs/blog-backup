package zlog

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var L *zap.Logger
var Z *zap.SugaredLogger

func Init() {
	if viper.GetBool("debug") {
		L, _ = zap.NewDevelopment()
	} else {
		L, _ = zap.NewProduction()
	}

	Z = L.Sugar()
}
