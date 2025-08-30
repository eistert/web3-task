package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Book 结构体，db tag 告诉 sqlx 列到字段的映射（类型安全地扫描）
type Book struct {
	ID     uint64  `db:"id"     json:"id"`
	Title  string  `db:"title"  json:"title"`
	Author string  `db:"author" json:"author"`
	Price  float64 `db:"price"  json:"price"` // 简化用 float64；生产可考虑“分”为单位（int64）或 decimal 库
}

func mustOpenDB2() *sqlx.DB {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		// 示例 DSN：按你的环境修改数据库名/密码
		dsn = "root:Root!123456@tcp(127.0.0.1:3306)/gorm_demo?parseTime=true&charset=utf8mb4&loc=Local"
	}

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatal("open db:", err)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)
	if err := db.Ping(); err != nil {
		log.Fatal("ping db:", err)
	}

	return db
}

// 用 sqlx 查询价格 > minPrice 的书籍（核心：类型安全映射到 []Book）
func queryBooksByPrice(ctx context.Context, db *sqlx.DB, minPrice float64) ([]Book, error) {
	const q = `
		SELECT id, title, author, price
		FROM books
		WHERE price > ?
		ORDER BY price DESC, id ASC`
	var list []Book
	if err := db.SelectContext(ctx, &list, q, minPrice); err != nil {
		return nil, err
	}
	return list, nil
}

func main() {
	db := mustOpenDB()
	r := gin.Default()

	// GET /books?min_price=50
	r.GET("/books", func(c *gin.Context) {
		minStr := c.DefaultQuery("min_price", "50")
		min, err := strconv.ParseFloat(minStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid min_price"})
			return
		}

		list, err := queryBooksByPrice(c.Request.Context(), db, min)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, list)
	})

	log.Println("listening :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
