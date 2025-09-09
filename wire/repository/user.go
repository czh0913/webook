package repository

import "github.com/czh0913/gocode/basic-go/wire/repository/dao"

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(userDAO *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: userDAO,
	}
}
