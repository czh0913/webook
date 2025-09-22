package repository

import (
	"context"
	"database/sql"
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/cache"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/dao"
	"time"
)

var (
	ErrDuplicate    = dao.ErrDuplicate
	ErrUserNotFound = dao.ErrRecordNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByWeChat(ctx context.Context, openID string) (domain.User, error)
}

type CacheUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, ca cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: ca,
	}
}

func (repo *CacheUserRepository) FindByWeChat(ctx context.Context, openID string) (domain.User, error) {
	u, err := repo.dao.FindByWeChat(ctx, openID)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(u), nil
}

func (repo *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	user := repo.domainToEntity(u)
	return repo.dao.Insert(ctx, user)
}

func (repo *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(u), nil
}
func (repo *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(u), nil
}

func (repo *CacheUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	//先从cache找 再从dao找找到了回写cache
	u, err := repo.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}
	//没有找到

	us, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = repo.entityToDomain(us)

	//if err != nil {
	//
	//}
	go func() {
		err = repo.cache.Set(ctx, u)

	}()

	return u, nil
}

func (repo *CacheUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		WeChatInfo: domain.WeChatInfo{
			UnionID: u.WeChatUnionID.String,
			OpenID:  u.WeChatOpenID.String,
		},
		Ctime: time.UnixMilli(u.Ctime),
	}
}

func (repo *CacheUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		WeChatOpenID: sql.NullString{
			String: u.WeChatInfo.OpenID,
			Valid:  u.WeChatInfo.OpenID != "",
		},
		WeChatUnionID: sql.NullString{
			String: u.WeChatInfo.UnionID,
			Valid:  u.WeChatInfo.UnionID != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
	}
}
