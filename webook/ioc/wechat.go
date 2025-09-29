package ioc

import (
	"github.com/czh0913/gocode/basic-go/webook/internal/service/oath2/wechat"
	"github.com/czh0913/gocode/basic-go/webook/pkg/logger"
)

func InitWeChatService(logger logger.Logger) wechat.Service {
	appId := "key"
	appSecret := ""
	return wechat.NewService(appId, appSecret, logger)
}
