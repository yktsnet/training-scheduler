package handlers

import (
	"testing"
	"training-scheduler/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(
		&models.User{},
		&models.Menu{},
		&models.Plan{},
		&models.Report{},
		&models.Progress{},
	); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}
	return db
}

func createTestMenu(t *testing.T, db *gorm.DB, name string, days int) models.Menu {
	t.Helper()
	menu := models.Menu{Name: name, Days: days}
	if err := db.Create(&menu).Error; err != nil {
		t.Fatalf("failed to create test menu: %v", err)
	}
	return menu
}

func createTestUser(t *testing.T, db *gorm.DB, emoji string) models.User {
	t.Helper()
	user := models.User{Emoji: emoji}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return user
}
