package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type jwtHandler struct {
	atKey []byte
	rtKey []byte
}

func newJWTHandler() jwtHandler {
	return jwtHandler{
		atKey: []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK"),
		rtKey: []byte("cAD3w4zXADS6YAR7fKJKCr6vuHprfk5n"),
	}
}

func (h jwtHandler) setLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.setJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}

	err = h.setRefreshJWTToken(ctx, uid, ssid)

	return nil
}

func (h jwtHandler) setJWTToken(ctx *gin.Context, uid int64, ssid string) error {
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
	tokenStr, err := token.SignedString(h.atKey)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)

	return nil
}

func (h jwtHandler) setRefreshJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	uc := RefreshClaims{
		uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			// 10 天过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 240)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString(h.rtKey)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)

	return nil
}

func ExtractToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Authorization")

	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

type RefreshClaims struct {
	uid  int64
	Ssid string
	jwt.RegisteredClaims
}
type UserClaims struct {
	jwt.RegisteredClaims
	Ssid      string
	Uid       int64
	UserAgent string
}
