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

func TestSaveReport_Creates(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, "🦊")
	h := &ReportsHandler{DB: db}
	router := gin.New()
	router.POST("/api/reports", h.SaveReport)

	body := `{"date": "2025-06-01", "content": "今日の日報"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/reports", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", fmt.Sprintf("%d", user.ID))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var report models.Report
	db.Where("user_id = ? AND date = ?", user.ID, "2025-06-01").First(&report)
	if report.Content != "今日の日報" {
		t.Errorf("expected '今日の日報', got %q", report.Content)
	}
}

func TestSaveReport_Upserts(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, "🐼")
	h := &ReportsHandler{DB: db}
	router := gin.New()
	router.POST("/api/reports", h.SaveReport)

	sendReport := func(content string) {
		body := fmt.Sprintf(`{"date": "2025-06-01", "content": %q}`, content)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/reports", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-Id", fmt.Sprintf("%d", user.ID))
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
	}

	sendReport("最初の日報")
	sendReport("更新された日報")

	var reports []models.Report
	db.Where("user_id = ? AND date = ?", user.ID, "2025-06-01").Find(&reports)
	if len(reports) != 1 {
		t.Errorf("expected 1 report (upsert), got %d", len(reports))
	}
	if reports[0].Content != "更新された日報" {
		t.Errorf("expected updated content, got %q", reports[0].Content)
	}
}

func TestSaveReport_NoAuth_Returns401(t *testing.T) {
	db := setupTestDB(t)
	h := &ReportsHandler{DB: db}
	router := gin.New()
	router.POST("/api/reports", h.SaveReport)

	body := `{"date": "2025-06-01", "content": "test"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/reports", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetReports_NoAuth_ReturnsEmpty(t *testing.T) {
	db := setupTestDB(t)
	h := &ReportsHandler{DB: db}
	router := gin.New()
	router.GET("/api/reports", h.GetReports)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/reports", nil)
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

func TestGetReports_OrderedByDateDesc(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, "🐯")
	db.Create(&models.Report{UserID: user.ID, Date: "2025-06-01", Content: "first"})
	db.Create(&models.Report{UserID: user.ID, Date: "2025-06-03", Content: "third"})
	db.Create(&models.Report{UserID: user.ID, Date: "2025-06-02", Content: "second"})

	h := &ReportsHandler{DB: db}
	router := gin.New()
	router.GET("/api/reports", h.GetReports)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/reports", nil)
	req.Header.Set("X-User-Id", fmt.Sprintf("%d", user.ID))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var reports []models.Report
	json.Unmarshal(w.Body.Bytes(), &reports)
	if len(reports) != 3 {
		t.Fatalf("expected 3 reports, got %d", len(reports))
	}
	if reports[0].Date != "2025-06-03" {
		t.Errorf("expected first date 2025-06-03, got %s", reports[0].Date)
	}
}
