package users

import (
	"errors"
	"gotrack/helpers/common"
	"gotrack/middlewares"
	"time"

	"github.com/gin-gonic/gin"
)

type Service interface {
	LoginService(ctx *gin.Context) (result LoginResponse, err error)
	SignUpService(ctx *gin.Context) (err error)
	Track(ctx *gin.Context) (*IPInfo, error)
}

type UserService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &UserService{
		repository,
	}
}

func (service *UserService) LoginService(ctx *gin.Context) (result LoginResponse, err error) {
	var userReq LoginRequest

	err = ctx.ShouldBind(&userReq)
	if err != nil {
		return
	}

	err = userReq.ValidateLogin()
	if err != nil {
		return
	}

	user, err := service.repository.Login(userReq)
	if err != nil {
		return
	}

	if common.IsEmptyField(user.ID) {
		err = errors.New("invalid account")
		return
	}

	matches := common.CheckPassword(user.Password, userReq.Password)
	if !matches {
		err = errors.New("wrong username or password")
		return
	}

	jwtToken, err := middlewares.GenerateJwtToken()
	if err != nil {
		return
	}

	middlewares.DummyRedis[jwtToken] = middlewares.UserLoginRedis{
		UserId:    int64(user.ID),
		Username:  user.Username,
		Role:      user.Role,
		LoginAt:   time.Now(),
		ExpiredAt: time.Now().Add(time.Hour * 2),
	}

	// Get the IP address from the request
	ipAddress := ctx.ClientIP()

	// ctx.JSON(http.StatusOK, gin.H{"ip": ipAddress})

	// Update or create IP info for the user
	if err = service.repository.UpdateIPEmployee(user.ID, ipAddress); err != nil {
		common.GenerateErrorResponse(ctx, "Failed to update IP information")
		return
	}

	result.Token = jwtToken

	return
}

func (service *UserService) SignUpService(ctx *gin.Context) (err error) {
	var userReq SignUpRequest

	err = ctx.ShouldBind(&userReq)
	if err != nil {
		return err
	}

	err = userReq.ValidateSignUp()
	if err != nil {
		return err
	}

	user, err := userReq.ConvertToModelForSignUp()
	if err != nil {
		return err
	}

	err = service.repository.SignUp(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) FindByID(id int) (User, error) {
	return s.repository.FindByID(uint(id))
}

// Track implements Service.
func (service *UserService) Track(ctx *gin.Context) (*IPInfo, error) {
	var request struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := ctx.BindJSON(&request); err != nil {
		return nil, errors.New("invalid request")
	}

	user, err := service.repository.FindByID(request.UserID)
	if err != nil {
		return nil, errors.New("id employee tidak ditemukan")
	}

	geoLocation, err := GetGeoLocation(user.IP)
	if err != nil {
		return nil, errors.New("failed to get geo location")
	}

	geoLocation.UserID = user.ID
	return geoLocation, nil
}
