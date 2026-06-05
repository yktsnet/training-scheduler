package handlers

import (
	"net/http"
	"strconv"
	"training-scheduler/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OverviewHandler struct {
	DB *gorm.DB
}

// GET /api/overview (選択されたメニューの進捗状況のみを取得)
func (h *OverviewHandler) GetOverview(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	uid, ok := parseUserID(userID)
	if !ok {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	var activeMenus []struct {
		models.Menu
	}
	err := h.DB.Table("menus").
		Select("menus.*").
		Joins("JOIN plans ON plans.menu_id = menus.id").
		Where("plans.user_id = ?", uid).
		Order("menus.id ASC").
		Scan(&activeMenus).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch active plans"})
		return
	}

	var progressList []models.Progress
	h.DB.Where("user_id = ?", uid).Find(&progressList)
	progMap := make(map[uint]models.Progress)
	for _, p := range progressList {
		progMap[p.MenuID] = p
	}

	type OverviewItem struct {
		MenuID      uint    `json:"menu_id"`
		MenuName    string  `json:"menu_name"`
		BaseDays    int     `json:"base_days"`
		StartDate   string  `json:"start_date"`
		TargetDays  int     `json:"target_days"`
		OffsetDays  float64 `json:"offset_days"`
		StatusMemo  string  `json:"status_memo"`
		IsCompleted bool    `json:"is_completed"`
	}

	var response []OverviewItem
	for _, m := range activeMenus {
		item := OverviewItem{
			MenuID:   m.ID,
			MenuName: m.Name,
			BaseDays: m.Days,
		}
		if p, ok := progMap[m.ID]; ok {
			item.StartDate = p.StartDate
			item.TargetDays = p.TargetDays
			item.OffsetDays = p.OffsetDays
			item.StatusMemo = p.StatusMemo
			item.IsCompleted = p.IsCompleted
		} else {
			item.TargetDays = m.Days
		}
		response = append(response, item)
	}
	c.JSON(http.StatusOK, response)
}

// POST /api/overview (進捗データの更新)
func (h *OverviewHandler) UpdateOverview(c *gin.Context) {
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

	var req models.Progress
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req.UserID = uid

	err := h.DB.Where("user_id = ? AND menu_id = ?", req.UserID, req.MenuID).
		Assign(req).
		FirstOrCreate(&req).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// parseUserID は X-User-Id ヘッダーの文字列を uint に変換する。
// 不正な値の場合は (0, false) を返す。
func parseUserID(idStr string) (uint, bool) {
	id, err := strconv.ParseUint(idStr, 10, 64)
	return uint(id), err == nil
}
