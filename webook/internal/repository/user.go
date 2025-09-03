package repository

import (
	"context"
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/cache"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, ca *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: ca,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (repo *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
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
	u = repo.toDomain(us)

	go func() {
		err = repo.cache.Set(ctx, u)
		if err != nil {
			
		}

	}()

	return u, err
}
