package handlers

import (
	"fmt"
	models "khalifgfrz/tickitz-be/internal/models/moviesAdd"
	"khalifgfrz/tickitz-be/internal/repository"
	"khalifgfrz/tickitz-be/pkg"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type HandlerMovieTime struct {
	repository.MovieTimeRepoInteface
	DB *sqlx.DB
}

func NewMovieTime(r repository.MovieTimeRepoInteface, db *sqlx.DB) *HandlerMovieTime {
	return &HandlerMovieTime{r, db}
}

func (h *HandlerMovieTime) PostMovieTimes(ctx *gin.Context) {
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

	movieTimes := models.MovieTime{}

	if err := ctx.ShouldBind(&movieTimes); err != nil {
		tx.Rollback()
		response.BadRequest("Create movie time failed, invalid input", err.Error())
		return
	}

	results, err := h.CreateMovieTime(tx, &movieTimes)
	if err != nil {
		tx.Rollback()
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}

	// Commit transaksi jika semua operasi berhasil
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to commit transaction", err.Error())
		return
	}

	response.Created("Movie time has been created", results)
}

func (h *HandlerMovieTime) GetMovieTimeById(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	id := ctx.Param("id")

	data, err := h.GetTimeByMovieId(id)
	if err != nil {
		response.InternalServerError("Internal Server Error", err.Error())
	}
	if data == nil {
		response.NotFound("Data Not Found", "No datas available for the given criteria")
		return
	}

	response.Success("Data retrieved successfully", data)
}

func (h *HandlerMovieTime) PatchMovieTime(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	id := ctx.Param("id")

	movieTime := models.MovieTime{}

	if err := ctx.ShouldBindJSON(&movieTime); err != nil {
		response.BadRequest("Update movie time failed, invalid input", err.Error())
		return
	}

	result, err := h.UpdateMovieTime(id, &movieTime)
	if err != nil {
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}

	response.Success("Data updated successfully", result)
}

func (h *HandlerMovieTime) DeleteMovieTimes(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.BadRequest("Delete genre movie failed", err.Error())
		return
	}

	results, err := h.DeleteMovieTime(id)
	if err != nil {
		response.InternalServerError("Internar server error", err.Error())
		return
	}
	response.Success("Genre movie has been deleted", results)
}
