package repository

import (
	"context"
	"database/sql"
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/cache"
	cachemocks "github.com/czh0913/gocode/basic-go/webook/internal/repository/cache/mocks"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/dao"
	daomocks "github.com/czh0913/gocode/basic-go/webook/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCacheUserRepository_FindById(t *testing.T) {
	now := time.Now()
	now = time.UnixMilli(now.UnixMilli())
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		Id       int64
		ctx      context.Context
		wantUser domain.User

		wantErr error
	}{
		{
			name: "缓存未命中，查找成功",
			Id:   1,
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userdao := daomocks.NewMockUserDAO(ctrl)
				usercache := cachemocks.NewMockUserCache(ctrl)

				usercache.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, cache.ErrRedisNotFound)
				userdao.EXPECT().FindById(gomock.Any(), int64(1)).Return(dao.User{
					Id: 1,
					Email: sql.NullString{
						String: "user@qq.com",
						Valid:  true,
					},
					Password: "xxx",
					Phone: sql.NullString{
						String: "123456789",
						Valid:  true,
					},
					Ctime: now.UnixMilli(),
					Utime: now.UnixMilli(),
				}, nil)
				usercache.EXPECT().Set(gomock.Any(), domain.User{
					Id:       1,
					Email:    "user@qq.com",
					Password: "xxx",
					Phone:    "123456789",
					Ctime:    now,
				}).Return(nil)
				return userdao, usercache
			},
			ctx: context.Background(),
			wantUser: domain.User{
				Id:       1,
				Email:    "user@qq.com",
				Password: "xxx",
				Phone:    "123456789",
				Ctime:    now,
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := NewUserRepository(tc.mock(ctrl))
			user, err := repo.FindById(tc.ctx, tc.Id)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
			time.Sleep(time.Second)
		})
	}
}
