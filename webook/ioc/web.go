package ioc

import (
	"context"
	"github.com/czh0913/gocode/basic-go/webook/internal/web"
	ijwt "github.com/czh0913/gocode/basic-go/webook/internal/web/jwt"
	"github.com/czh0913/gocode/basic-go/webook/internal/web/middleware"
	"github.com/czh0913/gocode/basic-go/webook/internal/web/middleware/logger"
	logger2 "github.com/czh0913/gocode/basic-go/webook/pkg/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func InitJWTHandler(cmd redis.Cmdable) ijwt.Handler {
	return ijwt.NewJWTHandler(cmd)
}

func InitGin(mdls []gin.HandlerFunc, hdl *web.UserHandler, oauth2WeChatHdl *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterRoutes(server)
	oauth2WeChatHdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisCllient redis.Cmdable, l logger2.Logger, handler ijwt.Handler) []gin.HandlerFunc {
	bd := logger.NewBuilder(func(ctx context.Context, al *logger.AccessLog) {
		l.Debug("请求信息", logger2.Field{
			Key:   "al",
			Value: al,
		})
	}).AllowRespBody().AllowReqBody(true)
	viper.OnConfigChange(func(in fsnotify.Event) {
		bl := viper.GetBool("web.logreq")
		bd.AllowReqBody(bl)
	})

	return []gin.HandlerFunc{
		corsHdl(),
		bd.Build(),
		middleware.NewLoginJWTMiddlewareBuilder(handler).JwtBuild(),

		//ratelimit.NewBuilder(redisCllient, time.Second, 100).Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:3000"}, // 前端地址
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"}, // 如果你有自定义 header
		// 允许带 cookie
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				//开发环境
				return true
			}
			if strings.HasSuffix(origin, ".webook.com") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
