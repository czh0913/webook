package failover

import (
	"context"
	"errors"
	"github.com/czh0913/gocode/basic-go/webook/internal/service/sms"
	"sync/atomic"
)

type TimeoutFailoverSMSService struct {
	svcs      []sms.Service
	idx       int32
	cnt       int32
	threshold int32
}

func NewTimeoutFailoverSMSService(svcs []sms.Service, threshold int32) sms.Service {
	return &TimeoutFailoverSMSService{
		svcs:      svcs,
		idx:       0,
		cnt:       0,
		threshold: threshold,
	}
}

// 超时次数达到阈值更换
func (t TimeoutFailoverSMSService) Send(ctx context.Context, biz string, args []string, number ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)

	if cnt > t.threshold {
		newIdx := (idx + 1) % int32(len(t.svcs))

		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = atomic.LoadInt32(&t.idx)
	}

	svc := t.svcs[idx]
	err := svc.Send(ctx, biz, args, number...)
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		atomic.AddInt32(&t.cnt, 1)
	case err == nil:
		atomic.StoreInt32(&t.cnt, 0)
	default:
	}
	return err
}
