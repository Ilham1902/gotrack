package main

import (
	"gotrack/database"
	"gotrack/helpers/swagger"
	"gotrack/modules/orders"
	"gotrack/modules/users"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	DBConnections *gorm.DB
	err           error
)

// @title GoTrack Documentation
// @version 1.0.0
// @description This is documentation GoTrack.
// @host gotrack-production-2b8d.up.railway.app

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	err = godotenv.Load("config/.env")
	if err != nil {
		panic("Error loading .env file")
	}

	database.Conn()
	db := database.DBConnections

	router := gin.Default()

	db.AutoMigrate(&users.User{}, &orders.Order{}, &orders.OrderDetail{}, &users.IPInfo{}, &users.DetailLocation{})

	swagger.Initiator(router)
	users.Initiator(router)
	orders.Initiator(router)

	router.Run(":" + os.Getenv("PORT"))
}
