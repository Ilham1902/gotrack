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

// GetByID godoc
// @Summary Get By ID users
// @Description Get By ID users with details provided in the request body.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security Bearer
// @Router /api/users/{id} [get]
func GetByID(ctx *gin.Context) {
	var (
		userRepo = NewRepository(database.DBConnections)
		userSrv  = NewService(userRepo)
	)

	data, err := userSrv.FindByID(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}

	common.GenerateSuccessResponseWithData(ctx, "successfully Get User data", data)
}

// Update godoc
// @Summary Update a new user
// @Description Updates a new user with details provided in the request body.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param UpdatePayload body UpdatePayload true "Update Request"
// @Security Bearer
// @Router /api/users/{id} [put]
func Update(ctx *gin.Context) {
	var (
		userRepo = NewRepository(database.DBConnections)
		userSrv  = NewService(userRepo)
	)

	err := userSrv.Update(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}

	common.GenerateSuccessResponse(ctx, "successfully updated User data")
}

// Delete godoc
// @Summary Delete a users by ID
// @Description Remove a users from the database by its ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security Bearer
// @Router /api/users/{id} [delete]
func Delete(ctx *gin.Context) {
	var (
		usersRepo = NewRepository(database.DBConnections)
		usersSrv  = NewService(usersRepo)
	)

	err := usersSrv.Delete(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}

	common.GenerateSuccessResponse(ctx, "successfully delete user")
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
