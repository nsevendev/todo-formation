package apiresponse

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
	Path      string      `json:"path,omitempty"`
	Method    string      `json:"method,omitempty"`
}

func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Status:    http.StatusOK,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
		Path:      c.Request.URL.Path,
		Method:    c.Request.Method,
	})
}

func Created(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, APIResponse{
		Status:    http.StatusCreated,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
		Path:      c.Request.URL.Path,
		Method:    c.Request.Method,
	})
}

func Error(c *gin.Context, status int, message string, err interface{}) {
	c.JSON(status, APIResponse{
		Status:    status,
		Message:   message,
		Error:     err,
		Timestamp: time.Now().Format(time.RFC3339),
		Path:      c.Request.URL.Path,
		Method:    c.Request.Method,
	})
}