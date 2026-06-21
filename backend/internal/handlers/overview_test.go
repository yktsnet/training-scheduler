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

func TestGetOverview_FallsBackToMenuDays(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, "🦁")
	menu := createTestMenu(t, db, "Vue入門", 7)
	db.Create(&models.Plan{UserID: user.ID, MenuID: menu.ID, Content: "計画"})
	// Progress を作らない → TargetDays は Menu.Days に fallback する

	h := &OverviewHandler{DB: db}
	router := gin.New()
	router.GET("/api/overview", h.GetOverview)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/overview", nil)
	req.Header.Set("X-User-Id", fmt.Sprintf("%d", user.ID))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var items []struct {
		MenuName   string `json:"menu_name"`
		TargetDays int    `json:"target_days"`
	}
	json.Unmarshal(w.Body.Bytes(), &items)
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].TargetDays != 7 {
		t.Errorf("expected target_days to fall back to menu days 7, got %d", items[0].TargetDays)
	}
}

func TestUpdateOverview_UpsertsProgress(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, "🐨")
	menu := createTestMenu(t, db, "テスト設計", 4)

	h := &OverviewHandler{DB: db}
	router := gin.New()
	router.POST("/api/overview", h.UpdateOverview)

	send := func(offset float64) {
		body := fmt.Sprintf(`{"menu_id": %d, "offset_days": %v, "status_memo": "順調"}`, menu.ID, offset)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/overview", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-Id", fmt.Sprintf("%d", user.ID))
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
	}

	send(1.0)
	send(3.0) // 同一 user+menu は upsert され重複しない

	var rows []models.Progress
	db.Where("user_id = ? AND menu_id = ?", user.ID, menu.ID).Find(&rows)
	if len(rows) != 1 {
		t.Fatalf("expected 1 progress row (upsert), got %d", len(rows))
	}
	if rows[0].OffsetDays != 3.0 {
		t.Errorf("expected offset_days 3.0, got %v", rows[0].OffsetDays)
	}
}

func TestUpdateOverview_NoAuth_Returns401(t *testing.T) {
	db := setupTestDB(t)
	h := &OverviewHandler{DB: db}
	router := gin.New()
	router.POST("/api/overview", h.UpdateOverview)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/overview", bytes.NewBufferString(`{"menu_id":1}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
