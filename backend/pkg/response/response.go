package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess         = 0
	CodeBadRequest      = 40000
	CodeUnauthorized    = 40100
	CodeForbidden       = 40300
	CodeNotFound        = 40400
	CodeConflict        = 40900
	CodeTooManyRequests = 42900
	CodeInternal        = 50000
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: CodeSuccess, Message: "success", Data: data})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{Code: code, Message: message, Data: nil})
}

func BadRequest(c *gin.Context, message string)    { Error(c, CodeBadRequest, message) }
func Unauthorized(c *gin.Context, message string)  { Error(c, CodeUnauthorized, message) }
func Forbidden(c *gin.Context, message string)     { Error(c, CodeForbidden, message) }
func NotFound(c *gin.Context, message string)       { Error(c, CodeNotFound, message) }
func InternalError(c *gin.Context, message string)  { Error(c, CodeInternal, message) }
