package main

import (
	"github.com/czh0913/gocode/basic-go/webook/config"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/cache"
	"github.com/czh0913/gocode/basic-go/webook/internal/repository/dao"
	"github.com/czh0913/gocode/basic-go/webook/internal/service"
	"github.com/czh0913/gocode/basic-go/webook/internal/service/memory"
	"github.com/czh0913/gocode/basic-go/webook/internal/web"
	"github.com/czh0913/gocode/basic-go/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func main() {

	server := initWebServer()
	db := initDB()
	rdb := initRedis()
	u := initUser(db, rdb)
	u.RegisterRoutes(server)

	server.Run("0.0.0.0:8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	// 1. 加 CORS 中间件
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello")
	})
	server.GET("/users/login", func(c *gin.Context) {
		c.String(http.StatusMethodNotAllowed, "请用 POST 方法访问登录接口")
	})
	//redisCllient := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr, // Redis 地址
	//})
	//
	//server.Use(ratelimit.NewBuilder(redisCllient, time.Second, 100).Build())

	server.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:3000"}, // 前端地址
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"x-jwt-token"}, // 如果你有自定义 header
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
	}))
	// store 存储的地方
	//store := cookie.NewStore([]byte("secret"))
	store := memstore.NewStore([]byte("kGokUbI4xPzYsQ33OFmtV3tQ66MypaN0"), []byte("QIdjVJJwb9C1gMwpzJDTiUNotTxjBkdk"))

	// 使用 Redis 存储会话
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "", "",
	//	[]byte("kGokUbI4xPzYsQ33OFmtV3tQ66MypaN0"),
	//	[]byte("QIdjVJJwb9C1gMwpzJDTiUNotTxjBkdk"))
	//
	//if err != nil {
	//	panic(err)
	//}

	//store := memstore.NewStore([]byte("kGokUbI4xPzYsQ33OFmtV3tQ66MypaN0"),
	//	[]byte("QIdjVJJwb9C1gMwpzJDTiUNotTxjBkdk"))

	server.Use(sessions.Sessions("mysession", store))
	//登录校验
	//server.Use(middleware.NewLoginMiddlewareBuilder().Build())
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().JwtBuild())
	return server
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		//只会在初始化过程种 panic
		//整个 goroutine 结束
		//一旦初始化过程出错，应用就不要启动了
		panic(err)
		//直接退出
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}

func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	uc := cache.NewUserCache(rdb)
	repo := repository.NewUserRepository(ud, uc)
	svc := service.NewUserService(repo)
	cacheCode := cache.NewCacheCode(rdb)
	codeRepo := repository.NewCodeRepository(cacheCode)
	smsSvc := memory.NewSMSService()
	codeSvc := service.NewCodeService(codeRepo, smsSvc)
	u := web.NewUserHandler(svc, codeSvc)

	return u
}

func initRedis() redis.Cmdable {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr, // 从配置里取地址
	})
	return rdb
}
