package users

import (
	"gotrack/database"
	"gotrack/helpers/common"

	"github.com/gin-gonic/gin"
)

func Initiator(router *gin.Engine) {
	api := router.Group("/api/users")
	{
		api.POST("/login", Login)
		api.POST("/signup", SignUp)
	}
}

// Login godoc
// @Tags Users
// @Summary User Login
// @Description This endpoint is used for user login
// @Accept json
// @Produce json
// @Param loginRequest body LoginRequest true "Login Request"
// @Success 200 {object} common.APIResponse{data=LoginResponse} "Success"
// @Failure 400 {object} common.APIResponse "Bad Request"
// @Failure 500 {object} common.APIResponse "Internal Server Error"
// @Router /api/users/login [post]
func Login(ctx *gin.Context) {
	var (
		userRepo = NewRepository(database.DBConnections)
		userSrv  = NewService(userRepo)
	)

	token, err := userSrv.LoginService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}

	common.GenerateSuccessResponseWithData(ctx, "successfully login", token)
}

// SignUp godoc
// @Tags Users
// @Summary User Signup
// @Description This endpoint is used for user signup
// @Accept json
// @Produce json
// @Param signUpRequest body SignUpRequest true "Sign Up Request"
// @Success 200 {object} common.APIResponse "Success"
// @Failure 400 {object} common.APIResponse "Bad Request"
// @Failure 500 {object} common.APIResponse "Internal Server Error"
// @Router /api/users/signup [post]
func SignUp(ctx *gin.Context) {
	var (
		userRepo = NewRepository(database.DBConnections)
		userSrv  = NewService(userRepo)
	)

	err := userSrv.SignUpService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}

	common.GenerateSuccessResponse(ctx, "awesome, successfully create user")
}
