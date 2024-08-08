package users

import (
	"errors"
	"fmt"
	"gotrack/helpers/common"
	"gotrack/middlewares"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Service interface {
	LoginService(ctx *gin.Context) (result LoginResponse, err error)
	SignUpService(ctx *gin.Context) (err error)
	FindByID(ctx *gin.Context) (user User, err error)
	GetAll(ctx *gin.Context) (result []User, err error)
	Update(ctx *gin.Context) (err error)
	Delete(ctx *gin.Context) (err error)
	Track(ctx *gin.Context) (interface{}, error)
}

type UserService struct {
	repository Repository
	validate   *validator.Validate
}

func NewService(repository Repository) Service {
	return &UserService{
		repository,
		validator.New(),
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

// GetAll implements Service.
func (service *UserService) GetAll(ctx *gin.Context) (result []User, err error) {
	search := ctx.Query("search")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))    // Default to page 1
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10")) // Default to limit 10

	return service.repository.GetAll(search, page, limit)
}

// Update implements Service.
func (service *UserService) Update(ctx *gin.Context) (err error) {

	request := new(UpdatePayload)

	if err = ctx.BindJSON(&request); err != nil {
		return errors.New("invalid request")
	}

	// Validate request data
	if err = service.validate.Struct(request); err != nil {
		return errors.New("validation failed: " + err.Error())
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid ID format")
	}

	user := User{
		Username: request.Username,
		Password: request.Password,
		Role:     request.Role,
		IP:       "",
	}

	if err = service.repository.Update(user, id); err != nil {
		return err
	}

	return nil
}

// Delete implements Service.
func (service *UserService) Delete(ctx *gin.Context) (err error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid ID format")
	}

	exists, err := service.repository.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user with given ID does not exist")
		}
		return err
	}

	if (exists == User{}) {
		return errors.New("user with given ID does not exist")
	}

	if err = service.repository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (service *UserService) FindByID(ctx *gin.Context) (User, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return User{}, fmt.Errorf("invalid ID format")
	}

	data, err := service.repository.FindByID(uint(id))
	if err != nil {
		return User{}, err
	}

	return data, nil
}

// Track implements Service.
func (service *UserService) Track(ctx *gin.Context) (interface{}, error) {
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

	// geoLocation.UserID = user.ID
	return geoLocation, nil
}
