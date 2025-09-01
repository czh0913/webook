package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/users/login" || ctx.Request.URL.Path == "/users/signup" {
			return
		}
		sess := sessions.Default(ctx)
		if sess == nil {
			//没有登录过
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		id := sess.Get("userId")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 20,
		})

		now := time.Now()
		// 说明还没有刷新，第一次登录
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
		// updateTime 是有的
		updateTimeVal, ok := updateTime.(time.Time)

		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if now.Sub(updateTimeVal) > time.Second*10 {
			sess.Set("update_time", now)
			sess.Save()

		}

	}
}
