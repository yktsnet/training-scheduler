package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"training-scheduler/internal/database"
	"training-scheduler/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandler struct {
	DB *gorm.DB
}

// getAdminToken は現在の ADMIN_PASSWORD (未設定時は admin123) の SHA256 ハッシュを返します
func getAdminToken() string {
	pwd := os.Getenv("ADMIN_PASSWORD")
	if pwd == "" {
		pwd = "admin123"
	}
	hash := sha256.Sum256([]byte(pwd))
	return hex.EncodeToString(hash[:])
}

// AdminAuthMiddleware は X-Admin-Token ヘッダーを検証するミドルウェアです
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Admin-Token")
		expectedToken := getAdminToken()

		if token == "" || token != expectedToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized admin access"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// POST /api/admin/login
func (h *AdminHandler) Login(c *gin.Context) {
	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	pwd := os.Getenv("ADMIN_PASSWORD")
	if pwd == "" {
		pwd = "admin123"
	}

	if req.Password != pwd {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// 認証成功時、ハッシュ化したトークンを返す
	c.JSON(http.StatusOK, gin.H{
		"token": getAdminToken(),
	})
}

// GET /api/admin/menus (すべてのメニューを取得)
func (h *AdminHandler) GetMenus(c *gin.Context) {
	var menus []models.Menu
	if err := h.DB.Order("id asc").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, menus)
}

// POST /api/admin/menus (新規メニュー作成)
func (h *AdminHandler) CreateMenu(c *gin.Context) {
	var menu models.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 新規登録のため ID をゼロ値にして自動採番させる
	menu.ID = 0

	if err := h.DB.Create(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// menu_config.json の書き出し
	if err := h.rewriteMenuConfigJSON(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync JSON: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// PUT /api/admin/menus/:id (メニュー更新)
func (h *AdminHandler) UpdateMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
		return
	}

	var req models.Menu
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req.ID = uint(id)

	var existing models.Menu
	if err := h.DB.First(&existing, req.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}

	if err := h.DB.Save(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// menu_config.json の書き出し
	if err := h.rewriteMenuConfigJSON(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync JSON: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// DELETE /api/admin/menus/:id (メニュー削除 ＆ 関連データ削除)
func (h *AdminHandler) DeleteMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
		return
	}

	var existing models.Menu
	if err := h.DB.First(&existing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		// 関連する計画 (plans) の削除
		if err := tx.Where("menu_id = ?", id).Delete(&models.Plan{}).Error; err != nil {
			return err
		}
		// 関連する進捗 (progresses) の削除
		if err := tx.Where("menu_id = ?", id).Delete(&models.Progress{}).Error; err != nil {
			return err
		}
		// メニュー自体の削除
		if err := tx.Delete(&models.Menu{}, id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// menu_config.json の書き出し
	if err := h.rewriteMenuConfigJSON(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync JSON: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// rewriteMenuConfigJSON はDBの最新メニューデータを internal/database/menu_config.json に書き出します
func (h *AdminHandler) rewriteMenuConfigJSON() error {
	var menus []models.Menu
	if err := h.DB.Order("id asc").Find(&menus).Error; err != nil {
		return err
	}

	data, err := json.MarshalIndent(menus, "", "  ")
	if err != nil {
		return err
	}

	externalPath := database.GetMenuConfigPath()
	return os.WriteFile(externalPath, data, 0644)
}
