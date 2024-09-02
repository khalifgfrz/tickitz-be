package handlers

import (
	models "khalifgfrz/tickitz-be/internal/models/moviesAdd"
	"khalifgfrz/tickitz-be/internal/repository"
	"khalifgfrz/tickitz-be/pkg"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HandlerLocation struct {
	repository.LocationRepoInterface
}

func NewLocation(r repository.LocationRepoInterface) *HandlerLocation {
	return &HandlerLocation{r}
}

func (h *HandlerLocation) PostLocation(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	Location := models.Locations{}

	if err := ctx.ShouldBind(&Location); err != nil {
		response.BadRequest("Insert Location  failed", err.Error())
		return
	}

	results, err := h.CreateLocation(&Location)
	if err != nil {
		response.InternalServerError("Internar server error", err.Error())
		return
	}
	response.Created("Location has been created", results)
}

func (h *HandlerLocation) GetLocations(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	data, err := h.GetAllLocations()
	if err != nil {
		response.InternalServerError("Internal Server Error", err.Error())
	}
	if len(data) == 0 {
		response.NotFound("Data Not Found", "No datas available for the given criteria")
		return
	}

	response.Success("Data retrieved successfully", data)
}

func (h *HandlerLocation) PatchLocations(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	var location models.Locations
	if err := ctx.ShouldBind(&location); err != nil {
		response.BadRequest("Update genre movie failed", err.Error())
		return
	}
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.BadRequest("Update genre movie failed", err.Error())
		return
	}
	results, err := h.UpdateLocation(id, &location)
	if err != nil {
		response.InternalServerError("Internal server error", err.Error())
		return
	}
	response.Success("Location movie has been updated", results)
}

func (h *HandlerLocation) DeleteLocations(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.BadRequest("Delete genre movie failed", err.Error())
		return
	}

	results, err := h.DeleteLocation(id)
	if err != nil {
		response.InternalServerError("Internar server error", err.Error())
		return
	}
	response.Success("Genre movie has been deleted", results)
}
