package database

import (
	_ "embed"
	"encoding/json"
	"os"
	"training-scheduler/internal/models"

	"gorm.io/gorm"
)

//go:embed menu_config.json
var menuConfigJSON []byte

// SeedMenus は外部JSONファイルを優先し、なければ埋め込みデータを使用してDBを同期します
func SeedMenus(db *gorm.DB) error {
	var menuItems []models.Menu
	var data []byte

	// 1. 外部ファイルを探す
	externalPath := "internal/database/menu_config.json"
	externalData, err := os.ReadFile(externalPath)

	if err == nil {
		data = externalData
	} else {
		data = menuConfigJSON
	}

	if err := json.Unmarshal(data, &menuItems); err != nil {
		return err
	}

	// --- 修正箇所：削除ロジックの追加 ---
	// JSONにあるIDのリストを作成
	var jsonIDs []uint
	for _, item := range menuItems {
		jsonIDs = append(jsonIDs, item.ID)
	}

	// JSONに含まれないIDのメニューをDBから削除
	if err := db.Where("id NOT IN ?", jsonIDs).Delete(&models.Menu{}).Error; err != nil {
		return err
	}
	// --------------------------------

	// 既存データの更新・新規作成 (Upsert)
	for _, item := range menuItems {
		var existing models.Menu
		result := db.First(&existing, item.ID)

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&item).Error; err != nil {
				return err
			}
		} else if result.Error == nil {
			// 既存IDがある場合は、JSONの内容で全フィールドを更新
			if err := db.Save(&item).Error; err != nil {
				return err
			}
		} else {
			return result.Error
		}
	}
	return nil
}
