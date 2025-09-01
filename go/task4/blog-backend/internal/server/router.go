package server

import (
	"github.com/eistert/web3-task/go/task4/blog-backend/internal/config"
	"github.com/eistert/web3-task/go/task4/blog-backend/internal/handlers"
	"github.com/eistert/web3-task/go/task4/blog-backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 把路由分组、挂中间件，组装 *gin.Engine
func NewRouter(db *gorm.DB, cfg config.Config) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api/v1")

	// health
	api.GET("/health", func(c *gin.Context) { c.String(200, "OK") })

	// auth
	authH := handlers.NewAuthHandler(db, cfg)
	api.POST("/auth/register", authH.Register)
	api.POST("/auth/login", authH.Login)

	// posts (public read)
	postH := handlers.NewPostHandler(db)
	api.GET("/posts", postH.List)
	api.GET("/posts/:id", postH.Get)
	// comments (public read)
	cmtH := handlers.NewCommentHandler(db)
	api.GET("/posts/:id/comments", cmtH.ListByPost)

	// authorized
	auth := api.Group("/")
	auth.Use(middleware.Auth(cfg))
	{
		auth.POST("/posts", postH.Create)
		auth.PUT("/posts/:id", postH.Update)
		auth.DELETE("/posts/:id", postH.Delete)

		auth.POST("/posts/:id/comments", cmtH.Create)
	}

	return r
}
