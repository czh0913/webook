package service

import (
	"context"
	"errors"
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository"
	repomocks "github.com/czh0913/gocode/basic-go/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func Test_userService_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) repository.UserRepository

		email string

		password string

		wantUser domain.User

		wantErr error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "user@qq.com").
					Return(domain.User{
						Email:    "user@qq.com",
						Password: "$2a$10$ejq21rxYSE4Bx6wnefRndeMNggr9ZteMmTOiFDmZJbBEZ4qcR69Ia",
						Phone:    "123456789",
						Ctime:    now,
					}, nil)

				return repo
			},
			email:    "user@qq.com",
			password: "00000000@user",
			wantUser: domain.User{
				Email:    "user@qq.com",
				Password: "$2a$10$ejq21rxYSE4Bx6wnefRndeMNggr9ZteMmTOiFDmZJbBEZ4qcR69Ia",
				Phone:    "123456789",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "user@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)

				return repo
			},
			email:    "user@qq.com",
			password: "00000000@user",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "DB 错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "user@qq.com").
					Return(domain.User{}, errors.New("DB 错误"))

				return repo
			},
			email:    "user@qq.com",
			password: "00000000@user",
			wantUser: domain.User{},
			wantErr:  errors.New("DB 错误"),
		},
		{
			name: "密码不对",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "user@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "user@qq.com",
			password: "000000@user",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userRepo := NewUserService(tc.mock(ctrl))

			u, err := userRepo.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantUser, u)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestEncrypted(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("00000000@user"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(hash))
	}

}
