# CLAUDE.md

@context/conventions.md

## コマンド

- Docker 起動: `docker compose up --build`（:5000、SQLite は named volume に永続化）
- ビルド（frontend embed 込み単一バイナリ）: `make build`
- テスト: `make test`（in-memory SQLite、外部依存なし）
- 開発: `make dev-dist`（初回スタブ）→ `make dev-back`（:5000）/ `make dev-front`（:5173 HMR）

## アーキテクチャの要点

- 単一 Go バイナリ。frontend のビルド資産を `go:embed`（`backend/dist`）で同梱する。
- 研修メニューの原本は `backend/internal/database/menu_config.json`。起動時に DB へ同期し、管理画面の編集はこの JSON へ書き戻す。DB ファイルは生成物。
- DB は SQLite（Pure Go ドライバ）。

## 検証手段

- backend: `make test`
- frontend: `cd frontend && npm ci && npm run build`（ビルド成功＝型/参照の静的確認）
