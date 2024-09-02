package routers

import (
	"khalifgfrz/tickitz-be/internal/handlers"
	middleware "khalifgfrz/tickitz-be/internal/middlewares"
	"khalifgfrz/tickitz-be/internal/repository"
	"khalifgfrz/tickitz-be/pkg"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func movieRouter(g *gin.Engine, d *sqlx.DB) {
	router := g.Group("/movie")

	var movieRepo repository.MovieRepoInterface = repository.NewMovieRepository(d)
	var genreRepo repository.GenreMovieRepoInterface = repository.NewGenreMovieRepository(d)
	var airingDateRepo repository.AiringDateRepoInterface = repository.NewAiringDateRepository(d)
	var airingTimeDateRepo repository.AiringTimeDateRepoInterface = repository.NewAiringTimeDateRepository(d, nil)
	var movieTimeRepo repository.MovieTimeRepoInteface = repository.NewMovieTimeRepository(d)
	var locationMovieTimeRepo repository.LocationMovieTimeRepoInterface = repository.NewLocationMovieRepository(d)
	var timesRepo repository.AiringTimeRepoInterface = repository.NewAiringTimeRepository(d)
	var locations repository.LocationRepoInterface = repository.NewLocationRepository(d)
	var gRepo repository.GenreRepoInterface = repository.NewGenresRepository(d)
	var cld pkg.Cloudinary = *pkg.NewCloudinaryUtil()
	handler := handlers.NewMovieRepository(movieRepo, genreRepo, airingDateRepo, airingTimeDateRepo, movieTimeRepo, locationMovieTimeRepo, timesRepo, locations, gRepo, cld, d)

	router.POST("/insert", middleware.Auth("admin"), handler.InsertMovies)
	router.GET("/", handler.GetMovies)
	router.GET("/:id", handler.GetMovieDetails)
	router.PATCH("/:id", middleware.Auth("admin"), handler.UpdateMovies)
	router.PATCH("/banner/:id", middleware.Auth("admin"), handler.BannerUpdate)
	router.DELETE("/:id", middleware.Auth("admin"), handler.MoviesDelete)

	//additional

	router.GET("/times", handler.GetAllAiringTime)
	router.GET("/locations", handler.GetLocations)
	router.GET("/genres", handler.GetGenres)

}
