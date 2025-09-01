package handlers

import (
	"net/http"

	"github.com/eistert/web3-task/go/task4/blog-backend/internal/config"
	"github.com/eistert/web3-task/go/task4/blog-backend/internal/models"
	"github.com/eistert/web3-task/go/task4/blog-backend/pkg/jwtutil"
	"github.com/eistert/web3-task/go/task4/blog-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB  *gorm.DB
	Cfg config.Config
}

func NewAuthHandler(db *gorm.DB, cfg config.Config) *AuthHandler {
	return &AuthHandler{DB: db, Cfg: cfg}
}

type registerReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Email    string `json:"email"    binding:"required,email"`
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 绑定请求 JSON → bcrypt 加密密码 → db.Create(&user)；
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	u := models.User{Username: req.Username, Password: string(hash), Email: req.Email}
	if err := h.DB.Create(&u).Error; err != nil {
		response.JSON(c, http.StatusBadRequest, nil, "username or email exists")
		return
	}

	response.JSON(c, http.StatusCreated, gin.H{"id": u.ID, "username": u.Username, "email": u.Email}, "")
}

// 校验用户名/密码 → GenerateToken → 返回 token。
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	var u models.User
	if err := h.DB.Where("username = ?", req.Username).First(&u).Error; err != nil {
		response.JSON(c, http.StatusUnauthorized, nil, "invalid username or password")
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) != nil {
		response.JSON(c, http.StatusUnauthorized, nil, "invalid username or password")
		return
	}

	tok, err := jwtutil.Issue(u.ID, u.Username, []byte(h.Cfg.JWTSecret), h.Cfg.TokenTTL)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, nil, "issue token failed")
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"token": tok}, "")
}
