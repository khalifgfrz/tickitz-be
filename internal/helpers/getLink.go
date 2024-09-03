package helpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func BuildPaginationLink(ctx *gin.Context, baseURL string, page, limit int, total int) (string, string) {
	// fmt.Printf("Building pagination links with baseURL: %s, page: %d, limit: %d, total: %d\n", baseURL, page, limit, total)

	nextPage := ""
	prevPage := ""
	totalPage := (total + limit - 1) / limit

	if page < totalPage {
		nextPage = fmt.Sprintf("%s?page=%d&limit=%d", baseURL, page+1, limit)
	}
	if page > 1 {
		prevPage = fmt.Sprintf("%s?page=%d&limit=%d", baseURL, page-1, limit)
	}

	// fmt.Printf("Next page URL: %s\n", nextPage)
	// fmt.Printf("Previous page URL: %s\n", prevPage)

	return nextPage, prevPage
}
