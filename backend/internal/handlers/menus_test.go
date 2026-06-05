package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"training-scheduler/internal/models"

	"github.com/gin-gonic/gin"
)

func TestGetMenus_ReturnsAll(t *testing.T) {
	db := setupTestDB(t)
	createTestMenu(t, db, "Go基礎", 5)
	createTestMenu(t, db, "Git入門", 3)

	h := &MenusHandler{DB: db}
	router := gin.New()
	router.GET("/api/menus", h.GetMenus)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/menus", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var menus []models.Menu
	json.Unmarshal(w.Body.Bytes(), &menus)
	if len(menus) != 2 {
		t.Errorf("expected 2 menus, got %d", len(menus))
	}
}

func TestSaveSelection_CreatesNewPlans(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, "🦅")
	menu := createTestMenu(t, db, "Go基礎", 3)

	h := &MenusHandler{DB: db}
	router := gin.New()
	router.POST("/api/menus/select", h.SaveSelection)

	body := fmt.Sprintf(`{"menu_ids": [%d]}`, menu.ID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/menus/select", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", fmt.Sprintf("%d", user.ID))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var plan models.Plan
	db.Where("user_id = ? AND menu_id = ?", user.ID, menu.ID).First(&plan)
	if plan.ID == 0 {
		t.Fatal("plan was not created")
	}
	// ヘッダー行確認
	if !strings.HasPrefix(plan.Content, "【Go基礎 研修計画（計3日間）】") {
		t.Errorf("unexpected header: %q", plan.Content)
	}
	// 日付テンプレート確認
	if !strings.Contains(plan.Content, "1日目：") || !strings.Contains(plan.Content, "3日目：") {
		t.Errorf("missing day templates: %q", plan.Content)
	}
}

// SaveSelection は既存プランのヘッダー1行目だけを更新し、ボディ内容を保持する
func TestSaveSelection_UpdatesHeaderOnly(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, "🦉")
	menu := createTestMenu(t, db, "旧メニュー名", 3)

	// ユーザーが書き込んだ内容を持つ既存プラン
	existingContent := "【旧メニュー名 研修計画（計3日間）】\n1日目：実装完了\n2日目：レビュー中\n3日目："
	db.Create(&models.Plan{UserID: user.ID, MenuID: menu.ID, Content: existingContent})

	// メニュー名・日数を変更
	db.Model(&menu).Updates(models.Menu{Name: "新メニュー名", Days: 4})

	h := &MenusHandler{DB: db}
	router := gin.New()
	router.POST("/api/menus/select", h.SaveSelection)

	body := fmt.Sprintf(`{"menu_ids": [%d]}`, menu.ID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/menus/select", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", fmt.Sprintf("%d", user.ID))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var plan models.Plan
	db.Where("user_id = ? AND menu_id = ?", user.ID, menu.ID).First(&plan)

	if !strings.HasPrefix(plan.Content, "【新メニュー名 研修計画（計4日間）】") {
		t.Errorf("header not updated: %q", plan.Content)
	}
	if !strings.Contains(plan.Content, "1日目：実装完了") {
		t.Errorf("body content was overwritten: %q", plan.Content)
	}
}

// 選択解除されたメニューのプランは削除される
func TestSaveSelection_DeletesDeselectedPlans(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, "🐬")
	menu1 := createTestMenu(t, db, "Menu1", 2)
	menu2 := createTestMenu(t, db, "Menu2", 2)

	db.Create(&models.Plan{UserID: user.ID, MenuID: menu1.ID, Content: "plan1"})
	db.Create(&models.Plan{UserID: user.ID, MenuID: menu2.ID, Content: "plan2"})

	h := &MenusHandler{DB: db}
	router := gin.New()
	router.POST("/api/menus/select", h.SaveSelection)

	// menu1 だけ選択 → menu2 は削除される
	body := fmt.Sprintf(`{"menu_ids": [%d]}`, menu1.ID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/menus/select", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", fmt.Sprintf("%d", user.ID))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var count int64
	db.Model(&models.Plan{}).Where("user_id = ? AND menu_id = ?", user.ID, menu2.ID).Count(&count)
	if count != 0 {
		t.Error("deselected plan was not deleted")
	}

	db.Model(&models.Plan{}).Where("user_id = ? AND menu_id = ?", user.ID, menu1.ID).Count(&count)
	if count != 1 {
		t.Error("selected plan should still exist")
	}
}

func TestSaveSelection_NoAuth_Returns401(t *testing.T) {
	db := setupTestDB(t)
	h := &MenusHandler{DB: db}
	router := gin.New()
	router.POST("/api/menus/select", h.SaveSelection)

	body := `{"menu_ids": [1]}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/menus/select", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
