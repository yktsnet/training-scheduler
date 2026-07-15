package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"training-scheduler/internal/models"

	"github.com/gin-gonic/gin"
)

// テスト用のクリーンアップ（テスト中に作られる一時的なJSONファイルを削除）
func cleanupTestJSON(t *testing.T) {
	t.Helper()
	// テスト時にカレントディレクトリ(backend/internal/handlers)内に作られる可能性のある
	// "internal/database/menu_config.json" を削除する
	path := "internal/database/menu_config.json"
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
		// 親ディレクトリも空なら削除
		os.Remove("internal/database")
		os.Remove("internal")
	}
}

func TestAdminLogin_Success(t *testing.T) {
	db := setupTestDB(t)
	h := &AdminHandler{DB: db}

	router := gin.New()
	router.POST("/api/admin/login", h.Login)

	// ADMIN_PASSWORD 環境変数を設定 (テスト終了時にクリーンアップ)
	origPassword := os.Getenv("ADMIN_PASSWORD")
	os.Setenv("ADMIN_PASSWORD", "testpass123")
	defer func() {
		if origPassword == "" {
			os.Unsetenv("ADMIN_PASSWORD")
		} else {
			os.Setenv("ADMIN_PASSWORD", origPassword)
		}
	}()

	body := `{"password": "testpass123"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] == "" {
		t.Error("expected token in response, got empty")
	}
}

func TestAdminLogin_Failure(t *testing.T) {
	db := setupTestDB(t)
	h := &AdminHandler{DB: db}

	router := gin.New()
	router.POST("/api/admin/login", h.Login)

	// 無効なパスワード
	body := `{"password": "wrongpassword"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAdminLogin_DefaultPasswordWhenEnvUnset(t *testing.T) {
	db := setupTestDB(t)
	h := &AdminHandler{DB: db}

	router := gin.New()
	router.POST("/api/admin/login", h.Login)

	origPassword, hadPassword := os.LookupEnv("ADMIN_PASSWORD")
	os.Unsetenv("ADMIN_PASSWORD")
	defer func() {
		if hadPassword {
			os.Setenv("ADMIN_PASSWORD", origPassword)
		}
	}()

	body := `{"password": "admin123"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 with default password when ADMIN_PASSWORD unset, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAdminAuthMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(AdminAuthMiddleware())
	router.GET("/api/admin/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 1. ヘッダーなし (401)
	w1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/api/admin/test", nil)
	router.ServeHTTP(w1, req1)
	if w1.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w1.Code)
	}

	// 2. 正しいトークン (200)
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/api/admin/test", nil)
	token := getAdminToken()
	req2.Header.Set("X-Admin-Token", token)
	router.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w2.Code)
	}
}

func TestAdminCreateMenu_Success(t *testing.T) {
	defer cleanupTestJSON(t)
	db := setupTestDB(t)
	h := &AdminHandler{DB: db}

	router := gin.New()
	router.POST("/api/admin/menus", h.CreateMenu)

	body := `{"name": "Docker入門", "days": 3, "description": "Dockerの基礎", "url": "https://example.com"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/menus", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var menu models.Menu
	json.Unmarshal(w.Body.Bytes(), &menu)
	if menu.ID == 0 || menu.Name != "Docker入門" {
		t.Errorf("invalid created menu: %+v", menu)
	}

	// DBに保存されているか確認
	var count int64
	db.Model(&models.Menu{}).Where("name = ?", "Docker入門").Count(&count)
	if count != 1 {
		t.Error("menu was not saved to DB")
	}
}

func TestAdminUpdateMenu_Success(t *testing.T) {
	defer cleanupTestJSON(t)
	db := setupTestDB(t)
	menu := createTestMenu(t, db, "Go基礎", 3)

	h := &AdminHandler{DB: db}
	router := gin.New()
	router.PUT("/api/admin/menus/:id", h.UpdateMenu)

	body := fmt.Sprintf(`{"name": "Go基礎(改)", "days": 4, "description": "Goの応用", "url": "https://example.com"}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/admin/menus/%d", menu.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var updated models.Menu
	db.First(&updated, menu.ID)
	if updated.Name != "Go基礎(改)" || updated.Days != 4 {
		t.Errorf("menu was not updated correctly: %+v", updated)
	}
}

func TestAdminDeleteMenu_Success(t *testing.T) {
	defer cleanupTestJSON(t)
	db := setupTestDB(t)
	menu := createTestMenu(t, db, "消されるメニュー", 2)
	user := createTestUser(t, db, "🐰")

	// 関連データを作成
	db.Create(&models.Plan{UserID: user.ID, MenuID: menu.ID, Content: "plan"})
	db.Create(&models.Progress{UserID: user.ID, MenuID: menu.ID, OffsetDays: 1})

	h := &AdminHandler{DB: db}
	router := gin.New()
	router.DELETE("/api/admin/menus/:id", h.DeleteMenu)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/admin/menus/%d", menu.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// 削除されたか確認
	var count int64
	db.Model(&models.Menu{}).Where("id = ?", menu.ID).Count(&count)
	if count != 0 {
		t.Error("menu was not deleted")
	}

	// 関連計画・進捗も削除されたか確認
	db.Model(&models.Plan{}).Where("menu_id = ?", menu.ID).Count(&count)
	if count != 0 {
		t.Error("associated plans were not deleted")
	}
	db.Model(&models.Progress{}).Where("menu_id = ?", menu.ID).Count(&count)
	if count != 0 {
		t.Error("associated progresses were not deleted")
	}
}