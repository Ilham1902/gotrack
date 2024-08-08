package orders

import (
	"errors"
	"fmt"
	"gotrack/helpers/common"
	"gotrack/middlewares"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	Create(ctx *gin.Context) (err error)
	GetAll(ctx *gin.Context) (result []Order, err error)
	GetById(ctx *gin.Context) (result Order, err error)
	Update(ctx *gin.Context) (err error)
	Delete(ctx *gin.Context) (err error)

	Delivery(ctx *gin.Context) (err error)
	Success(ctx *gin.Context) (err error)
}

type orderServices struct {
	repository Repository
	validate   *validator.Validate
}

// Create implements Service.
func (o *orderServices) Create(ctx *gin.Context) (err error) {
	var request = OrderRequest{}

	if err = ctx.BindJSON(&request); err != nil {
		return errors.New("invalid request")
	}

	// Validate request data
	if err = o.validate.Struct(request); err != nil {
		return errors.New("validation failed: " + err.Error())
	}

	dataUser, err := o.repository.FindEmployee(request.EmployeeID)
	if err != nil {
		return errors.New("employee not found")
	}

	if dataUser.Role == "owner" {
		return errors.New("this user is owner")
	}

	order := Order{
		EmployeeID:  request.EmployeeID,
		Customer:    request.Customer,
		Location:    request.Location,
		Description: request.Description,
		Status:      "Pending",
	}

	if err = o.repository.Create(&order); err != nil {
		return err
	}

	for i := range request.OrderDetails {
		request.OrderDetails[i].OrderID = int(order.ID)
		if common.IsEmptyField(request.OrderDetails[i].Item) {
			return errors.New("item Required")
		}
		if common.IsEmptyField(request.OrderDetails[i].Qty) {
			return errors.New("qty Required")
		}
	}

	if err = o.repository.CreateOrderDetails(request.OrderDetails); err != nil {
		o.repository.Delete(int(order.ID))
		return err
	}

	return
}

// Delete implements Service.
func (o *orderServices) Delete(ctx *gin.Context) (err error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid ID format")
	}

	exists, err := o.repository.IsOrderExists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("orders with ID does not exist")
	}

	var orderReq OrderRequest

	err = ctx.ShouldBind(&orderReq)
	if err != nil {
		return
	}

	if err = o.repository.Delete(id); err != nil {
		return err
	}

	return nil
}

// GetAll implements Service.
func (o *orderServices) GetAll(ctx *gin.Context) (result []Order, err error) {
	search := ctx.Query("search")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))    // Default to page 1
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10")) // Default to limit 10

	user, exists := ctx.Get("auth")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		ctx.Abort()
		return
	}

	loginData, ok := user.(middlewares.UserLoginRedis)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		ctx.Abort()
		return
	}

	return o.repository.GetAll(loginData.Role, int(loginData.UserId), search, page, limit)
}

// GetById implements Service.
func (o *orderServices) GetById(ctx *gin.Context) (result Order, err error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return Order{}, fmt.Errorf("invalid ID format")
	}

	exists, err := o.repository.IsOrderExists(id)
	if err != nil {
		return Order{}, err
	}
	if !exists {
		return Order{}, errors.New("orders with ID does not exist")
	}

	data, err := o.repository.GetByID(id)
	if err != nil {
		return Order{}, err
	}

	return data, nil
}

// Update implements Service.
func (o *orderServices) Update(ctx *gin.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid ID format")
	}

	exists, err := o.repository.IsOrderExists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("orders with ID does not exist")
	}

	var request OrderRequest
	if err = ctx.BindJSON(&request); err != nil {
		return errors.New("invalid request")
	}

	// Validate request data
	if err = o.validate.Struct(request); err != nil {
		return errors.New("validation failed: " + err.Error())
	}

	_, err = o.repository.FindEmployee(request.EmployeeID)
	if err != nil {
		return errors.New("employee not found")
	}

	order := Order{
		EmployeeID:  request.EmployeeID,
		Customer:    request.Customer,
		Location:    request.Location,
		Status:      request.Status,
		Description: request.Description,
	}

	// Convert OrderRequest details to OrderDetail
	var details []OrderDetail
	for _, detail := range request.OrderDetails {
		if common.IsEmptyField(detail.Item) {
			return errors.New("item required")
		}
		if common.IsEmptyField(detail.Qty) {
			return errors.New("qty required")
		}
		detail.OrderID = id
		details = append(details, detail)
	}

	if err = o.repository.Update(order, id, details); err != nil {
		return err
	}

	return nil
}

func (o *orderServices) Delivery(ctx *gin.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid ID format")
	}

	exists, err := o.repository.IsOrderExists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("orders with ID does not exist")
	}

	if err := o.repository.Delivery(id); err != nil {
		return err
	}

	return nil
}

// Success implements Service.
func (o *orderServices) Success(ctx *gin.Context) (err error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid ID format")
	}

	exists, err := o.repository.IsOrderExists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("orders with ID does not exist")
	}

	// Parse the form data
	if err = ctx.Request.ParseMultipartForm(5 << 20); err != nil {
		return errors.New("unable to parse form")
	}

	// Get the file from the form input
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		return errors.New("unable to get file from form")
	}
	defer file.Close()

	// Hash the file name
	fileName := HashFilename(header.Filename)

	// Save the file to the public folder
	publicDir := "./public"
	if err = os.MkdirAll(publicDir, os.ModePerm); err != nil {
		return errors.New("unable to create public directory")
	}
	dst, err := os.Create(filepath.Join(publicDir, fileName))
	if err != nil {
		return errors.New("unable to create file on server")
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return errors.New("unable to save file")
	}

	// Get IP address of the requester
	ip := ctx.ClientIP()

	if err := o.repository.Success(id, ip, fileName); err != nil {
		return err
	}

	return
}

func NewService(repository Repository) Service {
	return &orderServices{
		repository,
		validator.New(),
	}
}
