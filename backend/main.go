package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"
	"os"
	"training-scheduler/internal/database"
	"training-scheduler/internal/handlers"
	"training-scheduler/internal/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// frontend/dist の中身をバイナリに埋め込む
// ※実行前に frontend/dist を backend/dist にコピーしておく必要があります
//go:embed dist/*
var frontendFS embed.FS

func main() {
	// 1. データベース接続 (SQLite)
	db, err := gorm.Open(sqlite.Open("instance/database.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 2. マイグレーションの実行
	err = db.AutoMigrate(&models.User{}, &models.Menu{}, &models.Plan{}, &models.Report{}, &models.Progress{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 3. マスターデータの初期投入/同期 または デモモード初期化
	isDemoMode := os.Getenv("DEMO_MODE") == "true"
	if isDemoMode {
		log.Println("Demo mode is ENABLED. Seeding demo data...")
		if err := database.SeedDemoData(db); err != nil {
			log.Fatalf("failed to seed demo data: %v", err)
		}
		// 30分ごとの自動リセット監視 Goroutine を開始
		go func() {
			ticker := time.NewTicker(30 * time.Minute)
			for range ticker.C {
				log.Println("Checking database drift for demo reset...")
				if database.HasDemoDrift(db) {
					log.Println("Drift detected. Resetting database to initial demo state...")
					if err := database.SeedDemoData(db); err != nil {
						log.Printf("error during demo reset: %v\n", err)
					}
				} else {
					log.Println("No drift detected in demo database.")
				}
			}
		}()
	} else {
		if err := database.SeedMenus(db); err != nil {
			log.Fatalf("failed to seed menus: %v", err)
		}
	}

	// 4. Ginルーターの初期化
	r := gin.Default()

	// 5. CORSの設定 (開発環境用)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "X-User-Id", "X-Admin-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ハンドラーのインスタンス化
	menusHandler := &handlers.MenusHandler{DB: db}
	usersHandler := &handlers.UsersHandler{DB: db}
	plansHandler := &handlers.PlansHandler{DB: db}
	reportsHandler := &handlers.ReportsHandler{DB: db}
	overviewHandler := &handlers.OverviewHandler{DB: db}
	adminHandler := &handlers.AdminHandler{DB: db}

	// 6. ルーティングの登録
	// APIグループ化
	api := r.Group("/api")
	{
		api.GET("/menus", menusHandler.GetMenus)
		api.POST("/menus/select", menusHandler.SaveSelection)

		api.GET("/users", usersHandler.GetUsers)
		api.POST("/users", usersHandler.CreateUser)
		api.DELETE("/users/:id", usersHandler.DeleteUser)

		api.GET("/plans", plansHandler.GetPlans)
		api.POST("/plans/:id", plansHandler.UpdatePlan)

		api.GET("/reports", reportsHandler.GetReports)
		api.POST("/reports", reportsHandler.SaveReport)

		api.GET("/overview", overviewHandler.GetOverview)
		api.POST("/overview", overviewHandler.UpdateOverview)

		// 管理者用ルート
		api.POST("/admin/login", adminHandler.Login)
		admin := api.Group("/admin")
		admin.Use(handlers.AdminAuthMiddleware())
		{
			admin.GET("/menus", adminHandler.GetMenus)
			admin.POST("/menus", adminHandler.CreateMenu)
			admin.PUT("/menus/:id", adminHandler.UpdateMenu)
			admin.DELETE("/menus/:id", adminHandler.DeleteMenu)
		}
	}

	// 7. 静的ファイルとSPAルーティングの設定
	// embedしたFSから "dist" ディレクトリをルートとして取り出す
	publicFS, _ := fs.Sub(frontendFS, "dist")
	staticServer := http.FileServer(http.FS(publicFS))

	// API以外のリクエストをフロントエンドに振り分ける
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// /api で始まるリクエストがここに来た場合は、純粋な404を返す
		if strings.HasPrefix(path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API route not found"})
			return
		}

		// 静的ファイル（js, css, png等）が存在するか確認
		_, err := publicFS.Open(strings.TrimPrefix(path, "/"))
		if err == nil {
			// ファイルが存在すればそれを返す
			staticServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		// ファイルが存在しない、またはルートへのアクセスは index.html を返す (Vue Router用)
		indexData, _ := frontendFS.ReadFile("dist/index.html")
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexData)
	})

	// 5000番ポートでサーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Printf("Server starting on :%s (Serving API and Frontend)...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
