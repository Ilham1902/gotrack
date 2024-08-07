package orders

import (
	"gotrack/middlewares"

	"github.com/gin-gonic/gin"
)

func Initiator(router *gin.Engine) {
	api := router.Group("/api/order")
	api.Use(middlewares.JwtMiddleware())
	api.Use(middlewares.Logging())
	{
		api.POST("", Create)
		api.GET("/list", GetAll)
		api.PUT(":id", middlewares.AuthorizeRole("owner"), Update)
		api.DELETE(":id", middlewares.AuthorizeRole("owner"), Delete)

		api.POST("/delivery/:id", middlewares.AuthorizeRole("employee"), Delivery)
		api.POST("/success/:id", middlewares.AuthorizeRole("employee"), Success)
	}
}
