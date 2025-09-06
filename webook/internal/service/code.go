package service

import (
	"context"
	"fmt"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository"
	"github.com/czh0913/gocode/basic-go/webook/internal/service/sms"
	"math/rand"
)

const codeTplId = "123456789"

var (
	ErrCodeSendTooMany = repository.ErrCodeSendTooMany
)

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

/*
在 main 函数里面依赖注入
repo := repository.NewRedisCodeRepository(redisClient)
smsSvc := tencent.NewService(tencentClient, appId, signName)
codeSvc := service.NewCodeService(repo, smsSvc)


*/

func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {

		return err
	}
	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)

	return err
}

func (svc *CodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *CodeService) generateCode() string {
	num := rand.Intn(1000000)

	return fmt.Sprintf("%6d", num)
}
