package utils

import (
	"math"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
}

func OK(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Fail(c *gin.Context, code int, err string) {
	c.JSON(code, Response{
		Success: false,
		Error:   err,
	})
}

func Paginated(c *gin.Context, code int, data interface{}, total int64, page, limit int) {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	c.JSON(code, PaginatedResponse{
		Success: true,
		Data:    data,
		Pagination: Pagination{
			CurrentPage: page,
			PerPage:     limit,
			Total:       total,
			TotalPages:  totalPages,
		},
	})
}

func Offset(page, limit int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * limit
}
