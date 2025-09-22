package retelimit

import (
	"context"
	"fmt"
	"github.com/czh0913/gocode/basic-go/webook/internal/service/sms"
	"github.com/czh0913/gocode/basic-go/webook/pkg/ratelimit"
)

/*
 组合装饰器的	优点：可以只实现 需要实现的方法


*/

type RatelimitSMSServiceV1 struct {
	sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSServiceV1(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSMSServiceV1{
		Service: svc,
		limiter: limiter,
	}
}

func (s *RatelimitSMSServiceV1) Send(ctx context.Context, biz string, args []string, number ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		// 大概率 Redis 崩掉， 考虑限流还是不限流
		// 如果下游比较面对大流量比较弱，可以考虑限流
		// 如果下游比较强可以不用限流, 或者业务可用性要求很高
		// 包一下error
		return fmt.Errorf("短信服务是否出现问题，%w", err)
	}

	if limited {
		return errLimited
	}
	err = s.Service.Send(ctx, biz, args, number...)
	return err
}
