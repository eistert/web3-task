package main

import (
	"context"

	"gorm.io/gorm"
)

// PostWithCount 用于承接“文章 + 评论数”聚合结果
type PostWithCount struct {
	Post         `gorm:"embedded"` // 把 posts.* 扫到内嵌的 Post 字段
	CommentCount int64             `gorm:"column:comment_count" json:"comment_count"`
}

// QueryUserPostsWithComments 返回指定用户的所有文章，并预加载每篇文章的评论
func QueryUserPostsWithComments(ctx context.Context, db *gorm.DB, userID uint64) ([]Post, error) {
	var posts []Post
	err := db.WithContext(ctx).
		Preload("Comments", func(tx *gorm.DB) *gorm.DB {
			// 可选：评论按时间升序
			return tx.Order("comments.created_at ASC")
		}).
		// 如需把评论者也带出来：.Preload("Comments.User")
		Where("user_id = ?", userID).
		Order("posts.created_at DESC").
		Find(&posts).Error
	return posts, err
}

// QueryTopPostByCommentCount 返回评论数量最多的文章（含评论数），并可预加载作者/评论
func QueryTopPostByCommentCount(ctx context.Context, db *gorm.DB) (PostWithCount, error) {
	var out PostWithCount
	err := db.WithContext(ctx).
		Model(&Post{}). // 主表：posts
		Select("posts.*, COUNT(comments.id) AS comment_count").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Group("posts.id").
		Order("comment_count DESC, posts.id ASC").
		// 可选预加载：作者和（部分）评论
		Preload("User").
		// 如需把评论也带出，取消下一行注释；注意大量数据时可加 Limit
		// Preload("Comments", func(tx *gorm.DB) *gorm.DB { return tx.Order("comments.created_at DESC").Limit(20) }).
		Take(&out).Error
	return out, err
}
