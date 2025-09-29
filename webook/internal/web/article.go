package web

import (
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
	"github.com/czh0913/gocode/basic-go/webook/internal/service"
	ijwt "github.com/czh0913/gocode/basic-go/webook/internal/web/jwt"
	"github.com/czh0913/gocode/basic-go/webook/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc service.ArticleService
	l   logger.Logger
}

func NewArticleHandler(svc service.ArticleService, l logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", h.Edit)
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// ?????
	c := ctx.MustGet("claims")
	// 断言 UserClaims 指针
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		//ctx.AbortWithStatus(http.StatusUnauthorized)
		h.l.Error("未发现用户 sessino 信息")
		return
	}
	id, err := h.svc.Save(ctx, domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: claims.Uid,
		},
	})

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("保存帖子失败", logger.Error(err))
		// 日志
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})
}
