package errcode

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// D attaches a detailed message and error to the context.
func D(c *gin.Context, status int, msg ...string) {
	c.JSON(getHttpStatusFromCode(status), gin.H{
		"msg":  strings.Join(msg, "\n"),
		"code": status,
	})
}

// S attaches only a simple message
func S(c *gin.Context, httpStatus int, msg ...string) {
	c.JSON(httpStatus, gin.H{
		"msg": strings.Join(msg, "\n"),
	})
}

func getHttpStatusFromCode(status int) (httpStatus int) {
	if status < 100 {
		return http.StatusNotFound
	} else {
		return http.StatusBadRequest
	}
}
