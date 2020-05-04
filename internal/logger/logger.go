package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func init() {
	Log, _ = zap.NewProduction()
	defer Log.Sync()
	Log.Info("started logging service ...")
}
