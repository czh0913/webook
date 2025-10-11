package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
	"github.com/czh0913/gocode/basic-go/webook/internal/service"
	svcmocks "github.com/czh0913/gocode/basic-go/webook/internal/service/mocks"
	ijwt "github.com/czh0913/gocode/basic-go/webook/internal/web/jwt"
	"github.com/czh0913/gocode/basic-go/webook/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.ArticleService

		reqBody string

		wantCode int

		wantRes Result
	}{
		{
			name: "新建并发表",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				articleSvc := svcmocks.NewMockArticleService(ctrl)
				// 期待进入 ArticleService 的 Publish
				articleSvc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "标题",
					Content: "内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)

				return articleSvc
			},
			reqBody:  `{"title": "标题","content":"内容"}`,
			wantCode: 200,

			wantRes: Result{
				Data: float64(1),
				Msg:  "OK",
			},
		},
		{
			name: "publish 失败",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				articleSvc := svcmocks.NewMockArticleService(ctrl)
				// 期待进入 ArticleService 的 Publish
				articleSvc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "标题",
					Content: "内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("publish error"))

				return articleSvc
			},
			reqBody:  `{"title": "标题","content":"内容"}`,
			wantCode: 200,

			wantRes: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		// 修改已有的帖子，然后发表成功
		// Bind 错误
		// 找不到 user
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(T *testing.T) {
			// 创建一个控制器，检测mock使用情况
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			// 新建发表需要有登录态

			server.Use(func(ctx *gin.Context) {
				ctx.Set("claims", &ijwt.UserClaims{
					Uid: 123,
				})
			})

			h := NewArticleHandler(tc.mock(ctrl), &logger.NopLogger{})
			h.RegisterRoutes(server)
			// 新建请求
			req, err := http.NewRequest(http.MethodPost, "/articles/publish", bytes.NewBuffer([]byte(tc.reqBody)))
			// 设置请求格式是json
			req.Header.Set("Content-Type", "application/json")

			require.NoError(t, err)
			//
			resp := httptest.NewRecorder()

			// ServerHTTP http请求进入gin框架的入口
			// 这样调用GIN会处理这个请求
			// 响应写回到 resp 里面
			server.ServeHTTP(resp, req)
			code := resp.Code

			// 反序列化为结果
			var result Result
			err = json.Unmarshal(resp.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, code)
			assert.Equal(t, tc.wantRes, result)

		})
	}
}
