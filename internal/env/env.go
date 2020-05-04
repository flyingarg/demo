package env

import (
	"demo/internal/logger"
	"github.com/spf13/viper"
)

var runtimeViper *viper.Viper

func LoadEnv() {
	runtimeViper = viper.New()
	logger.Log.Info("loading environmental variables")

	runtimeViper.SetEnvPrefix("DEMOAPP")
	runtimeViper.AutomaticEnv()
}

func GetEnv(varName string) interface{} {
	return runtimeViper.Get(varName)
}
