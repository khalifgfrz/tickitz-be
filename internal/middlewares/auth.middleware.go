package middleware

import (
	"khalifgfrz/tickitz-be/pkg"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response := pkg.NewResponse(ctx)
		var header string

		if header = ctx.GetHeader("Authorization"); header == "" {
			response.Unauthorized("Unauthorized", nil)
			return
		}

		if !strings.Contains(header, "Bearer") {
			response.Unauthorized("Inavlid Bearer Token", nil)
			return
		}

		token := strings.Replace(header, "Bearer ", "", -1)

		check, err := pkg.VerifyToken(token)
		if err != nil {
			response.Unauthorized("Invalid Bearer Token", nil)
			return
		}

		roleAllowed := false
		for _, role := range roles {
			if check.Role == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			response.Unauthorized("Forbidden: Insufficient permissions", nil)
			ctx.Abort()
			return
		}

		ctx.Set("id", check.Id)
		ctx.Set("email", check.Email)
		ctx.Set("role", check.Role)
		ctx.Next()
	}
}
