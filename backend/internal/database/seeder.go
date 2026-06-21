package database

import (
	_ "embed"
	"encoding/json"
	"os"
	"time"
	"training-scheduler/internal/models"

	"gorm.io/gorm"
)

//go:embed menu_config.json
var menuConfigJSON []byte

const PlanContent = `【Gitで始めるバージョン管理 研修計画（計3日間）】
1日目：Git基本操作（完了！✅）
- リポジトリの作成、クローン、add, commit, pushの練習。
- ターミナルでの基本的なコマンド操作を体に覚え込ませる。

2日目：ブランチ管理・コンフリクト解消（完了！✅）
- ブランチの作成、切り替え、マージ（Fast-Forward / 3-way）。
- わざと競合を発生させ、競合マーカーを読み解いて手動で解消する。

3日目：GitHub Flowの実践（月曜日予定・進行中）
- Pull Request of 作成、レビュワーからの指摘修正。
- 本番用 main ブランチへのマージとローカルへのプル。`

const Report1Content = `本日から研修がスタートしました！
初日はGitの基本操作（リポジトリの作成、クローン、コミット、プッシュ）を学びました。
コマンドラインでの操作に少し慣れていない部分があり、プッシュの宛先設定でエラーが出ましたが、先輩にアドバイスをもらって無事にGitHub上に反映させることができました！
明日も集中して取り組みます。`

const Report2Content = `今日はGit研修の2日目でした。
ブランチを使ったマージ作業と、コンフリクト（競合）の解消を実践しました。
複数人で同じファイルを編集したときに競合が起きる仕組みがよく理解できました。
最初はマージ競合の画面に焦りましたが、差分を比較して手動で修正できるようになり、少し自信がつきました！
月曜日はGitHub上のPull Requestを用いたチーム開発フローを体験します。`

// GetMenuConfigPath はプロジェクトルート実行かbackendディレクトリ実行かに応じてmenu_config.jsonのパスを動的に解決します
func GetMenuConfigPath() string {
	if _, err := os.Stat("backend/internal/database/menu_config.json"); err == nil {
		return "backend/internal/database/menu_config.json"
	}
	return "internal/database/menu_config.json"
}

// GetDatabasePath はプロジェクトルート実行かbackendディレクトリ実行かに応じてdatabase.dbのパスを動的に解決します
func GetDatabasePath() string {
	if _, err := os.Stat("backend/internal/database/menu_config.json"); err == nil {
		return "backend/instance/database.db"
	}
	return "instance/database.db"
}

// SeedMenus は外部JSONファイルを優先し、なければ埋め込みデータを使用してDBを同期します
func SeedMenus(db *gorm.DB) error {
	var menuItems []models.Menu
	var data []byte

	// 1. 外部ファイルを探す
	externalPath := GetMenuConfigPath()
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
		if err := db.Limit(1).Find(&existing, item.ID).Error; err != nil {
			return err
		}

		if existing.ID == 0 {
			if err := db.Create(&item).Error; err != nil {
				return err
			}
		} else {
			// 既存IDがある場合は、JSONの内容で全フィールドを更新
			if err := db.Save(&item).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// SeedDemoData はデータベースを一度空にして、現在の日時を基準としたダミーのデモデータを再投入します
func SeedDemoData(db *gorm.DB) error {
	// テーブルの全レコード削除
	db.Exec("DELETE FROM plans")
	db.Exec("DELETE FROM reports")
	db.Exec("DELETE FROM progresses")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM menus")

	// メニューのシード
	if err := SeedMenus(db); err != nil {
		return err
	}

	// デモユーザー (🐶) の作成
	demoUser := models.User{ID: 1, Emoji: "🐶"}
	if err := db.Create(&demoUser).Error; err != nil {
		return err
	}

	// 動的な日付計算 (今日から3日前、2日前)
	today := time.Now().Local()
	format := "2006-01-02"
	dayMinus3 := today.AddDate(0, 0, -3).Format(format)
	dayMinus2 := today.AddDate(0, 0, -2).Format(format)

	// デモ研修計画の作成
	demoPlan := models.Plan{
		ID:      1,
		UserID:  1,
		MenuID:  1,
		Content: PlanContent,
	}
	if err := db.Create(&demoPlan).Error; err != nil {
		return err
	}

	// デモ日報の作成 (2日分)
	report1 := models.Report{
		ID:      1,
		UserID:  1,
		Date:    dayMinus3,
		Content: Report1Content,
	}
	report2 := models.Report{
		ID:      2,
		UserID:  1,
		Date:    dayMinus2,
		Content: Report2Content,
	}
	if err := db.Create(&report1).Error; err != nil {
		return err
	}
	if err := db.Create(&report2).Error; err != nil {
		return err
	}

	// デモ進捗データの作成
	demoProgress := models.Progress{
		ID:          1,
		UserID:      1,
		MenuID:      1,
		StartDate:   dayMinus3,
		TargetDays:  3,
		OffsetDays:  3.0,
		StatusMemo:  "2日目まで順調に完了。月曜日に最終日のGitHub Flowを実施して修了予定です。",
		IsCompleted: false,
	}
	if err := db.Create(&demoProgress).Error; err != nil {
		return err
	}

	return nil
}

// HasDemoDrift はデータベースの現在の状態が初期デモ状態から変更されているかをチェックします
func HasDemoDrift(db *gorm.DB) bool {
	var userCount int64
	var planCount int64
	var reportCount int64
	var progressCount int64
	var menuCount int64

	db.Model(&models.User{}).Count(&userCount)
	db.Model(&models.Plan{}).Count(&planCount)
	db.Model(&models.Report{}).Count(&reportCount)
	db.Model(&models.Progress{}).Count(&progressCount)
	db.Model(&models.Menu{}).Count(&menuCount)

	// menu_config.json (または埋め込みデータ) の期待メニュー件数を算出
	var menuItems []models.Menu
	var data []byte
	externalPath := GetMenuConfigPath()
	if externalData, err := os.ReadFile(externalPath); err == nil {
		data = externalData
	} else {
		data = menuConfigJSON
	}
	_ = json.Unmarshal(data, &menuItems)
	expectedMenuCount := int64(len(menuItems))

	if userCount != 1 || planCount != 1 || reportCount != 2 || progressCount != 1 || menuCount != expectedMenuCount {
		return true
	}

	var user models.User
	if err := db.Limit(1).Find(&user).Error; err != nil || user.ID == 0 || user.Emoji != "🐶" {
		return true
	}

	var plan models.Plan
	if err := db.Limit(1).Find(&plan).Error; err != nil || plan.ID == 0 || plan.UserID != 1 || plan.MenuID != 1 || plan.Content != PlanContent {
		return true
	}

	var progress models.Progress
	if err := db.Limit(1).Find(&progress).Error; err != nil ||
		progress.ID == 0 ||
		progress.UserID != 1 ||
		progress.MenuID != 1 ||
		progress.TargetDays != 3 ||
		progress.OffsetDays != 3.0 ||
		progress.StatusMemo != "2日目まで順調に完了。月曜日に最終日のGitHub Flowを実施して修了予定です。" ||
		progress.IsCompleted {
		return true
	}

	today := time.Now().Local()
	format := "2006-01-02"
	dayMinus3 := today.AddDate(0, 0, -3).Format(format)
	dayMinus2 := today.AddDate(0, 0, -2).Format(format)

	var r1 models.Report
	if err := db.Where("date = ?", dayMinus3).Limit(1).Find(&r1).Error; err != nil || r1.ID == 0 || r1.Content != Report1Content {
		return true
	}
	var r2 models.Report
	if err := db.Where("date = ?", dayMinus2).Limit(1).Find(&r2).Error; err != nil || r2.ID == 0 || r2.Content != Report2Content {
		return true
	}

	return false
}
