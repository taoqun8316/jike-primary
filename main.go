package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"jike/internal/repository"
	"jike/internal/repository/dao"
	"jike/internal/service"
	"jike/internal/web"
	"jike/internal/web/middleware"
	"strings"
	"time"
)

func main() {
	db := initDb()
	server := initWebServer()
	u := initUser(db)
	u.RegisterRoute(server)

	err := server.Run(":8081")
	if err != nil {
		panic(err)
	}
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:3000"},
		AllowMethods:     []string{"PUT", "GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			//开发环境
			if strings.Contains(origin, "http://127.0.0.1") {
				return true
			}
			//生产环境
			return strings.Contains(origin, "https://domain.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	server.Use(sessions.Sessions("jikeSession", cookie.NewStore([]byte("secret"))))
	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePaths("/users/login").
		IgnorePaths("/users/signup").
		Build())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	udao := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(udao)
	svc := service.NewUserService(repo)
	u := web.NewUsersHandler(svc)
	return u
}

func initDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/jike?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
