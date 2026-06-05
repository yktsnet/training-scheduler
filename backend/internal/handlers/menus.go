package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"training-scheduler/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MenusHandler struct {
	DB *gorm.DB
}

func (h *MenusHandler) GetMenus(c *gin.Context) {
	var menus []models.Menu
	if err := h.DB.Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, menus)
}

func (h *MenusHandler) SaveSelection(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-Id")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	var req struct {
		MenuIDs []uint `json:"menu_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		// チェックが外れたメニューの計画を削除
		if len(req.MenuIDs) > 0 {
			if err := tx.Where("user_id = ? AND menu_id NOT IN ?", userID, req.MenuIDs).Delete(&models.Plan{}).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Where("user_id = ?", userID).Delete(&models.Plan{}).Error; err != nil {
				return err
			}
		}

		// 既存のPlanを取得
		var existingPlans []models.Plan
		tx.Where("user_id = ?", userID).Find(&existingPlans)

		existingMap := make(map[uint]*models.Plan)
		for i := range existingPlans {
			existingMap[existingPlans[i].MenuID] = &existingPlans[i]
		}

		for _, mID := range req.MenuIDs {
			var menu models.Menu
			if err := tx.First(&menu, mID).Error; err != nil {
				continue
			}

			newHeader := fmt.Sprintf("【%s 研修計画（計%d日間）】", menu.Name, menu.Days)

			if plan, exists := existingMap[mID]; exists {
				// --- 修正箇所：既存プランのヘッダー（1行目）のみを更新 ---
				// 正規表現で最初の【...】部分を置換
				re := regexp.MustCompile(`^【.*?】`)
				updatedContent := re.ReplaceAllString(plan.Content, newHeader)
				
				if updatedContent != plan.Content {
					tx.Model(plan).Update("content", updatedContent)
				}
			} else {
				// 新規作成
				var sb strings.Builder
				for i := 0; i < menu.Days; i++ {
					sb.WriteString(fmt.Sprintf("%d日目：", i+1))
					if i < menu.Days-1 {
						sb.WriteString("\n")
					}
				}
				fullContent := newHeader + "\n" + sb.String()

				newPlan := models.Plan{
					UserID:  uint(userID),
					MenuID:  mID,
					Content: fullContent,
				}
				if err := tx.Create(&newPlan).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Synced successfully"})
}
