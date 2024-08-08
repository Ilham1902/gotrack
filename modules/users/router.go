package users

import (
	"gotrack/middlewares"

	"github.com/gin-gonic/gin"
)

func Initiator(router *gin.Engine) {
	api := router.Group("/api/users")
	{
		api.POST("/login", Login)
	}

	auth := router.Group("/api/users")
	auth.Use(middlewares.JwtMiddleware())
	auth.Use(middlewares.Logging())
	{
		auth.PUT(":id", Update)
		api.POST("/signup", middlewares.AuthorizeRole("owner"), SignUp)
		auth.GET("", middlewares.AuthorizeRole("owner"), GetList)
		auth.GET(":id", middlewares.AuthorizeRole("owner"), GetByID)
		auth.DELETE(":id", middlewares.AuthorizeRole("owner"), Delete)
		auth.POST("/track", middlewares.AuthorizeRole("owner"), Track)
	}
}
