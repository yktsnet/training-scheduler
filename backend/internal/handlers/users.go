package handlers

import (
	"net/http"
	"strings"
	"training-scheduler/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UsersHandler struct {
	DB *gorm.DB
}

// GET /api/users (アニマル一覧取得)
func (h *UsersHandler) GetUsers(c *gin.Context) {
	var users []models.User
	if err := h.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// POST /api/users (新しいアニマルでログイン/登録)
func (h *UsersHandler) CreateUser(c *gin.Context) {
	var req struct {
		Emoji   string `json:"emoji" binding:"required"`
		Initial string `json:"initial"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Emoji is required"})
		return
	}

	initial := strings.TrimSpace(req.Initial)
	initial = strings.ToUpper(initial)
	if len(initial) == 0 || len(initial) > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Initial must be between 1 and 3 characters"})
		return
	}

	// すでに使われているかチェック
	var existing models.User
	if err := h.DB.Where("emoji = ?", req.Emoji).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already taken"})
		return
	}

	newUser := models.User{
		Emoji:   req.Emoji,
		Initial: initial,
	}
	if err := h.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newUser)
}

// DELETE /api/users/:id (アニマルデータの削除)
func (h *UsersHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// 関連データの削除 (Flaskの delete().delete() 相当)
		if err := tx.Where("user_id = ?", id).Delete(&models.Plan{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&models.Report{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&models.Progress{}).Error; err != nil {
			return err
		}
		// 最後にユーザー自身を削除
		if err := tx.Delete(&models.User{}, id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
