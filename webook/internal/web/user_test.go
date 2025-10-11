package web

import (
	"bytes"
	"errors"
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
	"github.com/czh0913/gocode/basic-go/webook/internal/service"
	svcmocks "github.com/czh0913/gocode/basic-go/webook/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncrypt(t *testing.T) {
	password := "hello#world123"
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(encrypted, []byte(password))
	assert.NoError(t, err)
}

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.UserService

		reqBody string

		wantCode int

		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "czh@qq.com",
					Password: "00000000@user",
				}).Return(nil)
				return usersvc
			},

			reqBody: `{"email": "czh@qq.com","password": "00000000@user", "confirmPassword": "00000000@user"}`,

			wantCode: http.StatusOK,

			wantBody: "注册成功",
		},
		{
			name: "参数错误, Bind 失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)

				return usersvc
			},

			reqBody: `{"email": "czh@qq.com","password": "00000000@user", "confirmPassword": "}`,

			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},

			reqBody: `{"email": "czh.com","password": "00000000@user", "confirmPassword": "00000000@user"}`,

			wantCode: http.StatusOK,

			wantBody: "非法邮箱格式",
		},
		{
			name: "两次输入密码不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)

				return usersvc
			},

			reqBody: `{"email": "czh@qq.com","password": "00000000@user", "confirmPassword": "000000500@user"}`,

			wantCode: http.StatusOK,

			wantBody: "两次输入密码不对",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)

				return usersvc
			},

			reqBody: `{"email": "czh@qq.com","password": "00000000user", "confirmPassword": "00000000user"}`,

			wantCode: http.StatusOK,

			wantBody: "密码必须包含字母、数字、特殊字符，并且不少于八位",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "czh@qq.com",
					Password: "00000000@user",
				}).Return(service.ErrDuplicateEmail)
				return usersvc
			},

			reqBody: `{"email": "czh@qq.com","password": "00000000@user", "confirmPassword": "00000000@user"}`,

			wantCode: http.StatusOK,

			wantBody: "邮箱冲突，请换一个",
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "czh@qq.com",
					Password: "00000000@user",
				}).Return(errors.New("error "))
				return usersvc
			},

			reqBody: `{"email": "czh@qq.com","password": "00000000@user", "confirmPassword": "00000000@user"}`,

			wantCode: http.StatusOK,

			wantBody: "系统错误",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(T *testing.T) {
			// 创建一个控制器，检测mock使用情况
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			h := NewUserHandler(tc.mock(ctrl), nil, nil)
			h.RegisterRoutes(server)
			// 新建请求
			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			// 设置请求格式是json
			req.Header.Set("Content-Type", "application/json")

			require.NoError(t, err)
			//
			resp := httptest.NewRecorder()

			// ServerHTTP http请求进入gin框架的入口
			// 这样调用GIN会处理这个请求
			// 响应写回到 resp 里面
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())

		})
	}

}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersvc := svcmocks.NewMockUserService(ctrl)

	usersvc.EXPECT().Signup(gomock.Any(), gomock.Any()).
		Return(errors.New("错误"))

}
