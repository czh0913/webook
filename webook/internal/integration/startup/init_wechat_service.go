package startup

import (
	"github.com/czh0913/gocode/basic-go/webook/internal/service/oath2/wechat"
	"github.com/czh0913/gocode/basic-go/webook/pkg/logger"
)

// InitPhantomWechatService 没啥用的虚拟的 wechatService
func InitPhantomWechatService(l logger.Logger) wechat.Service {
	return wechat.NewService("", "", l)
}
