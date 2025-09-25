package middleware

import (
	"encoding/gob"
	ijwt "github.com/czh0913/gocode/basic-go/webook/internal/web/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(ijwt ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: ijwt,
	}
}

func (l LoginJWTMiddlewareBuilder) JwtBuild() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/users/login" ||
			ctx.Request.URL.Path == "/users/signup" ||
			ctx.Request.URL.Path == "/users/refresh_token" ||
			ctx.Request.URL.Path == "/users/login_sms/code/send" ||
			ctx.Request.URL.Path == "/oauth2/wechat/authurl" ||
			ctx.Request.URL.Path == "/oauth2/wechat/callback" ||
			ctx.Request.URL.Path == "/users/login_sms" {
			ctx.Next()
			return
		}

		// 我现在用 JWT 来校验

		tokenStr := l.ExtractToken(ctx)
		claims := &ijwt.UserClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK"), nil
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
		err = l.CheckSession(ctx, claims.Ssid)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		ctx.Set("claims", claims)
	}
}
