package middleware

import (
	"net/http"
	"strings"

	"github.com/eistert/web3-task/go/task4/blog-backend/internal/config"
	"github.com/eistert/web3-task/go/task4/blog-backend/pkg/jwtutil"
	"github.com/eistert/web3-task/go/task4/blog-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// Gin 中间件（JWT 鉴权等）
func Auth(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头 Authorization: Bearer <token> 解析 JWT；
		authz := c.GetHeader("Authorization")
		parts := strings.SplitN(authz, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.JSON(c, http.StatusUnauthorized, nil, "missing or invalid Authorization header")
			c.Abort()
			return
		}
		claims, err := jwtutil.Parse(parts[1], []byte(cfg.JWTSecret))
		if err != nil {
			response.JSON(c, http.StatusUnauthorized, nil, "invalid token")
			c.Abort()
			return
		}

		// 把用户 ID 写入 c.Set("user_id", xxx)，供后续 handler 使用；
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
