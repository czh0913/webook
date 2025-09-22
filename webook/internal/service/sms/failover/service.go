package failover

import (
	"context"
	"errors"
	"github.com/czh0913/gocode/basic-go/webook/internal/service/sms"
	"log"
	"sync/atomic"
)

type FailOverSMSService struct {
	svcs []sms.Service

	idx uint64
}

func NewFailOverSMSService(svcs []sms.Service) sms.Service {
	return &FailOverSMSService{
		svcs: svcs,
	}
}

// 轮询
func (f *FailOverSMSService) Send(ctx context.Context, biz string, args []string, number ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, biz, args, number...)
		if err == nil {
			return nil
		}
		log.Println(err)
		//可能超时但也发送成功，导致后续用户接收两个验证码
	}
	return errors.New("供应商全部失败")
}

func (f *FailOverSMSService) SendV1(ctx context.Context, tpl string, args []string, number ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	lenth := uint64(len(f.svcs))
	for i := idx; i < uint64(lenth)+idx; i++ {
		svc := f.svcs[int(i%lenth)]
		err := svc.Send(ctx, tpl, args, number...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			{
				return err
			}

		default:
			//输出日志
		}
	}

	return errors.New("供应商全部失败")

}
