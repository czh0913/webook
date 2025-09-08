package middleware

import (
	"encoding/gob"
	"github.com/czh0913/gocode/basic-go/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l LoginJWTMiddlewareBuilder) JwtBuild() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/users/login" ||
			ctx.Request.URL.Path == "/users/signup" ||
			ctx.Request.URL.Path == "/users/login_sms/code/send" ||
			ctx.Request.URL.Path == "/users/login_sms" {
			ctx.Next()
			return
		}

		// 我现在用 JWT 来校验

		tokenHeader := ctx.GetHeader("Authorization")

		if tokenHeader == "" {
			// 没有登录过
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			// 没登陆
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		tokenStr := segs[1]
		// ParseWithClaims 里面会修改把解析了之后的数据赋值给 claims ，会修改claims 需要传指针

		claims := &web.UserClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})

		if err != nil {
			// 没登陆
			ctx.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		if !token.Valid || token == nil || claims.Uid == 0 {
			// 没登陆
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 在确认 calims 存在之后
		if claims.UserAgent != ctx.Request.UserAgent() {
			// 用户代理不匹配，可能是 CSRF 攻击
			// 需要加监控
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user agent"})
			return
		}

		now := time.Now()

		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Minute))
			tokenStr, err = token.SignedString(web.JWTKey)
			if err != nil {
				// 记录日志 生成token 失败
				log.Println(" jwt 续约失败", err)
			}
			//println("续约 JWT token")
			ctx.Header("x-jwt-token", tokenStr)
		}

		ctx.Set("claims", claims)
	}
}
