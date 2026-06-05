package models

// User (アニマルユーザー)
type User struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Emoji string `gorm:"type:varchar(10);unique;not null" json:"emoji"`
}

// Menu (研修メニュー)
type Menu struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	Name          string `gorm:"type:varchar(100);not null" json:"name"`
	Days          int    `gorm:"not null" json:"days"`
	DocLink       string `gorm:"type:varchar(255)" json:"doc_link"`
	Summary       string `gorm:"type:text" json:"summary"`
	Skills        string `gorm:"type:varchar(200)" json:"skills"`
	Prerequisites string `gorm:"type:varchar(200)" json:"prerequisites"`
	Difficulty    int    `json:"difficulty"`
}

// Plan (研修計画)
type Plan struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	UserID  uint   `gorm:"not null" json:"user_id"`
	MenuID  uint   `gorm:"not null" json:"menu_id"`
	Content string `gorm:"type:text" json:"content"`
	Menu    Menu   `gorm:"foreignKey:MenuID" json:"menu,omitempty"`
}

// Report (日報)
// UserID + Date の複合ユニーク制約 → OnConflict upsert に必要
type Report struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	UserID  uint   `gorm:"not null;uniqueIndex:idx_report_user_date" json:"user_id"`
	Date    string `gorm:"type:varchar(10);uniqueIndex:idx_report_user_date" json:"date"`
	Content string `gorm:"type:text" json:"content"`
}

// Progress (進捗状況)
type Progress struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	UserID      uint    `gorm:"not null" json:"user_id"`
	MenuID      uint    `gorm:"not null" json:"menu_id"`
	StartDate   string  `gorm:"type:varchar(10)" json:"start_date"`
	TargetDays  int     `gorm:"default:7" json:"target_days"`
	OffsetDays  float64 `gorm:"default:0.0" json:"offset_days"`
	StatusMemo  string  `gorm:"type:text" json:"status_memo"`
	IsCompleted bool    `gorm:"default:false" json:"is_completed"`
	Menu        Menu    `gorm:"foreignKey:MenuID" json:"menu"`
}
