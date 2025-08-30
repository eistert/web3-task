package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

/*
题目1：使用SQL扩展库进行查询
假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
要求 ：
编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
*/

// Employee 与表字段一一对应（db tag 告诉 sqlx 如何映射）
type Employee struct {
	ID         uint64  `db:"id"          json:"id"`
	Name       string  `db:"name"        json:"name"`
	Department string  `db:"department"  json:"department"`
	Salary     float64 `db:"salary"      json:"salary"`
}

// 按部门查询（要求一）
func queryEmployeesByDept(ctx context.Context, db *sqlx.DB, dept string) ([]Employee, error) {
	var list []Employee
	const q = `
		SELECT id, name, department, salary
		FROM employees
		WHERE department = ?`
	if err := db.SelectContext(ctx, &list, q, dept); err != nil {
		return nil, err
	}
	return list, nil
}

// 查询工资最高的员工（要求二）
func queryTopSalary(ctx context.Context, db *sqlx.DB) (Employee, error) {
	var e Employee
	// 方式1：ORDER BY + LIMIT 1
	const q = `
		SELECT id, name, department, salary
		FROM employees
		ORDER BY salary DESC
		LIMIT 1`
	err := db.GetContext(ctx, &e, q)
	return e, err
	// 方式2：MAX 子查询（效果等价）
	// const q2 = `
	// 	SELECT id, name, department, salary
	// 	FROM employees
	// 	WHERE salary = (SELECT MAX(salary) FROM employees)
	// 	LIMIT 1`
}

// 连接 sqlx.DB
func mustOpenDB() *sqlx.DB {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		// 示例 DSN（根据你的环境修改）
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

func main() {
	db := mustOpenDB()
	r := gin.Default()

	// GET /employees?dept=技术部
	r.GET("/employees", func(c *gin.Context) {
		dept := c.Query("dept")
		if dept == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing dept"})
			return
		}
		list, err := queryEmployeesByDept(c.Request.Context(), db, dept)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, list)
	})

	// GET /employees/top-salary
	r.GET("/employees/top-salary", func(c *gin.Context) {
		e, err := queryTopSalary(c.Request.Context(), db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, e)
	})

	log.Println("listening :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
