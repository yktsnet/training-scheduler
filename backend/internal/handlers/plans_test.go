package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"training-scheduler/internal/models"

	"github.com/gin-gonic/gin"
)

func TestGetPlans_ReturnsWithMenuName(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, "🦊")
	menu := createTestMenu(t, db, "Go入門", 5)
	db.Create(&models.Plan{UserID: user.ID, MenuID: menu.ID, Content: "標準ライブラリを読む"})

	h := &PlansHandler{DB: db}
	router := gin.New()
	router.GET("/api/plans", h.GetPlans)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/plans", nil)
	req.Header.Set("X-User-Id", fmt.Sprintf("%d", user.ID))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var result []struct {
		MenuName string `json:"menu_name"`
		Content  string `json:"content"`
	}
	json.Unmarshal(w.Body.Bytes(), &result)
	if len(result) != 1 {
		t.Fatalf("expected 1 plan, got %d", len(result))
	}
	if result[0].MenuName != "Go入門" {
		t.Errorf("expected menu_name 'Go入門', got %q", result[0].MenuName)
	}
}

func TestGetPlans_NoAuth_ReturnsEmpty(t *testing.T) {
	db := setupTestDB(t)
	h := &PlansHandler{DB: db}
	router := gin.New()
	router.GET("/api/plans", h.GetPlans)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/plans", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var result []interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	if len(result) != 0 {
		t.Errorf("expected empty without auth, got %d items", len(result))
	}
}

func TestUpdatePlan_NonexistentID_Returns404(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, "🐿")

	h := &PlansHandler{DB: db}
	router := gin.New()
	router.POST("/api/plans/:id", h.UpdatePlan)

	body := `{"content": "更新"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/plans/9999", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", fmt.Sprintf("%d", user.ID))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for nonexistent plan id, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdatePlan_OnlyOwnerCanUpdate(t *testing.T) {
	db := setupTestDB(t)
	owner := createTestUser(t, db, "🐼")
	other := createTestUser(t, db, "🐯")
	menu := createTestMenu(t, db, "DB設計", 3)
	plan := models.Plan{UserID: owner.ID, MenuID: menu.ID, Content: "初期"}
	db.Create(&plan)

	h := &PlansHandler{DB: db}
	router := gin.New()
	router.POST("/api/plans/:id", h.UpdatePlan)

	// 他人が更新しようとすると 404（RowsAffected == 0）で、内容は変わらない
	body := `{"content": "改ざん"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/plans/%d", plan.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", fmt.Sprintf("%d", other.ID))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for non-owner, got %d", w.Code)
	}
	var fresh models.Plan
	db.First(&fresh, plan.ID)
	if fresh.Content != "初期" {
		t.Errorf("non-owner must not modify content, got %q", fresh.Content)
	}
}
