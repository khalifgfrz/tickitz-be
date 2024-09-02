package routers

import (
	"khalifgfrz/tickitz-be/internal/handlers"
	middleware "khalifgfrz/tickitz-be/internal/middlewares"
	"khalifgfrz/tickitz-be/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func orderRouter(g *gin.Engine, d *sqlx.DB) {
	router := g.Group("/order")

	var orderRepo repository.OrderRepositoryInterface = repository.NewOrderRepository(d)
	var orderDetailsRepo repository.OrderDetailsRepositoryInterface = repository.NewOrderDetailsRepository(d)
	var paymentRepo repository.PaymentsRepoInterface = repository.NewPaymentsRepository(d)
	var sales repository.SalesRepoInterface = repository.NewSalesRepository(d)
	var movies repository.MovieRepoInterface = repository.NewMovieRepository(d)
	handler := handlers.NewOrderHandler(orderRepo, orderDetailsRepo, paymentRepo, sales, movies)

	router.POST("/", middleware.Auth("user"), handler.CreateOrder)
	router.GET("/history", middleware.Auth("user"), handler.FetchHistory)
	router.GET("/", middleware.Auth("admin"), handler.FetchAll)
	router.GET("/:id", handler.FetchDetail)

	//additional
	router.GET("/payments", handler.GetPayments)
	router.GET("/dashboards", handler.FetchSales)
}
