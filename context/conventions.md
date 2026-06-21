# Conventions

命名・コード規約・スタイル（どう書くか）。

## Go (backend)

- パッケージ構成は `internal/` 配下。`handlers`（HTTP）/ `models`（GORM 構造体）/ `database`（接続・seeder）に分ける。
- ハンドラはリソース単位で1ファイル（`users.go` / `plans.go` / `reports.go` / `menus.go` / `overview.go` / `admin_handlers.go`）。
- モデルのフィールドは Go 側 PascalCase、JSON タグは snake_case（例: `OffsetDays float64 json:"offset_days"`）。GORM 制約はタグで宣言する。
- 複合ユニークなど非自明な制約はフィールド直上にコメントで意図を書く（例: Report の `uniqueIndex:idx_report_user_date`）。
- テストは同パッケージの `*_test.go`。共通セットアップは `testutil_test.go` に置き、in-memory SQLite を使う（外部依存なし）。

## Vue (frontend)

- 画面は `src/views/{Name}View.vue`、再利用部品は `src/components/` に PascalCase。
- ルーティングは `src/router/index.js` に集約。
- API レスポンスの snake_case をそのまま受ける（backend の JSON タグに一致させる）。

## データの原本

- 研修メニューの原本は `backend/internal/database/menu_config.json`。起動時に DB へ同期し、管理画面の編集はこの JSON にも書き戻す。DB ファイル自体は生成物として ignore する。

## コミット

- conventional commits（`feat` / `fix` / `docs` / `style` / `chore`）。
