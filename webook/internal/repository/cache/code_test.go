package cache

import (
	"context"
	"errors"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/cache/redismocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name string

		mock  func(ctrl *gomock.Controller) redis.Cmdable
		ctx   context.Context
		biz   string
		phone string
		code  string

		wantErr error
	}{
		{
			name: "设置成功",

			biz:   "验证码登录",
			phone: "123456789",
			code:  "123456",

			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdable := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())

				cmd.SetVal(int64(0))

				cmdable.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:验证码登录:123456789"},
					[]any{"123456"}).Return(cmd)

				return cmdable
			},
			ctx:     context.Background(),
			wantErr: nil,
		},
		{
			name: "redis 错误",

			biz:   "验证码登录",
			phone: "123456789",
			code:  "123456",

			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdable := redismocks.NewMockCmdable(ctrl)

				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(errors.New("mock redis 错误"))

				cmdable.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:验证码登录:123456789"},
					[]any{"123456"}).Return(cmd)

				return cmdable
			},
			ctx:     context.Background(),
			wantErr: errors.New("mock redis 错误"),
		},
		{
			name: "发送太频繁",

			biz:   "验证码登录",
			phone: "123456789",
			code:  "123456",

			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdable := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())

				cmd.SetVal(int64(-1))

				cmdable.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:验证码登录:123456789"},
					[]any{"123456"}).Return(cmd)

				return cmdable
			},
			ctx:     context.Background(),
			wantErr: ErrCodeSendTooMany,
		},
		{
			name: "系统错误",

			biz:   "验证码登录",
			phone: "123456789",
			code:  "123456",

			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdable := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())

				cmd.SetVal(int64(-2))

				cmdable.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:验证码登录:123456789"},
					[]any{"123456"}).Return(cmd)

				return cmdable
			},
			ctx:     context.Background(),
			wantErr: errors.New("系统错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := NewCacheCode(tc.mock(ctrl))
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)

			assert.Equal(t, tc.wantErr, err)
		})
	}
}
