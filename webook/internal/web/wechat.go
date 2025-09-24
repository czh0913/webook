package web

import (
	"errors"
	"fmt"
	"github.com/czh0913/gocode/basic-go/webook/internal/service"
	"github.com/czh0913/gocode/basic-go/webook/internal/service/oath2/wechat"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"time"

	"net/http"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	jwtHandler
	stateKey []byte
}

func NewOAuth2WeChatHandler(svc wechat.Service, userSvc service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:      svc,
		userSvc:  userSvc,
		stateKey: []byte("RtyCrTBkkTS2U6XCawU7kmnWNPaup4nf"),
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("oauth2/wechat") // 更新路由路径以匹配规范
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New()
	url, err := h.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造扫码登录 URL 失败",
		})
		return
	}

	if err := h.setToken(ctx, state); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}

func (h *OAuth2WechatHandler) setToken(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, StateClaims{
		state: state,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	})

	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		return err
	}

	ctx.SetCookie("jwt-state", tokenStr, 600,
		"/oauth2/wechat/callback", "", false, true)
	return nil
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	// 实现回调逻辑
	code := ctx.Query("code")
	err := h.verifyState(ctx)

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "登录失败",
		})

		return
	}

	info, err := h.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})

		return
	}
	// 在这设置 JWT-Token
	u, err := h.userSvc.FindOrCreatByWeChat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}

	err = h.setJWTToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})

}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	tokenState, err := ctx.Cookie("jwt-state")
	if err != nil {

		return fmt.Errorf("token获取失败, %w", err)
	}

	var sc StateClaims
	token, err := jwt.ParseWithClaims(tokenState, &sc, func(token *jwt.Token) (any, error) {
		return h.stateKey, nil
	})
	if err != nil || token.Valid {
		return fmt.Errorf("解码失败, %w", err)
	}
	if sc.state != state {

		return errors.New("state 不相同")
	}
	return nil
}

type StateClaims struct {
	state string
	jwt.RegisteredClaims
}
