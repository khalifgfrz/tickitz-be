package handlers

import (
	models "khalifgfrz/tickitz-be/internal/models/moviesAdd"
	"khalifgfrz/tickitz-be/internal/repository"
	"khalifgfrz/tickitz-be/pkg"

	"github.com/gin-gonic/gin"
)

type HandlerGenre struct {
	repository.GenreRepoInterface
}

func NewGenre(r repository.GenreRepoInterface) *HandlerGenre {
	return &HandlerGenre{r}
}

func (h *HandlerGenre) PostGenre(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	genre := models.Genres{}

	if err := ctx.ShouldBind(&genre); err != nil {
		response.BadRequest("Insert genre  failed", err.Error())
		return
	}

	results, err := h.CreateGenres(&genre)
	if err != nil {
		response.InternalServerError("Internar server error", err.Error())
		return
	}
	response.Created("Genre has been created", results)
}

func (h *HandlerGenre) GetGenres(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	data, err := h.GetAllGenres()
	if err != nil {
		response.InternalServerError("Internal Server Error", err.Error())
	}
	if len(data) == 0 {
		response.NotFound("Data Not Found", "No datas available for the given criteria")
		return
	}

	response.Success("Data retrieved successfully", data)
}
