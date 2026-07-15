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

	body := `{"emoji": "🦁", "initial": "yt"}`
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
	if result.Initial != "YT" {
		t.Errorf("expected initial YT (uppercased), got %s", result.Initial)
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

	body := `{"emoji": "🐰", "initial": "BB"}`

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

func TestCreateUser_InitialValidation(t *testing.T) {
	db := setupTestDB(t)
	h := &UsersHandler{DB: db}
	router := gin.New()
	router.POST("/api/users", h.CreateUser)

	tests := []struct {
		name       string
		body       string
		expectedCode int
	}{
		{
			name:       "empty initial",
			body:       `{"emoji": "🦊", "initial": ""}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:       "too long initial",
			body:       `{"emoji": "🦊", "initial": "ABCD"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:       "whitespace initial",
			body:       `{"emoji": "🦊", "initial": "   "}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:       "valid 3 chars",
			body:       `{"emoji": "🦊", "initial": "abc"}`,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			if w.Code != tc.expectedCode {
				t.Errorf("expected %d, got %d: %s", tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestCreateUser_EmojiRequired(t *testing.T) {
	db := setupTestDB(t)
	h := &UsersHandler{DB: db}
	router := gin.New()
	router.POST("/api/users", h.CreateUser)

	body := `{"initial": "yt"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing emoji, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDeleteUser_NonexistentID_NoOp(t *testing.T) {
	db := setupTestDB(t)
	h := &UsersHandler{DB: db}
	router := gin.New()
	router.DELETE("/api/users/:id", h.DeleteUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/users/9999", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 (no-op) for nonexistent user, got %d: %s", w.Code, w.Body.String())
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
