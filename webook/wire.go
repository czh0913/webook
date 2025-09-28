//go:build wireinject

package main

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

		// 基础设施层（外部依赖初始化）
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitSMSService,
		ioc.InitWeChatService,
		ioc.InitGin,
		ioc.InitMiddlewares,
		ioc.InitJWTHandler,
		ioc.InitLogger,
		// DAO 层（数据库访问对象）
		dao.NewUserDAO,

		// 缓存层
		cache.NewUserCache,
		cache.NewCacheCode,

		// Repository 层（聚合 DAO + Cache）
		repository.NewCodeRepository,
		repository.NewUserRepository,

		// Service 层（业务逻辑）
		service.NewCodeService,
		service.NewUserService,

		// Web 层（接口 Handler）
		web.NewOAuth2WeChatHandler,
		web.NewUserHandler,
	)
	return new(gin.Engine)
}
