package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"
	"io"
	"time"
)

type MiddleWareBuilder struct {
	allowReqBody  *atomic.Bool
	allowRespBody bool
	loggerFun     func(ctx context.Context, al *AccessLog)
}

func NewBuilder(loggerFun func(ctx context.Context, al *AccessLog)) *MiddleWareBuilder {
	return &MiddleWareBuilder{
		loggerFun:    loggerFun,
		allowReqBody: atomic.NewBool(false),
	}
}

func (m *MiddleWareBuilder) Build() gin.HandlerFunc {
	start := time.Now()
	return func(ctx *gin.Context) {
		url := ctx.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		// 设置 Method URL
		al := &AccessLog{
			Method: ctx.Request.Method,
			Url:    url,
		}

		//判断 reqBody
		if ctx.Request.Body != nil && m.allowReqBody.Load() {
			// 由于 Request.Body 的类型是 io.ReadCloser 使用io包的读操作一般会用一次就没了，
			// 也就是读出来了，Body里面没内容了
			body, _ := io.ReadAll(ctx.Request.Body)
			// 我们需要放回 Body 去
			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
			if len(body) > 1024 {
				body = body[:1024]
			}
			al.ReqBody = string(body)

		}
		if m.allowRespBody {
			ctx.Writer = reponseWriter{
				al:             al,
				ResponseWriter: ctx.Writer,
			}
		}

		defer func() {
			al.Duration = time.Since(start).String()

			m.loggerFun(ctx, al)
		}()

		ctx.Next()

	}
}

func (m *MiddleWareBuilder) AllowReqBody(ok bool) *MiddleWareBuilder {
	m.allowReqBody.Store(ok)
	return m
}

func (m *MiddleWareBuilder) AllowRespBody() *MiddleWareBuilder {
	m.allowRespBody = true
	return m
}

type reponseWriter struct {
	al                 *AccessLog //反向引用 ，这里为啥用 AccessLog ？
	gin.ResponseWriter            // 装饰器模式 ？ 为啥要用装饰器
}

func (r reponseWriter) WriteHeader(statusCode int) {
	r.al.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r reponseWriter) Write(data []byte) (int, error) {
	r.al.RespBody = string(data)
	return r.ResponseWriter.Write(data)
}

func (r reponseWriter) WriteString(data string) (int, error) {
	r.al.RespBody = data
	return r.ResponseWriter.WriteString(data)
}

type AccessLog struct {
	Method     string
	Url        string
	Duration   string
	ReqBody    string
	RespBody   string
	statusCode int
}

// 新建结构体，里面要有存储拦截信息的字段，然后组合需要用到的一组接口，然后需要的方法（装饰器模式），
// 然后在使用时替换自己定义的结构体
