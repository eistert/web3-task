package handlers

import (
	"errors"
	"net/http"

	"github.com/eistert/web3-task/go/task4/blog-backend/internal/models"
	"github.com/eistert/web3-task/go/task4/blog-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 文章 CRUD（需要鉴权，校验作者）
type PostHandler struct{ DB *gorm.DB }

func NewPostHandler(db *gorm.DB) *PostHandler { return &PostHandler{DB: db} }

type postCreateReq struct {
	Title   string `json:"title"   binding:"required"`
	Content string `json:"content" binding:"required"`
}
type postUpdateReq struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

func (h *PostHandler) Create(c *gin.Context) {
	var req postCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	uid := c.GetUint("userID")
	p := models.Post{Title: req.Title, Content: req.Content, UserID: uid}
	if err := h.DB.Create(&p).Error; err != nil {
		response.JSON(c, http.StatusInternalServerError, nil, "create post failed")
		return
	}
	if err := h.DB.Preload("User").First(&p, p.ID).Error; err != nil {
		response.JSON(c, http.StatusInternalServerError, nil, "load author failed")
		return
	}
	response.JSON(c, http.StatusCreated, p, "")
}

func (h *PostHandler) List(c *gin.Context) {
	var posts []models.Post
	if err := h.DB.Preload("User").Order("created_at DESC").Find(&posts).Error; err != nil {
		response.JSON(c, http.StatusInternalServerError, nil, "list posts failed")
		return
	}
	response.JSON(c, http.StatusOK, posts, "")
}

func (h *PostHandler) Get(c *gin.Context) {
	var p models.Post
	if err := h.DB.Preload("User").First(&p, "id = ?", c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.JSON(c, http.StatusNotFound, nil, "post not found")
			return
		}
		response.JSON(c, http.StatusInternalServerError, nil, "get post failed")
		return
	}
	response.JSON(c, http.StatusOK, p, "")
}

func (h *PostHandler) Update(c *gin.Context) {
	var p models.Post
	if err := h.DB.First(&p, "id = ?", c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.JSON(c, http.StatusNotFound, nil, "post not found")
			return
		}
		response.JSON(c, http.StatusInternalServerError, nil, "get post failed")
		return
	}
	if p.UserID != c.GetUint("userID") {
		response.JSON(c, http.StatusForbidden, nil, "only author can update")
		return
	}
	var req postUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	updates := map[string]any{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if len(updates) == 0 {
		response.JSON(c, http.StatusBadRequest, nil, "nothing to update")
		return
	}
	if err := h.DB.Model(&p).Updates(updates).Error; err != nil {
		response.JSON(c, http.StatusInternalServerError, nil, "update post failed")
		return
	}
	if err := h.DB.Preload("User").First(&p, p.ID).Error; err != nil {
		response.JSON(c, http.StatusInternalServerError, nil, "load author failed")
		return
	}
	response.JSON(c, http.StatusOK, p, "")
}

func (h *PostHandler) Delete(c *gin.Context) {
	var p models.Post
	if err := h.DB.First(&p, "id = ?", c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.JSON(c, http.StatusNotFound, nil, "post not found")
			return
		}
		response.JSON(c, http.StatusInternalServerError, nil, "get post failed")
		return
	}
	if p.UserID != c.GetUint("userID") {
		response.JSON(c, http.StatusForbidden, nil, "only author can delete")
		return
	}
	if err := h.DB.Delete(&p).Error; err != nil {
		response.JSON(c, http.StatusInternalServerError, nil, "delete post failed")
		return
	}
	c.Status(http.StatusNoContent)
}
