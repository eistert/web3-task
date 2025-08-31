package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/************ 程序入口 ************/
func main1() {
	dsn := getDSN()
	db := mustOpenDB(dsn)

	if err := autoMigrate(db); err != nil {
		log.Fatal("migrate:", err)
	}

	r := buildRouter(db)

	log.Println("listening :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

/************ 基础设施：DSN / DB / 迁移 / 路由 ************/

func getDSN() string {
	if v := os.Getenv("MYSQL_DSN"); v != "" {
		return v
	}
	// 默认 DSN（按需修改）
	return "root:Root!123456@tcp(127.0.0.1:3306)/gorm_demo?parseTime=true&charset=utf8mb4&loc=Local"
}

func mustOpenDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("open db:", err)
	}
	return db
}

func autoMigrate(db *gorm.DB) error {
	// 假设 User / Post / Comment 已定义
	return db.AutoMigrate(&User{}, &Post{}, &Comment{})
}

func buildRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 健康检查
	r.GET("/health", healthHandler)

	// 造测试数据（仅开发期使用）
	r.POST("/seed", seedHandler(db))

	// 查询：某用户的文章 + 评论
	r.GET("/users/:id/posts-with-comments", userPostsWithCommentsHandler(db))

	// 查询：评论最多的文章
	r.GET("/posts/top-by-comments", topPostByCommentsHandler(db))

	return r
}

/************ Handlers（只做解析入参 & 输出响应） ************/

func healthHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func seedHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := User{Name: "Alice", Email: "alice@example.com"}
		if err := db.Create(&u).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		p := Post{
			UserID:  &u.ID,
			Title:   "Hello GORM",
			Content: "This is my first post.",
		}
		if err := db.Create(&p).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		cm := Comment{
			PostID:  p.ID,
			UserID:  &u.ID,
			Content: "Nice post!",
		}
		if err := db.Create(&cm).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id":    u.ID,
			"post_id":    p.ID,
			"comment_id": cm.ID,
		})
	}
}

func userPostsWithCommentsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil || id64 == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		list, err := QueryUserPostsWithComments(c.Request.Context(), db, uint64(id64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

func topPostByCommentsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := QueryTopPostByCommentCount(c.Request.Context(), db)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "no posts"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}
