package handlers

import (
	"fmt"
	"khalifgfrz/tickitz-be/internal/models"
	"khalifgfrz/tickitz-be/internal/repository"
	"khalifgfrz/tickitz-be/pkg"
	"math/rand"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	repository.UserRepositoryInterface
	repository.AuthRepositoryInterface
	pkg.Cloudinary
}

func NewAuthHandler(userRepo repository.UserRepositoryInterface, authRepo repository.AuthRepositoryInterface, cld pkg.Cloudinary) *AuthHandler {
	return &AuthHandler{userRepo, authRepo, cld}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	body := models.Auth{}

	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest("Register failed", "Error")
		return
	}

	_, err := govalidator.ValidateStruct(&body)
	if err != nil {
		response.BadRequest("Register failed", "Error")
		return
	}

	body.Password, err = pkg.HashPassword(body.Password)
	if err != nil {
		response.BadRequest("Register failed", "Error")
		return
	}

	result, err := h.CreateData(&body)
	if err != nil {
		response.BadRequest("Register failed", "Error")
		return
	}

	response.Created("Register success", result)
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	body := models.Auth{}

	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest("Login failed", "Error")
		return
	}

	_, err := govalidator.ValidateStruct(&body)
	if err != nil {
		response.BadRequest("Login failed", "Error")
		return
	}

	result, err := h.GetByEmail(body.Email)
	if err != nil {
		response.BadRequest("Login failed", "Error")
		return
	}

	err = pkg.VerifyPassword(result.Password, body.Password)
	if err != nil {
		response.Unauthorized("Wrong password", "Error")
		return
	}

	jwt := pkg.NewJWT(result.Id, result.Email, result.Role)
	token, err := jwt.GenerateToken()
	if err != nil {
		response.Unauthorized("Failed generate token", "Error")
		return
	}

	response.Success("Login success", gin.H{
        "token": token,
        "user": gin.H{
            "email":      result.Email,
            "role":       result.Role,
            "image":       result.Image,
        },
    })
}

func (h *AuthHandler) Update(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	userID, exists := ctx.Get("id")
	if !exists {
		response.NotFound("User doesn't exist", nil)
		return
	}
	id := userID.(string)
	body := models.User{}
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest("Update data failed", "Error")
		return
	}
	file, header, err := ctx.Request.FormFile("image")
	if err == nil {
		mimeType := header.Header.Get("Content-Type")
		if mimeType != "image/jpg" && mimeType != "image/jpeg" && mimeType != "image/png" {
			response.BadRequest("Create User failed, upload file failed, file is not supported", nil)
			return
		}

		if header.Size > 2*1024*1024 {
			response.BadRequest("Create User failed, upload file failed, file size exceeds 2 MB", nil)
			return
		}

		randomNumber := rand.Int()
		fileName := fmt.Sprintf("user-image-%d", randomNumber)
		uploadResult, err := h.UploadFile(ctx, file, fileName)
		if err != nil {
			response.BadRequest("Create User failed, upload file failed", "Error")
			return
		}
		imageURL := uploadResult.SecureURL
		body.Image = imageURL
	}

	if body.Password != "" {
		body.Password, err = pkg.HashPassword(body.Password)
		if err != nil {
			response.BadRequest("Update data failed", "Error")
			return
		}
	}

	result, err := h.UpdateData(&body, id)
	if err != nil {
		response.InternalServerError("Update data failed", "Error")
		return
	}

	response.Success("Update data success", result)
}

// func (h *AuthHandler) FetchAll(ctx *gin.Context) {
// 	response := pkg.NewResponse(ctx)

// 	result, err := h.GetAllData()
// 	if err != nil {
// 		response.InternalServerError("get data failed", "Error")
// 		return
// 	}

// 	response.Success("get data success", result)
// }

func (h *AuthHandler) FetchDetail(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	userID, exists := ctx.Get("id")
	if !exists {
		response.NotFound("User doesn't exist", nil)
		return
	}
	id := userID.(string)
	result, err := h.GetDetailData(id)
	if err != nil {
		response.InternalServerError("Get data failed", "Error")
		return
	}

	response.Success("Get data success", result)
}

// func (h *AuthHandler) Delete(ctx *gin.Context) {
// 	response := pkg.NewResponse(ctx)
// 	userID, exists := ctx.Get("id")
// 	if !exists {
// 		response.NotFound("User doesn't exist", nil)
// 		return
// 	}
// 	id := userID.(string)
// 	result, err := h.DeleteData(id)
// 	if err != nil {
// 		response.InternalServerError("Delete data failed", "Error")
// 		return
// 	}

// 	response.Success("Delete data success", result)
// }
