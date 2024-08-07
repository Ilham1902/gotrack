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
		api.POST("/track", middlewares.AuthorizeRole("owner"), Track)
	}
}
