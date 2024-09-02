package handlers

import (
	"khalifgfrz/tickitz-be/internal/models"
	"khalifgfrz/tickitz-be/internal/models/moviesAdd"
	"khalifgfrz/tickitz-be/internal/repository"
	"khalifgfrz/tickitz-be/pkg"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	repository.OrderRepositoryInterface
	repository.OrderDetailsRepositoryInterface
	repository.PaymentsRepoInterface
	repository.SalesRepoInterface
	repository.MovieRepoInterface
}

func NewOrderHandler(
	orderRepo repository.OrderRepositoryInterface,
	orderDetailsRepo repository.OrderDetailsRepositoryInterface,
	payments repository.PaymentsRepoInterface,
	sales repository.SalesRepoInterface,
	movies repository.MovieRepoInterface,
) *OrderHandler {
	return &OrderHandler{orderRepo, orderDetailsRepo, payments, sales, movies}
}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	body := models.Order{}

	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest("Create order failed 1", "Error")
		return
	}

	orderID, err := h.CreateData(&body)
	if err != nil {
		response.BadRequest("Create order failed 2", "Error")
		return
	}

	result, err := h.CreateOrderDetails(orderID, body.Orders)
	if err != nil {
		response.BadRequest("Create order failed 3", "Failed to create order details")
		return
	}

	response.Created("Create order success", result)
}

func (h *OrderHandler) FetchAll(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	orders, err := h.GetAllData()
	if err != nil {
		response.InternalServerError("Get data failed", "Error")
		return
	}
	for i := range *orders {
		order := &(*orders)[i]

		orderDetails, err := h.GetDetailOrder(order.Id)
		if err != nil {
			response.InternalServerError("Get data failed", "Error")
			return
		}
		order.Orders = orderDetails
	}
	response.Success("Get data success", orders)
}

func (h *OrderHandler) FetchDetail(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	id := ctx.Param("id")

	order, err := h.GetDetailData(id)

	if err != nil {
		response.InternalServerError("Get data failed", "Error")
		return
	}

	orderID := order.Id

	orderDetails, err := h.GetDetailOrder(orderID)

	if err != nil {
		response.InternalServerError("Get data failed", "Error")
		return
	}

	order.Orders = orderDetails

	response.Success("Get data success", order)
}

func (h *OrderHandler) FetchHistory(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	userID, exists := ctx.Get("id")

	if !exists {
		response.InternalServerError("User id not found", nil)
		return
	}

	id := userID.(string)

	history, err := h.GetHistoryOrder(id)
	if err != nil {
		response.InternalServerError("Get data failed", "Error")
		return
	}

	// Iterate through each order and get its details
	for i, order := range history {
		orderID := order.Id

		orderDetails, err := h.GetDetailOrder(orderID)
		if err != nil {
			response.InternalServerError("Get data failed", "Error")
			return
		}

		// Populate the Orders field with the details
		history[i].Orders = orderDetails
	}

	response.Success("Get data success", history)
}

//additional

func (h *OrderHandler) GetPayments(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	data, err := h.GetAllPayments()
	if err != nil {
		response.InternalServerError("Internal Server Error", err.Error())
	}
	if len(data) == 0 {
		response.NotFound("Data Not Found", "No datas available for the given criteria")
		return
	}

	response.Success("Data retrieved successfully", data)
}

func (h *OrderHandler) FetchSales(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	// Parse page and limit query parameters
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		response.BadRequest("Invalid or missing 'page' parameter", err.Error())
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		response.BadRequest("Invalid or missing 'limit' parameter", err.Error())
		return
	}

	// Prepare the query for fetching movies
	query := models.MoviesQuery{
		Page:  page,
		Limit: limit,
	}

	// Fetch movies
	movies, _, err := h.GetAllMovies(&query)
	if err != nil {
		response.InternalServerError("Get movie data failed", err.Error())
		return
	}

	// Fetch sales data
	salesData, err := h.GetAllSales()
	if err != nil {
		response.InternalServerError("Get sales data failed", err.Error())
		return
	}

	// Combine movies with sales data
	type MovieSales struct {
		ID         string               `json:"id"`
		Title      string               `json:"title"`
		DailySales []moviesAdd.GetSales `json:"daily_sales"`
	}

	var dashboard struct {
		Movies []MovieSales `json:"movies"`
	}

	for _, movie := range *movies {
		movieSales := MovieSales{
			ID:         movie.Id,
			Title:      movie.Title,
			DailySales: salesData,
		}

		dashboard.Movies = append(dashboard.Movies, movieSales)
	}

	// Send the response
	response.GetAllSuccess("Sales data fetched successfully", dashboard, nil)
}
