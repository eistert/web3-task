package models

import "gorm.io/gorm"

// GORM 模型（User/Post/Comment，含关系）
type User struct {
	gorm.Model
	Username string `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:255;not null"            json:"-"`
	Email    string `gorm:"size:128;uniqueIndex;not null" json:"email"`
	Posts    []Post `json:"-"` // 一对多
	// 可选：PostCount 统计字段（用于钩子维护）
}

type Post struct {
	gorm.Model
	Title   string `gorm:"size:200;not null"   json:"title"`
	Content string `gorm:"type:longtext;not null" json:"content"`
	UserID  uint   `gorm:"index;not null"      json:"user_id"`
	User    User   `json:"author"`
}

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null" json:"content"`
	UserID  uint   `gorm:"index;not null"     json:"user_id"`
	User    User   `json:"user"`
	PostID  uint   `gorm:"index;not null"     json:"post_id"`
	Post    Post   `json:"post"`
}
