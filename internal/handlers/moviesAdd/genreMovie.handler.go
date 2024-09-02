package handlers

import (
	models "khalifgfrz/tickitz-be/internal/models/moviesAdd"
	"khalifgfrz/tickitz-be/internal/repository"
	"khalifgfrz/tickitz-be/pkg"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type HandlerGenreMovie struct {
	repository.GenreMovieRepoInterface
	DB *sqlx.DB
}

func NewGenreMovie(r repository.GenreMovieRepoInterface, db *sqlx.DB) *HandlerGenreMovie {
	return &HandlerGenreMovie{r, db}
}

func (h *HandlerGenreMovie) PostGenreMovie(ctx *gin.Context) {
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
		}
	}()

	genreMovies := models.GenreMovie{}

	if err := ctx.ShouldBind(&genreMovies); err != nil {
		tx.Rollback() // Rollback transaksi jika binding data gagal
		response.BadRequest("Insert genre movie failed", err.Error())
		return
	}

	// Gunakan transaksi untuk menyisipkan genre movie
	results, err := h.InsertGenreMovie(tx, &genreMovies)
	if err != nil {
		tx.Rollback() // Rollback transaksi jika insert gagal
		response.InternalServerError("Internal server error", err.Error())
		return
	}

	// Commit transaksi jika semua operasi berhasil
	if err := tx.Commit(); err != nil {
		tx.Rollback() // Rollback jika commit gagal
		response.InternalServerError("Failed to commit transaction", err.Error())
		return
	}

	response.Created("Genre movie has been created", results)
}

func (h *HandlerGenreMovie) DeleteGenreMovies(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	id := ctx.Param("id")

	results, err := h.DeleteGenreMovie(id)
	if err != nil {
		response.InternalServerError("Internar server error", err.Error())
		return
	}
	response.Success("Genre movie has been deleted", results)
}
