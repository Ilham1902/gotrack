package orders

import (
	"gotrack/database"
	"gotrack/helpers/common"

	"github.com/gin-gonic/gin"
)

// Create godoc
// @Summary Create a new order
// @Description Creates a new order with details provided in the request body.
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body OrderRequest true "Order data"
// @Success 200 {object} common.APIResponse{data=OrderRequest} "Successfully created the order"
// @Failure 400 {object} common.APIResponse "Bad request, invalid data provided"
// @Failure 500 {object} common.APIResponse "Internal server error"
// @Router /api/order [post]
func Create(ctx *gin.Context) {
	var (
		orderRepo = NewRepository(database.DBConnections)
		orderSrv  = NewService(orderRepo)
	)

	err := orderSrv.Create(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}

	common.GenerateSuccessResponse(ctx, "successfully added Order data")
}

// GetAllOrders godoc
// @Summary Get all orders
// @Description Get all orders with search and pagination
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param search query string false "Search term"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} common.APIResponse "Successfully retrieved orders"
// @Failure 400 {object} common.APIResponse "Bad request, invalid data provided"
// @Failure 500 {object} common.APIResponse "Internal server error"
// @Router /api/orders [get]
func GetAll(ctx *gin.Context) {
	var (
		orderRepo = NewRepository(database.DBConnections)
		orderSrv  = NewService(orderRepo)
	)

	data, err := orderSrv.GetAll(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}

	common.GenerateSuccessResponseWithListData(ctx, "successfully Get Order data", int64(len(data)), data)
}

// Update godoc
// @Summary Update a new order
// @Description Updates a new order with details provided in the request body.
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body OrderRequest true "Order data"
// @Success 200 {object} common.APIResponse{data=OrderRequest} "Successfully updated the order"
// @Failure 400 {object} common.APIResponse "Bad request, invalid data provided"
// @Failure 500 {object} common.APIResponse "Internal server error"
// @Router /api/order/{id} [put]
func Update(ctx *gin.Context) {
	var (
		orderRepo = NewRepository(database.DBConnections)
		orderSrv  = NewService(orderRepo)
	)

	err := orderSrv.Update(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}

	common.GenerateSuccessResponse(ctx, "successfully updated Order data")
}

// Delete godoc
// @Tags Order
// @Summary Delete a order by ID
// @Description Remove a order from the database by its ID
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} common.APIResponse "Success"
// @Failure 400 {object} common.APIResponse "Invalid ID format"
// @Failure 404 {object} common.APIResponse "Order not found"
// @Failure 500 {object} common.APIResponse "Internal server error"
// @Security Bearer
// @Router /api/order/{id} [delete]
func Delete(ctx *gin.Context) {
	var (
		ordersRepo = NewRepository(database.DBConnections)
		ordersSrv  = NewService(ordersRepo)
	)

	err := ordersSrv.Delete(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}

	common.GenerateSuccessResponse(ctx, "successfully delete order")
}
