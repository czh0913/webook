//go:build wireinject

package integration

import (
	"github.com/czh0913/gocode/basic-go/webook/internal/repository"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/cache"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/dao"
	"github.com/czh0913/gocode/basic-go/webook/internal/service"
	"github.com/czh0913/gocode/basic-go/webook/internal/web"
	"github.com/czh0913/gocode/basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,

		dao.NewUserDAO,

		cache.NewUserCache,
		cache.NewCacheCode,

		repository.NewCodeRepository,
		repository.NewUserRepository,

		service.NewCodeService,
		service.NewUserService,

		ioc.InitSMSService,

		web.NewUserHandler,

		//gin.Default,
		ioc.InitGin,
		ioc.InitMiddlewwares,
	)
	return new(gin.Engine)
}
