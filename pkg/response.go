package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Responder struct {
	C *gin.Context
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Meta    *Meta       `json:"meta,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type Meta struct {
	Total     int    `json:"totalData,omitempty"`
	TotalPage int    `json:"totalPage,omitempty"`
	Page      int    `json:"page,omitempty"`
	NextPage  string `json:"nextPage,omitempty"`
	PrevPage  string `json:"prevPage,omitempty"`
}

func NewResponse(ctx *gin.Context) *Responder {
	return &Responder{C: ctx}
}

func (r *Responder) Success(message string, data interface{}) {
	r.C.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	})
}

func (r *Responder) Created(message string, data interface{}) {
	r.C.JSON(http.StatusCreated, Response{
		Status:  http.StatusCreated,
		Message: message,
		Data:    data,
	})
}

func (r *Responder) GetAllSuccess(message string, data interface{}, meta *Meta) {
	r.C.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: message,
		Meta:    meta,
		Data:    data,
	})
}

func (r *Responder) BadRequest(message string, err interface{}) {
	r.C.JSON(http.StatusBadRequest, Response{
		Status:  http.StatusBadRequest,
		Message: message,
		Error:   err,
	})
	r.C.Abort()
}

func (r *Responder) Unauthorized(message string, err interface{}) {
	r.C.JSON(http.StatusUnauthorized, Response{
		Status:  http.StatusUnauthorized,
		Message: message,
		Error:   err,
	})
	r.C.Abort()
}

func (r *Responder) NotFound(message string, err interface{}) {
	r.C.JSON(http.StatusNotFound, Response{
		Status:  http.StatusNotFound,
		Message: message,
		Error:   err,
	})
	r.C.Abort()
}

func (r *Responder) InternalServerError(message string, err interface{}) {
	r.C.JSON(http.StatusInternalServerError, Response{
		Status:  http.StatusInternalServerError,
		Message: message,
		Error:   err,
	})
	r.C.Abort()
}
