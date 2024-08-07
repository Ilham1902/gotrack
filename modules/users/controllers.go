package users

import (
	"gotrack/database"
	"gotrack/helpers/common"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Tags Users
// @Summary User Login
// @Description This endpoint is used for user login
// @Accept json
// @Produce json
// @Param loginRequest body LoginRequest true "Login Request"
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

// Track Employee godoc
// @Tags Users
// @Summary Track Employee
// @Description This endpoint is used for Track Employee
// @Accept json
// @Produce json
// @Param TrackRequest body TrackRequest true "Sign Up Request"
// @Security Bearer
// @Router /api/users/track [post]
func Track(ctx *gin.Context) {
	var (
		userRepo = NewRepository(database.DBConnections)
		userSrv  = NewService(userRepo)
	)

	data, err := userSrv.Track(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}

	common.GenerateSuccessResponseWithData(ctx, "awesome, successfully create user", data)
}

// GetAllUser godoc
// @Summary Get all users
// @Description Get all users with search and pagination
// @Tags Users
// @Accept  json
// @Produce  json
// @Param search query string false "Search term"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Security Bearer
// @Router /api/users [get]
func GetList(ctx *gin.Context) {
	var (
		userrRepo = NewRepository(database.DBConnections)
		userrSrv  = NewService(userrRepo)
	)

	data, err := userrSrv.GetAll(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}

	common.GenerateSuccessResponseWithListData(ctx, "successfully Get User data", int64(len(data)), data)
}
