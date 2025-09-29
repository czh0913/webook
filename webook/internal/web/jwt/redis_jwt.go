package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"
)

var (
	AtKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")
	RtKey = []byte("cAD3w4zXADS6YAR7fKJKCr6vuHprfk5n")
)

type RedisJWTHandler struct {
	cmd redis.Cmdable
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd: cmd,
	}
}
func (h *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Authorization")

	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

func (h *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}

	err = h.setRefreshJWTToken(ctx, uid, ssid)

	return nil
}

func (h *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "失效token")
	ctx.Header("x-refresh-token", "失效token")
	c, _ := ctx.Get("Claims")
	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return errors.New("")
	}

	return h.cmd.Set(ctx, claims.Ssid, "", time.Hour*10*24).Err()

}

func (h *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	val, err := h.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	switch err {
	case redis.Nil:
		return nil
	case nil:
		if val == 0 {
			return nil
		}
		return errors.New("session 已经无效了")
	default:
		return err
	}
}

func (h *RedisJWTHandler) setRefreshJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	uc := RefreshClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			// 10 天过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 240)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString(RtKey)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)

	return nil
}

func (h *RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	uc := UserClaims{
		Uid:       uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		Ssid:      ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			// 1 分钟过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)

	return nil
}
