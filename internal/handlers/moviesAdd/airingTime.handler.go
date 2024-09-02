package handlers

import (
	"khalifgfrz/tickitz-be/internal/repository"
	"khalifgfrz/tickitz-be/pkg"

	"github.com/gin-gonic/gin"
)

type HandlerAiringTime struct {
	repository.AiringTimeRepoInterface
}

func NewAiringTime(r repository.AiringTimeRepoInterface) *HandlerAiringTime {
	return &HandlerAiringTime{r}
}

func (h *HandlerAiringTime) GetAllAiringTime(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	data, err := h.GetAiringTime()
	if err != nil {
		response.InternalServerError("Internal Server Error", err.Error())
	}
	if len(data) == 0 {
		response.NotFound("Data Not Found", "No datas available for the given criteria")
		return
	}

	response.Success("Data retrieved successfully", data)
}
