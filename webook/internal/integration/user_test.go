package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/czh0913/gocode/basic-go/webook/internal/web"
	"github.com/czh0913/gocode/basic-go/webook/ioc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUerHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := InitWebServer()
	rdb := ioc.InitRedis()
	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		reqBody string

		wantCode int

		wantBody web.Result
	}{
		{

			name: "发送成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				//清理数据
				val, err := rdb.Get(ctx, "phone_code:login:123456789").Result()
				assert.NoError(t, err)
				assert.True(t, len(val) == 6)

				_, err = rdb.Del(ctx, "phone_code:login:123456789").Result()
				assert.NoError(t, err)
			},

			reqBody:  `{"phone" : "123456789"}`,
			wantCode: 200,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
		{

			name: "发送频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

				_, err := rdb.Set(ctx, "phone_code:login:123456789", "123456", time.Minute*9+time.Second*30).Result()
				cancel()
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				//清理数据
				val, err := rdb.Get(ctx, "phone_code:login:123456789").Result()
				assert.NoError(t, err)

				assert.Equal(t, "123456", val)
				_, err = rdb.Del(ctx, "phone_code:login:123456789").Result()
				assert.NoError(t, err)
			},

			reqBody:  `{"phone" : "123456789"}`,
			wantCode: 200,
			wantBody: web.Result{
				Msg: "发送太频繁，请稍后再试",
			},
		},
		{

			name: "系统错误",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

				_, err := rdb.Set(ctx, "phone_code:login:123456789", "123456", 0).Result()
				cancel()
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				//清理数据
				val, err := rdb.Get(ctx, "phone_code:login:123456789").Result()
				assert.NoError(t, err)

				assert.Equal(t, "123456", val)
				_, err = rdb.Del(ctx, "phone_code:login:123456789").Result()
				assert.NoError(t, err)
			},

			reqBody:  `{"phone" : "123456789"}`,
			wantCode: 200,
			wantBody: web.Result{
				Msg:  "系统错误",
				Code: 5,
			},
		},
		{

			name: "手机号码为空",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},

			reqBody:  `{"phone" : ""}`,
			wantCode: 200,
			wantBody: web.Result{
				Msg:  "输入错误",
				Code: 4,
			},
		},
		{

			name: "数据格式错误",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},

			reqBody:  `{"phone" :}`,
			wantCode: 400,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send", bytes.NewBuffer([]byte(tc.reqBody)))

			req.Header.Set("Content-Type", "application/json")

			require.NoError(t, err)
			//
			resp := httptest.NewRecorder()

			// ServerHTTP http请求进入gin框架的入口
			// 这样调用GIN会处理这个请求
			// 响应写回到 resp 里面
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)

			if resp.Code != 200 {
				return
			}

			var webRes web.Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, webRes)
			tc.after(t)
		})
	}
}
