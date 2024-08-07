package users

import (
	"gotrack/middlewares"

	"github.com/gin-gonic/gin"
)

func Initiator(router *gin.Engine) {
	api := router.Group("/api/users")
	{
		api.POST("/login", Login)
		api.POST("/signup", SignUp)
	}

	auth := router.Group("/api/users")
	auth.Use(middlewares.JwtMiddleware())
	auth.Use(middlewares.Logging())
	{
		auth.POST("/track", middlewares.AuthorizeRole("owner"), Track)
		auth.POST("", middlewares.AuthorizeRole("owner"), GetList)
	}
}
