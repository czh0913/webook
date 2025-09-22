package auth

import (
	"context"
	"errors"
	"github.com/czh0913/gocode/basic-go/webook/internal/service/sms"
	"github.com/golang-jwt/jwt/v5"
)

type SMSService struct {
	svc sms.Service
	key string
}

func NewSMSService() sms.Service {
	return SMSService{}
}

// Send 发送，其中 biz 必须是线下申请的代表业务方的token

func (s SMSService) Send(ctx context.Context, biz string, args []string, number ...string) error {
	var claims Claims

	token, err := jwt.ParseWithClaims(biz, &claims, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})

	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("token 不合法")
	}

	return s.svc.Send(ctx, claims.Tpl, args, number...)
}

type Claims struct {
	jwt.RegisteredClaims
	Tpl string
}
