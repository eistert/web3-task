package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ★ Hook 1：文章创建后，给作者的 post_count 加 1（原子更新，避免并发问题）
func (p *Post) AfterCreate(tx *gorm.DB) error {
	log.Printf("[AfterCreate] post=%d user=%v", p.ID, p.UserID)
	if p.UserID == nil {
		return nil // 无作者则不统计
	}

	return tx.Model(&User{}).
		Where("id = ?", *p.UserID).
		UpdateColumn("post_count", gorm.Expr("post_count + 1")).Error
}

// ★ Hook 2：评论删除后，如该文章无评论，则把评论状态置为“无评论”
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	log.Printf("[AfterDelete] comment=%d post=%d", c.ID, c.PostID)
	var cnt int64
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt == 0 {
		return tx.Model(&Post{}).
			Where("id = ?", c.PostID).
			UpdateColumn("comment_status", "无评论").Error
	}
	return nil
}

/************ 简单路由用于验证钩子 ************/
type createUserReq struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required,email"`
}
type createPostReq struct {
	UserID  uint64 `json:"user_id"  binding:"required,gt=0"`
	Title   string `json:"title"    binding:"required"`
	Content string `json:"content"  binding:"required"`
}
type createCommentReq struct {
	PostID  uint64  `json:"post_id"  binding:"required,gt=0"`
	UserID  *uint64 `json:"user_id"` // 可空
	Content string  `json:"content"  binding:"required"`
}

func main() {
	// DSN：建议用环境变量覆盖
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		// 确保数据库已创建，如 blog；或改为你的库名
		dsn = "root:Root!123456@tcp(127.0.0.1:3306)/blog?parseTime=true&charset=utf8mb4&loc=Local"
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("open db:", err)
	}

	if err := db.AutoMigrate(&User{}, &Post{}, &Comment{}); err != nil {
		log.Fatal("migrate:", err)
	}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.String(http.StatusOK, "OK") })

	// 新建用户
	r.POST("/users", func(c *gin.Context) {
		var req createUserReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u := User{Name: req.Name, Email: req.Email}
		if err := db.Create(&u).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, u)
	})

	// 新建文章（触发 AfterCreate：post_count+1）
	r.POST("/posts", func(c *gin.Context) {
		var req createPostReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		p := Post{UserID: &req.UserID, Title: req.Title, Content: req.Content}
		if err := db.Create(&p).Error; err != nil { // ← 插入成功后，GORM 自动调用 p.AfterCreate(tx)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, p)
	})

	// 新建评论（可选：你也可以加一个 Comment.AfterCreate，把状态改为“有评论”）
	r.POST("/comments", func(c *gin.Context) {
		var req createCommentReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cm := Comment{PostID: req.PostID, UserID: req.UserID, Content: req.Content}
		// 如果你想创建评论时把 Post.CommentStatus 设为“有评论”，可以在这里顺手更新：
		// _ = db.Model(&Post{}).Where("id = ?", req.PostID).UpdateColumn("comment_status", "有评论")
		if err := db.Create(&cm).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, cm)
	})

	// 删除评论（触发 AfterDelete：若剩余 0 条，置为“无评论”）
	r.DELETE("/comments/:id", func(c *gin.Context) {
		var cm Comment
		if err := db.First(&cm, "id = ?", c.Param("id")).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := db.Delete(&cm).Error; err != nil { // ← 删除成功后，GORM 自动调用 cm.AfterDelete(tx)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})

	// 查看用户（验证 post_count）
	r.GET("/users/:id", func(c *gin.Context) {
		var u User
		if err := db.Preload("Posts").First(&u, "id = ?", c.Param("id")).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, u)
	})

	// 查看文章（验证 comment_status）
	r.GET("/posts/:id", func(c *gin.Context) {
		var p Post
		if err := db.Preload("Comments").First(&p, "id = ?", c.Param("id")).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, p)
	})

	log.Println("listening :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
