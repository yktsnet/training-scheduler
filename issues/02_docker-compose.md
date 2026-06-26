## Docker Compose 対応
id: 02
skill: pr-workflow
branch-slug: docker-compose
github_issue: 3
status: open
type: feat
対象: Dockerfile (新規), compose.yaml (新規), .dockerignore (新規)
内容: マルチステージ Dockerfile と compose.yaml を追加し、`docker compose up` で本番相当のシングルバイナリが起動できるようにする。
確認: `make test` が通ること、Dockerfile 内のビルドステージが Makefile の build ターゲットと同等であること
---

## 設計

### Dockerfile（マルチステージ）

1. **Stage 1: frontend-build** — `node:22-alpine`
   - `frontend/` を COPY し `npm ci && npm run build` で `dist/` を生成
2. **Stage 2: backend-build** — `golang:1.25-alpine`
   - `backend/` を COPY し、Stage 1 の `dist/` を `backend/dist/` へ COPY
   - `CGO_ENABLED=0 go build -o training-app .`
3. **Stage 3: runtime** — `gcr.io/distroless/static-debian12`
   - Stage 2 のバイナリだけ COPY
   - `EXPOSE 5000`、`ENTRYPOINT ["/training-app"]`

### compose.yaml

- サービス1つ（`app`）
- ビルドコンテキスト `.`、Dockerfile `Dockerfile`
- ポート `5000:5000`
- SQLite データ永続化: named volume を `/data` にマウントし、環境変数 `DB_PATH=/data/training.db` で指定

### .dockerignore

- `node_modules/`, `backend/dist/`, `backend/training-app`, `instance/`, `.git/`, `issues/`, `context/`, `src/`（スクリーンショット）

## 備考

- SQLite は Pure Go ドライバ（`glebarez/sqlite`）なので CGO 不要。distroless で動く。
- DB パス環境変数の受け取りは既存コードの確認が必要。対応が必要なら本 Issue 内で最小限の変更を行う。

## 実装順序

1. .dockerignore
2. Dockerfile
3. compose.yaml
4. （必要なら）DB パス環境変数対応
