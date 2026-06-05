package handlers

import (
	"net/http"
	"training-scheduler/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlansHandler struct {
	DB *gorm.DB
}

// GET /api/plans (ユーザーに紐づく全計画をメニュー名付きで取得)
func (h *PlansHandler) GetPlans(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	// Python版の辞書形式に合わせるための構造体
	type PlanResponse struct {
		ID       uint   `json:"id"`
		MenuName string `json:"menu_name"`
		Content  string `json:"content"`
	}

	var results []PlanResponse

	// PlanとMenuを結合して取得 (SELECT plans.id, menus.name as menu_name, ...)
	err := h.DB.Table("plans").
		Select("plans.id, menus.name as menu_name, plans.content").
		Joins("join menus on plans.menu_id = menus.id").
		Where("plans.user_id = ?", userID).
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// POST /api/plans/:id (特定の計画の内容を更新 - オートセーブ用)
func (h *PlansHandler) UpdatePlan(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")
	planID := c.Param("id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// IDとUserIDの両方が一致するレコードを更新 (他人のデータを守る)
	result := h.DB.Model(&models.Plan{}).
		Where("id = ? AND user_id = ?", planID, userID).
		Update("content", req.Content)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
