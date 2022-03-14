package util

import "github.com/gin-gonic/gin"

func NewError(ctx *gin.Context, status int, err error) {
	var error = HttpError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, error)
}

type HttpError struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" example:"status bad request"`
}
