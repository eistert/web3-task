package main

import (
	"context"  // 标准库上下文，给数据库/超时/链路传递用
	"errors"   // 自定义错误
	"log"      // 打日志
	"net/http" // HTTP 状态码
	"os"       // 读环境变量
	"sort"     // 排序，保证加锁顺序一致

	// 你的数据模型：Account/Transaction/TransferReq
	model "github.com/eistert/web3-task/go/task3/sql-statement/txstruct"

	"github.com/gin-gonic/gin" // Web 框架
	"gorm.io/driver/mysql"     // MySQL 驱动
	"gorm.io/gorm"
	"gorm.io/gorm/clause" // GORM 的 SQL 子句，这里用来加行级锁
)

func main() {
	// 示例 DSN：root:pass@tcp(127.0.0.1:3306)/bank?parseTime=true&charset=utf8mb4&loc=Local
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:Root!123456@tcp(127.0.0.1:3306)/gorm_demo?parseTime=true&charset=utf8mb4&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("open db:", err)
	}

	// 自动建表（与上面的 DDL 对齐；生产环境请用迁移脚本）
	if err := db.AutoMigrate(&model.Account{}, &model.Transaction{}); err != nil {
		log.Fatal("migrate:", err)
	}

	// 自动建/更新表结构（开发demo可用，生产建议迁移脚本）。
	r := gin.Default()

	r.POST("/transfer", func(c *gin.Context) {
		var req model.TransferReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := transfer(c.Request.Context(), db, req.FromID, req.ToID, req.Amount); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Println("listening :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

/************ 事务核心 ************/
// transfer 在一个事务中完成：余额检查 -> 扣款 -> 加款 -> 写交易流水
/*
用 标准库 context（不是 *gin.Context）传递取消/超时等；
传入 GORM 的 *gorm.DB；
fromID/toID/amount 是业务参数。
*/
func transfer(ctx context.Context, db *gorm.DB, fromID, toID uint64, amount int64) error {
	// 前置校验 自转和非正金额直接拒绝。
	if fromID == toID {
		return errors.New("cannot transfer to self")
	}
	if amount <= 0 {
		return errors.New("invalid amount")
	}

	// 开启事务
	// Transaction 会在闭包返回 nil 时提交，返回非 nil 时回滚。
	// WithContext(ctx) 把请求的链路上下文传入数据库操作。
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 固定顺序加锁，避免并发死锁
		// 把两个账户ID排序，稍后加锁按固定顺序进行。这样 A⇄B 和 B⇄A 并发时不互相等待造成死锁。
		ids := []uint64{fromID, toID}
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

		var accs []model.Account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id IN ?", ids).
			Find(&accs).Error; err != nil {
			return err
		}

		if len(accs) != 2 {
			return gorm.ErrRecordNotFound
		}

		// 映射出 from/to（因为上面是 IN 查询，顺序需还原）
		var fromAcc, toAcc *model.Account
		for i := range accs {
			switch accs[i].ID {
			case fromID:
				fromAcc = &accs[i]
			case toID:
				toAcc = &accs[i]
			}
		}

		if fromAcc == nil || toAcc == nil {
			return gorm.ErrRecordNotFound
		}

		// 余额校验
		if fromAcc.Balance < amount {
			return gorm.ErrInvalidData // 可自定义错误：余额不足
		}

		// 原子更新余额（使用表达式避免读-改-写竞态）
		// UPDATE accounts SET balance = balance - ? WHERE id=?。
		if err := tx.Model(&model.Account{}).
			Where("id = ?", fromID).
			Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Account{}).
			Where("id = ?", toID).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		// 写流水
		txRec := model.Transaction{
			FromAccountID: fromID,
			ToAccountID:   toID,
			Amount:        amount,
		}
		if err := tx.Create(&txRec).Error; err != nil {
			return err
		}

		return nil // 提交事务
	})
}
