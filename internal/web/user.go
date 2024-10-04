package web

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"jike/internal/domain"
	"jike/internal/service"
	"net/http"
	"time"
)

type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUsersHandler(svc *service.UserService) *UserHandler {
	const (
		EmailRegexPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		PasswordRegexPattern = `^[A-Za-z\d]{8,}$`
	)
	emailExp := regexp.MustCompile(EmailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(PasswordRegexPattern, regexp.None)

	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) RegisterRoute(server *gin.Engine) {
	ug := server.Group("users")
	{
		ug.POST("/signup", u.SignUp)
		ug.POST("/login", u.LoginJwt)
		ug.POST("/logout", u.LogoutJwt)
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

	//调用service层
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrDuplicateEmail) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "邮箱冲突",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": req.Email,
	})
	return
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId int64 `json:"user_id"`
}

func (u *UserHandler) LoginJwt(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.InvalidEmailOrPassword) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "邮箱或者密码不对",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "系统错误",
		})
		return
	}
	//设置JWT
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
		UserId: user.Id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedString, err := token.SignedString([]byte("etn&/1dTiCN;Th(tH/@<Xi&7>exV?<[*"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "系统错误",
		})
		return
	}

	ctx.Header("x-jwt-token", signedString)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "登录成功",
		"data": map[string]interface{}{
			"user_id": user.Id,
			"email":   user.Email,
		},
	})
	return
}

func (u *UserHandler) LogoutJwt(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1,
		//Secure:   true,//https
		//HttpOnly:true,
	})
	sess.Save()
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "退出登录成功",
	})
	return
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.InvalidEmailOrPassword) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "邮箱或者密码不对",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "系统错误",
		})
		return
	}
	//设置session
	sess := sessions.Default(ctx)
	sess.Set("user_id", user.Id)
	sess.Set("update_time", time.Now().UnixMilli())
	sess.Options(sessions.Options{
		MaxAge: 86400,
		//Secure:   true,//https
		//HttpOnly:true,
	})
	sess.Save()

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "登录成功",
	})
	return
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1,
		//Secure:   true,//https
		//HttpOnly:true,
	})
	sess.Save()
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "退出登录成功",
	})
	return
}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {
	c, _ := ctx.Get("claims")
	_, ok := c.(*UserClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "系统错误",
		})
	}

}
