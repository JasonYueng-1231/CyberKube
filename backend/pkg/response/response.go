package response

import "github.com/gin-gonic/gin"

type R struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(200, R{Code: 0, Message: "success", Data: data})
}

func Error(c *gin.Context, httpStatus int, code int, msg string) {
    c.JSON(httpStatus, R{Code: code, Message: msg})
}

