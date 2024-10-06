package main

import (
	"github.com/Salvatore-Giordano/gin-redis-ip-limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	ginRedis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"jike/config"
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
		ExposeHeaders:    []string{"x-jwt-token"},
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
	redisStore, err := ginRedis.NewStore(10, "tcp", config.Config.Redis.Addr, "", []byte("etn&/1dTiCN;Th(tH/@<Xi&7>exV?<[*"),
		[]byte("*t:{y{xYKb@nTX21eH*v{c.8D\"/;Lu(1"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("jikeSession", redisStore))
	server.Use(middleware.NewLoginJwtMiddlewareBuilder().
		IgnorePaths("/users/login").
		IgnorePaths("/users/signup").
		Build())
	server.Use(iplimiter.NewRateLimiterMiddleware(redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	}), "general", 200, 60*time.Second))

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
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
