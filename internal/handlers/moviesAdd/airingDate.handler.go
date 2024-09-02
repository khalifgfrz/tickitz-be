package handlers

import (
	"fmt"
	models "khalifgfrz/tickitz-be/internal/models/moviesAdd"
	"khalifgfrz/tickitz-be/internal/repository"
	"khalifgfrz/tickitz-be/pkg"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type HandlerAiringDate struct {
	repository.AiringDateRepoInterface
	DB *sqlx.DB
}

func NewAiringDate(r repository.AiringDateRepoInterface, db *sqlx.DB) *HandlerAiringDate {
	return &HandlerAiringDate{r, db}
}

func (h *HandlerAiringDate) PostAiringDate(ctx *gin.Context) {
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

	airingDates := models.AiringDate{}

	if err := ctx.ShouldBind(&airingDates); err != nil {
		tx.Rollback()
		response.BadRequest("Insert airing date failed", err.Error())
		return
	}

	results, err := h.CreateAiringDate(tx, &airingDates)
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

	response.Created("Airing date inserted successfully", results)
}

func (h *HandlerAiringDate) GetAllAiringDate(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	data, err := h.GetAiringDate()
	if err != nil {
		response.InternalServerError("Internal Server Error", err.Error())
	}
	if data == nil {
		response.NotFound("Data Not Found", "No datas available for the given criteria")
		return
	}

	response.Success("Data retrieved successfully", data)
}

func (h *HandlerAiringDate) GetAiringDateInput(ctx *gin.Context) {
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

	input := models.AiringDate{}
	if err := ctx.ShouldBind(&input); err != nil {
		tx.Rollback()
		response.BadRequest("Failed to retrieve data, invalid input", err.Error())
		return
	}

	data, err := h.GetAiringDateByInput(tx, &input)
	if err != nil {
		tx.Rollback()
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}

	if data == nil {
		tx.Rollback()
		response.NotFound("Data Not Found", "No data available for the given criteria")
		return
	}

	// Commit transaksi jika semua operasi berhasil
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to commit transaction", err.Error())
		return
	}

	response.Success("Data retrieved successfully", data)
}
