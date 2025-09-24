package ioc

import (
	"github.com/czh0913/gocode/basic-go/webook/internal/service/oath2/wechat"
)

func InitWeChatService() wechat.Service {
	appId := "key"
	appSecret := ""
	return wechat.NewService(appId, appSecret)
}
