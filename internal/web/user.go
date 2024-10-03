package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUsersHandler() *UserHandler {
	const (
		EmailRegexPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		PasswordRegexPattern = `^[A-Za-z\d]{8,}$`
	)
	emailExp := regexp.MustCompile(EmailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(PasswordRegexPattern, regexp.None)

	return &UserHandler{
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) RegisterRoute(server *gin.Engine) {
	ug := server.Group("users")
	{
		ug.POST("/signup", u.SignUp)
		ug.POST("/login", u.Login)
		ug.POST("/edit", u.Edit)
		ug.POST("/profile", u.Profile)
	}
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ComfirmPassword string `json:"comfirmPassword"`
	}

	var req SignUpReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}

	const (
		EmailRegexPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		PasswordRegexPattern = `^[A-Za-z\d]{8,}$`
	)

	//验证邮箱
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "你的邮箱格式错误",
		})
		return
	}
	//验证密码
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "密码必须包含字母和数字，并且至少8位",
		})
		return
	}
	if req.Password != req.ComfirmPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "两次输入密码不一致",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": req.Email,
	})
}

func (u *UserHandler) Login(ctx *gin.Context) {

}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {

}
