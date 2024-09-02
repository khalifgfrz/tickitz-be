package handlers

import (
	"fmt"
	models "khalifgfrz/tickitz-be/internal/models/moviesAdd"
	"khalifgfrz/tickitz-be/internal/repository"
	"khalifgfrz/tickitz-be/pkg"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type HandlerAiringTimeDate struct {
	repository.AiringTimeDateRepoInterface
	DB *sqlx.DB
}

func NewAiringTimeDate(r repository.AiringTimeDateRepoInterface, db *sqlx.DB) *HandlerAiringTimeDate {
	return &HandlerAiringTimeDate{r, db}
}

func (h *HandlerAiringTimeDate) PostAiringTimeDate(ctx *gin.Context) {
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

	airingTimeDates := models.AiringTimeDate{}

	if err := ctx.ShouldBind(&airingTimeDates); err != nil {
		tx.Rollback()
		response.BadRequest("Insert airing date failed", err.Error())
		return
	}

	results, err := h.InsertAiringTimeDate(tx, &airingTimeDates)
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

	response.Created("Airing time date inserted successfully", results)
}

func (h *HandlerAiringTimeDate) GetAllAiringTimeDate(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	data, err := h.GetAiringTimeDate()
	if err != nil {
		response.InternalServerError("Internal Server Error", err.Error())
	}
	if len(data) == 0 {
		response.NotFound("Data Not Found", "No datas available for the given criteria")
		return
	}

	response.Success("Data retrieved successfully", data)
}
