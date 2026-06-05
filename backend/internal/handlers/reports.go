package handlers

import (
	"net/http"
	"training-scheduler/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReportsHandler struct {
	DB *gorm.DB
}

// GET /api/reports (ユーザーの全日報を取得)
func (h *ReportsHandler) GetReports(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	var reports []models.Report
	if err := h.DB.Where("user_id = ?", userID).Order("date desc").Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reports)
}

// POST /api/reports (日報の保存・更新)
func (h *ReportsHandler) SaveReport(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid, ok := parseUserID(userID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	var req models.Report
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req.UserID = uid

	// (user_id, date) の複合ユニーク制約が models.go に定義済みのため有効
	err := h.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"content"}),
	}).Create(&req).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
