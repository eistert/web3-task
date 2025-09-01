package response

import "github.com/gin-gonic/gin"

// 统一返回格式：{ "data": ..., "error": "" }
func JSON(c *gin.Context, code int, data any, errMsg string) {
	c.JSON(code, gin.H{
		"data":  data,
		"error": errMsg,
	})
}
