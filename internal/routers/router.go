package routers

import (
	middleware "khalifgfrz/tickitz-be/internal/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *gin.Engine {
	router := gin.Default()

	// config := cors.DefaultConfig()


	// router.Use(cors.New(config))

	// router.Use(cors.New(cors.Config{
    //     AllowAllOrigins:  true,
    //     AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    //     AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "application/json"},
    //     ExposeHeaders:    []string{"Content-Length"},
    //     AllowCredentials: true,
    //     MaxAge:           12 * time.Hour,
    // }))

	router.Use(middleware.CORSMiddleware())

	authRouter(router, db)
	orderRouter(router, db)
	movieRouter(router, db)

	return router
}
