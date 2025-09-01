package main

import (
	"log"
	"os"

	"github.com/eistert/web3-task/go/task4/blog-backend/internal/config"
	"github.com/eistert/web3-task/go/task4/blog-backend/internal/database"
	"github.com/eistert/web3-task/go/task4/blog-backend/internal/models"
	"github.com/eistert/web3-task/go/task4/blog-backend/internal/server"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// 程序入口：装配配置、数据库、路由，启动 HTTP
func main() {
	_ = godotenv.Load(".env") // 没有也不报错
	cfg := config.Load()

	db := database.MustOpen(cfg)
	autoMigrate(db)

	r := server.NewRouter(db, cfg)
	log.Printf("listening :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil { // 5. 启动 HTTP
		log.Fatal(err)
	}
}

func autoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}); err != nil {
		log.Fatal("migrate:", err)
	}
	_ = os.Setenv("GORM_DIALECT", "mysql/sqlite") // 仅示意，无实际作用
}
