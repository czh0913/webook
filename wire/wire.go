//go:build wireinject

// 让 wire 来注入这里的代码
package wire

import (
	"github.com/czh0913/gocode/basic-go/wire/repository"
	"github.com/czh0913/gocode/basic-go/wire/repository/dao"
	"github.com/google/wire"
)

func InitRepository() *repository.UserRepository {
	//传入各个组件的初始化方法
	//只用这里传入初始化方法，wire会自动帮你初始化到wire_gen.go
	//传入Build参数顺序没影响
	wire.Build(repository.NewUserRepository,
		dao.NewUserDAO,
		InitDB)

	return new(repository.UserRepository)
}
