package handlers

import (
	"fmt"
	models "khalifgfrz/tickitz-be/internal/models/moviesAdd"
	"khalifgfrz/tickitz-be/internal/repository"
	"khalifgfrz/tickitz-be/pkg"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type HandlerLocationMovieTime struct {
	repository.LocationMovieTimeRepoInterface
	DB *sqlx.DB
}

func NewLocationMovieTime(r repository.LocationMovieTimeRepoInterface, db *sqlx.DB) *HandlerLocationMovieTime {
	return &HandlerLocationMovieTime{r, db}
}

func (h *HandlerLocationMovieTime) PostLocationMovieTime(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	// Mulai transaksi
	tx, err := h.DB.Beginx()
	if err != nil {
		response.InternalServerError("Failed to start transaction", err.Error())
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			response.InternalServerError("Transaction rolled back due to panic", fmt.Sprintf("%v", r))
		}
	}()

	data := models.LocationMovieTime{}

	if err := ctx.ShouldBind(&data); err != nil {
		tx.Rollback()
		response.BadRequest("Insert Location failed", err.Error())
		return
	}

	results, err := h.CreateLocationMovie(tx, &data)
	if err != nil {
		tx.Rollback()
		response.InternalServerError("Internal server error", err.Error())
		return
	}

	// Commit transaksi jika semua operasi berhasil
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to commit transaction", err.Error())
		return
	}

	response.Created("Location has been created", results)
}

func (h *HandlerLocationMovieTime) GetLocationMovieById(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	id := ctx.Param("id")

	data, err := h.GetMovieLocTimeById(id)
	if err != nil {
		response.InternalServerError("Internal Server Error", err.Error())
	}
	if data == nil {
		response.NotFound("Data Not Found", "No datas available for the given criteria")
		return
	}

	response.Success("Data retrieved successfully", data)
}

func (h *HandlerLocationMovieTime) DeleteLocationMovies(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	id := ctx.Param("id")

	results, err := h.DeleteLocationMovie(id)
	if err != nil {
		response.InternalServerError("Internar server error", err.Error())
		return
	}
	response.Success("Genre movie has been deleted", results)
}
