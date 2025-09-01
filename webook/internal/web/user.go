package web

import (
	"fmt"
	"github.com/czh0913/gocode/basic-go/webook/internal/domain"
	"github.com/czh0913/gocode/basic-go/webook/internal/service"
	"github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// UserHandle 定义所有user的路由
type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp2.Regexp
	passwordExp *regexp2.Regexp
}

// 新建一个UserHandler结构体
func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = `^\w+([-+.]?\w+)*@\w+([-.]?\w+)*\.\w+([-.]?\w+)*$`
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$`
	)

	emailExp := regexp2.MustCompile(emailRegexPattern, regexp2.None)
	passwordExp := regexp2.MustCompile(passwordRegexPattern, regexp2.None)

	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) RegisterRoutesV1(ug *gin.RouterGroup) {

	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)

}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit)
	//ug.GET("/profile", u.Profile)
	ug.GET("/profile", u.ProfileJWT)

}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"` // 原变量名拼写保留
		Password        string `json:"password"`
	}

	var req SignUpReq

	// 自动绑定 JSON 请求体
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusBadRequest, "请求参数错误")
		return
	}

	// 校验邮箱格式
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误1")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}

	// 校验密码一致性
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入密码不一致")
		return
	}

	// 校验密码格式
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误2")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码格式不对，必须大于8位，包含数字，特殊字符，字母")
		return
	}

	//调用一下 svc 的方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrUserDuplicate {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}

	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误3")
		return
	}

	ctx.String(http.StatusOK, "注册成功！")

}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq

	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusBadRequest, "请求参数错误: %v", err)
		return
	}

	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}

	// 在这里使用 JWT 设置登录态
	// 生成 JWT token

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)), // 设置过期时间为1分钟
		},
		Uid:       user.Id,
		UserAgent: ctx.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("kGokUbI4xPzYsQ33OFmtV3tQ66MypaN0"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误 生成 token 失败")
		return
	}

	ctx.Header("x-jwt-token", tokenStr)

	ctx.String(http.StatusOK, "登录成功")

	return
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq

	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}

	//设置session 在ctx里面挂载 session
	sess := sessions.Default(ctx)
	//可以随便设置session里面的值了
	sess.Set("userId", user.Id)

	sess.Options(sessions.Options{
		Secure:   true,
		HttpOnly: true,
		//
		MaxAge: 20,
	})
	sess.Save()
	if err != nil {
		fmt.Println("保存 session 失败:", err)
	}

	ctx.String(http.StatusOK, "登录成功")

	return
}
func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "test")
}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, ok := ctx.Get("claims")

	// 可以断定必然有 claims
	if !ok {
		// 监控这里
		//ctx.String(http.StatusOK, "系统错误 没有claims ")
		ctx.JSON(http.StatusOK, gin.H{
			"error": "系统错误 没有claims ",
		})
		return
	}
	// 断言 claims 的类型 因为 Get 方法返回的是 interface{} 类型
	claims, ok := c.(*UserClaims)
	if claims == nil {
	}
	if !ok {
		//ctx.String(http.StatusOK, "系统错误 claims 断言失败")
		ctx.JSON(http.StatusOK, gin.H{
			"error": "系统错误 claims 断言失败 ",
		})
		return
	}

	ctx.String(http.StatusOK, "test")
}

type UserClaims struct {
	jwt.RegisteredClaims

	// 声明自己的要放进去 token 里面的数据

	Uid int64

	UserAgent string
}
