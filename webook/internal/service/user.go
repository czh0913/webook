package service

import (
	"context"
	"errors"
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicate
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
)

type UserService interface {
	Signup(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindOrCreat(ctx context.Context, phone string) (domain.User, error)
	FindOrCreatByWeChat(ctx context.Context, wechatInfo domain.WeChatInfo) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)

	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 检查密码对不对
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil

}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	u, err := svc.repo.FindById(ctx, id)
	return u, err
}

func (svc *userService) FindOrCreat(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)

	if err != nil && err != repository.ErrDuplicate {
		return u, err
	}

	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *userService) FindOrCreatByWeChat(ctx context.Context, wechatInfo domain.WeChatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWeChat(ctx, wechatInfo.OpenID)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	u = domain.User{
		WeChatInfo: wechatInfo,
	}
	err = svc.repo.Create(ctx, u)

	if err != nil && err != repository.ErrDuplicate {
		return u, err
	}

	return svc.repo.FindByWeChat(ctx, wechatInfo.OpenID)
}
