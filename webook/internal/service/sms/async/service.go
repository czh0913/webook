package async

import (
	"context"
	"errors"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository"
	"github.com/czh0913/gocode/basic-go/webook/internal/service/sms"
	"sync"
	"sync/atomic"
)

var (
	ErrAsyncNotFound = repository.ErrAsyncNotFound
)
var wg sync.WaitGroup

type AsyncSMSService struct {
	smsSvc     sms.Service
	asyncRepo  repository.AsyncSMSRepository
	cnt        int
	midTimeCnt int
}

func NewAsyncSMSService(svc sms.Service, async repository.AsyncSMSRepository, cnt int, midTimeCnt int) sms.Service {

	return &AsyncSMSService{
		smsSvc:     svc,
		asyncRepo:  async,
		cnt:        cnt,
		midTimeCnt: midTimeCnt,
	}
}

func (a AsyncSMSService) storeReq(ctx context.Context, biz string, args []string, number ...string) error {
	for _, val := range number {
		err := a.asyncRepo.Store(ctx, biz, args, val, "Pending")
		if err != nil {
			return err
		}
	}

	return nil
}

func (a AsyncSMSService) startReq(ctx context.Context) error {
	defer wg.Done()
	asynces, err := a.asyncRepo.Find(ctx)
	if err != nil {
		return err
	}
	if len(asynces) == 0 {
		return nil
	}

	for _, val := range asynces {
		err := a.smsSvc.Send(ctx, val.Biz, val.Args, val.Number)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a AsyncSMSService) Send(ctx context.Context, biz string, args []string, number ...string) error {
	err := a.smsSvc.Send(ctx, biz, args, number...)
	if err == nil {
		return nil
	}

	// 触发限流或者服务商错误
	err = a.storeReq(ctx, biz, args, number...)
	if err != nil {
		return err
	}

	var success atomic.Bool
	var wg sync.WaitGroup

	for idx := 0; idx < a.cnt; idx++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := a.startReq(ctx); err == nil {
				success.Store(true)
			}
		}()
	}

	wg.Wait()
	if success.Load() {
		return nil
	}
	return errors.New("all retries failed")
}
