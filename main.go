package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"jike/internal/web"
	"strings"
	"time"
)

func main() {
	server := gin.Default()

	u := web.NewUsersHandler()
	u.RegisterRoute(server)

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

	server.Run(":8081")
}
