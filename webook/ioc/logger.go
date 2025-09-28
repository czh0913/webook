package ioc

import (
	"github.com/czh0913/gocode/basic-go/webook/pkg/logger"
	"go.uber.org/zap"
)

func InitLogger() logger.Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return logger.NewZapLogger(l)
}
