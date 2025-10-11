//go:build wireinject

package startup

import (
	"github.com/czh0913/gocode/basic-go/webook/internal/repository"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/article"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/cache"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/dao"
	"github.com/czh0913/gocode/basic-go/webook/internal/service"
	"github.com/czh0913/gocode/basic-go/webook/internal/web"
	ijwt "github.com/czh0913/gocode/basic-go/webook/internal/web/jwt"
	"github.com/czh0913/gocode/basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)
var userSvcProvider = wire.NewSet(
	dao.NewUserDAO,
	cache.NewUserCache,
	repository.NewUserRepository,
	service.NewUserService)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdProvider,
		userSvcProvider,
		//articlSvcProvider,
		cache.NewCacheCode,
		dao.NewGORMArticleDAO,
		repository.NewCodeRepository,
		article.NewCachedArticleRepository,
		// service 部分
		// 集成测试我们显式指定使用内存实现
		ioc.InitSMSService,

		// 指定啥也不干的 wechat service
		InitPhantomWechatService,
		service.NewCodeService,
		service.NewArticleService,

		// handler 部分
		web.NewUserHandler,         // 用户操作的 Handler
		web.NewOAuth2WeChatHandler, // 微信扫码的 Handler
		web.NewArticleHandler,      // 文章操作的 Handler

		ijwt.NewRedisJWTHandler,

		// gin 的中间件
		ioc.InitMiddlewares,

		// Web 服务器
		ioc.InitGin,
	)
	// 随便返回一个
	return gin.Default()
}

// 测试用的 ArticleHandler
func InitArticleHandler() *web.ArticleHandler {
	wire.Build(thirdProvider, dao.NewGORMArticleDAO, service.NewArticleService, article.NewCachedArticleRepository, web.NewArticleHandler)
	return &web.ArticleHandler{}
}

func InitUserSvc() service.UserService {
	wire.Build(thirdProvider, userSvcProvider)
	return service.NewUserService(nil, nil)
}

func InitJwtHdl() ijwt.Handler {
	wire.Build(thirdProvider, ijwt.NewRedisJWTHandler)
	return ijwt.NewRedisJWTHandler(nil)
}
