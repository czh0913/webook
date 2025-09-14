package ioc

import (
	"github.com/czh0913/gocode/basic-go/webook/internal/service/memory"
	"github.com/czh0913/gocode/basic-go/webook/internal/service/sms"
)

func InitSMSService() sms.Service {
	return memory.NewService()
}
