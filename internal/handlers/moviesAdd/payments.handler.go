package handlers

import (
	"khalifgfrz/tickitz-be/internal/repository"
	"khalifgfrz/tickitz-be/pkg"

	"github.com/gin-gonic/gin"
)

type HandlerPayment struct {
	repository.PaymentsRepoInterface
}

func NewPayment(r repository.PaymentsRepoInterface) *HandlerPayment {
	return &HandlerPayment{r}
}

func (h *HandlerPayment) GetPayments(ctx *gin.Context) {
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
