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

func TestGetUsers_Empty(t *testing.T) {
	db := setupTestDB(t)
	h := &UsersHandler{DB: db}
	router := gin.New()
	router.GET("/api/users", h.GetUsers)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var result []models.User
	json.Unmarshal(w.Body.Bytes(), &result)
	if len(result) != 0 {
		t.Errorf("expected empty list, got %d items", len(result))
	}
}

func TestCreateUser_Success(t *testing.T) {
	db := setupTestDB(t)
	h := &UsersHandler{DB: db}
	router := gin.New()
	router.POST("/api/users", h.CreateUser)

	body := `{"emoji": "🦁"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var result models.User
	json.Unmarshal(w.Body.Bytes(), &result)
	if result.Emoji != "🦁" {
		t.Errorf("expected emoji 🦁, got %s", result.Emoji)
	}
	if result.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestCreateUser_DuplicateEmoji(t *testing.T) {
	db := setupTestDB(t)
	h := &UsersHandler{DB: db}
	router := gin.New()
	router.POST("/api/users", h.CreateUser)

	body := `{"emoji": "🐰"}`

	// 1回目: 成功
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("first create failed: %s", w.Body.String())
	}

	// 2回目: 重複エラー
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for duplicate, got %d", w.Code)
	}
}

func TestDeleteUser_CascadesRelatedData(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, "🐻")
	menu := createTestMenu(t, db, "Test Menu", 3)

	db.Create(&models.Plan{UserID: user.ID, MenuID: menu.ID, Content: "test plan"})
	db.Create(&models.Report{UserID: user.ID, Date: "2025-01-01", Content: "test report"})
	db.Create(&models.Progress{UserID: user.ID, MenuID: menu.ID})

	h := &UsersHandler{DB: db}
	router := gin.New()
	router.DELETE("/api/users/:id", h.DeleteUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/users/%d", user.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var planCount, reportCount, progressCount int64
	db.Model(&models.Plan{}).Where("user_id = ?", user.ID).Count(&planCount)
	db.Model(&models.Report{}).Where("user_id = ?", user.ID).Count(&reportCount)
	db.Model(&models.Progress{}).Where("user_id = ?", user.ID).Count(&progressCount)

	if planCount != 0 {
		t.Errorf("plans not deleted: got %d", planCount)
	}
	if reportCount != 0 {
		t.Errorf("reports not deleted: got %d", reportCount)
	}
	if progressCount != 0 {
		t.Errorf("progress not deleted: got %d", progressCount)
	}
}
