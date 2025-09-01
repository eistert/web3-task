package handlers

import (
	"errors"
	"net/http"

	"github.com/eistert/web3-task/go/task4/blog-backend/internal/models"
	"github.com/eistert/web3-task/go/task4/blog-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 评论创建/查询

type CommentHandler struct{ DB *gorm.DB }

func NewCommentHandler(db *gorm.DB) *CommentHandler { return &CommentHandler{DB: db} }

type commentCreateReq struct {
	Content string `json:"content" binding:"required"`
}

func (h *CommentHandler) Create(c *gin.Context) {
	var req commentCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	var p models.Post
	if err := h.DB.First(&p, "id = ?", c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.JSON(c, http.StatusNotFound, nil, "post not found")
			return
		}
		response.JSON(c, http.StatusInternalServerError, nil, "get post failed")
		return
	}
	cm := models.Comment{Content: req.Content, UserID: c.GetUint("userID"), PostID: p.ID}
	if err := h.DB.Create(&cm).Error; err != nil {
		response.JSON(c, http.StatusInternalServerError, nil, "create comment failed")
		return
	}
	if err := h.DB.Preload("User").First(&cm, cm.ID).Error; err != nil {
		response.JSON(c, http.StatusInternalServerError, nil, "load user failed")
		return
	}
	response.JSON(c, http.StatusCreated, cm, "")
}

func (h *CommentHandler) ListByPost(c *gin.Context) {
	var list []models.Comment
	if err := h.DB.Preload("User").
		Where("post_id = ?", c.Param("id")).
		Order("created_at ASC").
		Find(&list).Error; err != nil {
		response.JSON(c, http.StatusInternalServerError, nil, "list comments failed")
		return
	}
	response.JSON(c, http.StatusOK, list, "")
}
