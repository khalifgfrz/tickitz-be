package handlers

import (
	"fmt"
	"khalifgfrz/tickitz-be/internal/models"
	"khalifgfrz/tickitz-be/internal/models/moviesAdd"
	"khalifgfrz/tickitz-be/internal/repository"
	"khalifgfrz/tickitz-be/pkg"
	"math/rand"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type HandlerMovie struct {
	repository.MovieRepoInterface
	repository.GenreMovieRepoInterface
	repository.AiringDateRepoInterface
	repository.AiringTimeDateRepoInterface
	repository.MovieTimeRepoInteface
	repository.LocationMovieTimeRepoInterface
	repository.AiringTimeRepoInterface
	repository.LocationRepoInterface
	repository.GenreRepoInterface
	pkg.Cloudinary
	DB *sqlx.DB
}

func NewMovieRepository(
	mr repository.MovieRepoInterface,
	gmr repository.GenreMovieRepoInterface,
	ad repository.AiringDateRepoInterface,
	atd repository.AiringTimeDateRepoInterface,
	mt repository.MovieTimeRepoInteface,
	lmt repository.LocationMovieTimeRepoInterface,
	tm repository.AiringTimeRepoInterface,
	l repository.LocationRepoInterface,
	g repository.GenreRepoInterface,
	cld pkg.Cloudinary,
	db *sqlx.DB,
) *HandlerMovie {
	return &HandlerMovie{mr, gmr, ad, atd, mt, lmt, tm, l, g, cld, db}
}

func SplitCommaSeparatedInts(input string) ([]int, error) {
	var result []int
	if input == "" {
		return result, nil
	}

	parts := strings.Split(input, ",")
	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		value, err := strconv.Atoi(trimmedPart)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}

func (h *HandlerMovie) InsertMovies(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	var movies models.MoviesBody

	if h.DB == nil {
		response.InternalServerError("Database connection is not initialized", nil)
		return
	}

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
	if err := ctx.ShouldBind(&movies); err != nil {
		tx.Rollback()
		response.BadRequest("Create movie failed, invalid input", err.Error())
		return
	}
	genres, err := SplitCommaSeparatedInts(*movies.Genres)
	if err != nil {
		tx.Rollback()
		response.BadRequest("Invalid genre ID format", err.Error())
		return
	}

	airingTimes, err := SplitCommaSeparatedInts(*movies.AiringTime)
	if err != nil {
		tx.Rollback()
		response.BadRequest("Invalid airing time format", err.Error())
		return
	}

	locations, err := SplitCommaSeparatedInts(*movies.Locations)
	if err != nil {
		tx.Rollback()
		response.BadRequest("Invalid location format", err.Error())
		return
	}

	_, err = govalidator.ValidateStruct(&movies)
	if err != nil {
		tx.Rollback()
		response.BadRequest("Create movie failed", err.Error())
		return
	}

	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		tx.Rollback()
		response.BadRequest("Create movie failed, upload file is required", nil)
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType != "image/jpg" && mimeType != "image/jpeg" && mimeType != "image/png" {
		tx.Rollback()
		response.BadRequest("Create movie failed, upload file failed, file type is not supported", nil)
		return
	}

	if header.Size > 2*1024*1024 {
		tx.Rollback()
		response.BadRequest("Create movie failed, upload file failed, file size exceeds 2 MB", nil)
		return
	}

	randomNumber := rand.Int()
	fileName := fmt.Sprintf("movie-image-%d", randomNumber)
	uploadResult, err := h.UploadFile(ctx, file, fileName)
	if err != nil {
		tx.Rollback()
		response.BadRequest("Create movie failed, upload file failed", err.Error())
		return
	}

	imageURL := uploadResult.SecureURL
	movies.Image = &imageURL

	results, err := h.CreateMovie(tx, &movies)
	if err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if strings.Contains(err.Error(), "unique_name") {
				response.BadRequest("Create movie failed, movie name already exists", err.Error())
				return
			}
		}
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}

	for _, genreId := range genres {
		genreMovie := moviesAdd.GenreMovie{
			Movie_id: results.Id,
			Genre_id: genreId,
		}
		if _, err := h.InsertGenreMovie(tx, &genreMovie); err != nil {
			tx.Rollback()
			response.InternalServerError("Failed to insert genre data", err.Error())
			return
		}
	}

	for _, dateRange := range strings.Split(*movies.AiringDate, ",") {
		dates := strings.Split(dateRange, ",")
		var airingDate moviesAdd.AiringDate

		if len(dates) == 2 {
			airingDate = moviesAdd.AiringDate{
				Start_date: strings.TrimSpace(dates[0]),
				End_date:   strings.TrimSpace(dates[1]),
			}
		} else if len(dates) == 1 {
			airingDate = moviesAdd.AiringDate{
				Start_date: strings.TrimSpace(dates[0]),
				End_date:   strings.TrimSpace(dates[0]),
			}
		} else {
			tx.Rollback()
			response.BadRequest("Invalid date range provided", "Expected one or two dates in the format [yyyy-mm-dd] or [yyyy-mm-dd, yyyy-mm-dd]")
			return
		}

		existingAiringDate, err := h.GetAiringDateByInput(tx, &airingDate)
		if err != nil {
			tx.Rollback()
			response.InternalServerError("Failed to check airing date", err.Error())
			return
		}

		if existingAiringDate != nil {
			airingDate.Id = existingAiringDate.Id
		} else {
			insertedAiringDates, err := h.CreateAiringDate(tx, &airingDate)
			if err != nil {
				tx.Rollback()
				response.InternalServerError("Failed to insert airing date", err.Error())
				return
			}
			if len(insertedAiringDates) > 0 {
				airingDate.Id = insertedAiringDates[0].Id
			} else {
				tx.Rollback()
				response.InternalServerError("No airing date was inserted", "Expected a valid ID")
				return
			}
		}

		for _, airingTimeId := range airingTimes {
			newAiringTimeDate := moviesAdd.AiringTimeDate{
				Airing_time_id: airingTimeId,
				Date_id:        airingDate.Id,
			}
			insertedAiringTimeDateId, err := h.InsertAiringTimeDate(tx, &newAiringTimeDate)
			if err != nil {
				tx.Rollback()
				response.InternalServerError("Failed to insert airing time date", err.Error())
				return
			}

			movieTime := moviesAdd.MovieTime{
				Movie_id:            results.Id,
				Airing_time_date_id: insertedAiringTimeDateId.Id,
			}
			insertedMovieTime, err := h.CreateMovieTime(tx, &movieTime)
			if err != nil {
				tx.Rollback()
				response.InternalServerError("Failed to insert movie time", err.Error())
				return
			}

			for _, locationId := range locations {
				movieLocation := moviesAdd.LocationMovieTime{
					Movie_time_id: insertedMovieTime.ID,
					Location_id:   locationId,
				}
				if _, err := h.CreateLocationMovie(tx, &movieLocation); err != nil {
					tx.Rollback()
					response.InternalServerError("Failed to insert movie location", err.Error())
					return
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		response.InternalServerError("Internal Server Error", "Failed to commit transaction")
		return
	}

	response.Created("Data created", results)
}

func (h *HandlerMovie) GetMovies(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "16")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		response.BadRequest("Invalid or missing 'page' parameter", err.Error())
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		response.BadRequest("Invalid or missing 'limit' parameter", err.Error())
		return
	}

	query := models.MoviesQuery{
		Page:  page,
		Limit: limit,
	}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.BadRequest("Invalid query parameter", err.Error())
		return
	}

	movies, total, err := h.GetAllMovies(&query)
	if err != nil {
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}

	if len(*movies) == 0 {
		response.NotFound("Movie Not Found", "No movies available for the given criteria")
		return
	}

	totalPages := (total + query.Limit - 1) / query.Limit
	meta := &pkg.Meta{
		Total:     total,
		TotalPage: totalPages,
		Page:      query.Page,
		NextPage:  0,
		PrevPage:  0,
	}

	if query.Page+1 <= totalPages {
		meta.NextPage = query.Page + 1
	}

	if query.Page > 1 {
		meta.PrevPage = query.Page - 1
	}

	response.GetAllSuccess("Data fetched", movies, meta)
}

func (h *HandlerMovie) GetMovieDetails(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	id := ctx.Param("id")

	data, err := h.GetDetailMovie(id)
	if err != nil {
		response.NotFound("Movie Not Found", err.Error())
		return
	}

	if data == nil {
		response.NotFound("Movie Not Found", "No movie available for the given id")
		return
	}

	response.Success("Movie details retrieved successfully", data)
}

func (h *HandlerMovie) MoviesDelete(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	id := ctx.Param("id")

	data, err := h.DeleteMovie(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("Movie with id %s not found", id) {
			response.NotFound("Movie Not Found", err.Error())
			return
		}
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}

	response.Success("Movie deleted successfully", data)
}

func (h *HandlerMovie) BannerUpdate(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	id := ctx.Param("id")

	var bannerMovie models.MoviesBanner

	file, header, err := ctx.Request.FormFile("banner")
	if err != nil {
		response.BadRequest("Update banner movie failed, upload file is required", nil)
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType != "image/jpg" && mimeType != "image/jpeg" && mimeType != "image/png" {
		response.BadRequest("Update movie failed, upload file failed, file type is not supported", nil)
		return
	}

	if header.Size > 2*1024*1024 {
		response.BadRequest("Update movie failed, upload file failed, file size exceeds 2 MB", nil)
		return
	}

	randomNumber := rand.Int()
	fileName := fmt.Sprintf("movie-banner-%d", randomNumber)
	uploadResult, err := h.UploadFile(ctx, file, fileName)
	if err != nil {
		response.BadRequest("Update movie failed, upload file failed", err.Error())
		return
	}

	imageURL := uploadResult.SecureURL
	bannerMovie.Banner = imageURL

	results, err := h.UpdateBannerMovie(id, &bannerMovie)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if strings.Contains(err.Error(), "unique_name") {
				response.BadRequest("Update movie failed, movie name already exists", err.Error())
				return
			}
		}
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}
	response.Success("Data Updated", results)
}

func (h *HandlerMovie) UpdateMovies(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	id := ctx.Param("id")
	if id == "" {
		response.BadRequest("Invalid or missing 'id' parameter", nil)
		return
	}

	var movies models.MoviesBody

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

	if err := ctx.ShouldBind(&movies); err != nil {
		tx.Rollback()
		response.BadRequest("Invalid input", err.Error())
		return
	}

	_, err = govalidator.ValidateStruct(&movies)
	if err != nil {
		tx.Rollback()
		response.BadRequest("Validation failed", err.Error())
		return
	}

	if ctx.Request.MultipartForm != nil {
		file, header, err := ctx.Request.FormFile("image")

		if err != nil && err.Error() != "http: no such file" {
			tx.Rollback()
			response.BadRequest("Failed to retrieve file", nil)
			return
		}

		if file != nil {
			mimeType := header.Header.Get("Content-Type")
			if mimeType != "image/jpg" && mimeType != "image/jpeg" && mimeType != "image/png" {
				tx.Rollback()
				response.BadRequest("Unsupported file type", nil)
				return
			}

			if header.Size > 2*1024*1024 {
				tx.Rollback()
				response.BadRequest("File size exceeds 2 MB", nil)
				return
			}

			randomNumber := rand.Int()
			fileName := fmt.Sprintf("movie-image-%d", randomNumber)
			uploadResult, err := h.UploadFile(ctx, file, fileName)
			if err != nil {
				tx.Rollback()
				response.InternalServerError("Failed to upload file", err.Error())
				return
			}

			movies.Image = &uploadResult.SecureURL
		}
	}

	updatedMovie, err := h.UpdateMovieDetails(tx, id, &movies)
	if err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to update movie details", err.Error())
		return
	}

	if movies.Genres != nil {
		genres, err := SplitCommaSeparatedInts(*movies.Genres)
		if err != nil {
			tx.Rollback()
			response.BadRequest("Invalid genre ID format", err.Error())
			return
		}

		if err := h.UpdateGenreMovie(tx, id, genres); err != nil {
			tx.Rollback()
			response.InternalServerError("Failed to update movie genres", err.Error())
			return
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to commit transaction", err.Error())
		return
	}

	response.Success("Movie updated successfully", updatedMovie)
}

//additional

func (h *HandlerMovie) GetAllAiringTime(ctx *gin.Context) {
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

func (h *HandlerMovie) GetLocations(ctx *gin.Context) {
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

func (h *HandlerMovie) GetGenres(ctx *gin.Context) {
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
