## 変更内容

マルチステージ Dockerfile と compose.yaml を追加し、`docker compose up` で本番相当のシングルバイナリが起動できるようにした。

- **`.dockerignore`**: `node_modules/`・`backend/dist/`・`.git/` 等のビルド不要物を除外
- **`Dockerfile`**（マルチステージ）:
  - Stage 1 `frontend-build` (`node:22-alpine`): `npm ci && npm run build` で `dist/` を生成
  - Stage 2 `backend-build` (`golang:1.25-alpine`): Stage 1 の dist を `backend/dist/` へ COPY、`CGO_ENABLED=0 go build -o training-app .`
  - Stage 3 `runtime` (`gcr.io/distroless/static-debian12`): バイナリのみ配置、`EXPOSE 5000`
- **`compose.yaml`**: `app` サービス 1 つ、ポート `5000:5000`、named volume `db-data` を `/data` にマウント、`DB_PATH=/data/training.db` で SQLite を永続化
- **`backend/internal/database/seeder.go`**: `GetDatabasePath()` に `DB_PATH` 環境変数サポートを追加（issue 備考「対応が必要なら最小限の変更を行う」に基づく）

## 静的確認結果

- `make test` → 全 28 件パス
- Dockerfile のビルドステージが Makefile の `build` ターゲットと同等であることを確認
  - Makefile: `npm ci && npm run build` → `cp -r frontend/dist backend/dist` → `go build -o training-app .`
  - Dockerfile Stage1: `npm ci && npm run build` ✅ Stage2: `COPY --from=frontend-build dist/` + `CGO_ENABLED=0 go build -o training-app .` ✅
- `GetDatabasePath()` は `os` を既にインポート済みのため import 追加不要
- `git diff --name-only HEAD~1` の出力:
  ```
  .dockerignore
  Dockerfile
  backend/internal/database/seeder.go
  compose.yaml
  ```

## 検証手順

```bash
docker compose up --build
# ブラウザで http://localhost:5000 を確認
```
