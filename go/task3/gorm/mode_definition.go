package main

import (
	"time"
)

/************ GORM 模型 ************/

type User struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:64;not null"         json:"name"`
	Email     string    `gorm:"size:128;uniqueIndex;not null" json:"email"`
	PostCount int64     `gorm:"not null;default:0"       json:"post_count"` // ← 统计字段
	Posts     []Post    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"posts"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Post struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        *uint64   `gorm:"index"                    json:"user_id"`
	User          *User     `json:"user"`
	Title         string    `gorm:"size:200;not null"        json:"title"`
	Content       string    `gorm:"type:longtext;not null"   json:"content"`
	CommentStatus string    `gorm:"size:16;not null;default:'无评论'" json:"comment_status"` // ← 评论状态
	Comments      []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"comments"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Comment struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    uint64    `gorm:"index;not null"           json:"post_id"`
	Post      Post      `json:"post"`
	UserID    *uint64   `gorm:"index"                    json:"user_id"`
	User      *User     `json:"user"`
	Content   string    `gorm:"type:text;not null"       json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
