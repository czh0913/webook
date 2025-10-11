package startup

import (
	"github.com/czh0913/gocode/basic-go/webook/pkg/logger"
)

func InitLog() logger.Logger {
	return &logger.NopLogger{}
}
