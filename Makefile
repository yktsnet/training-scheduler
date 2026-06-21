.PHONY: build test clean dev-back dev-front dev-dist

# 単一バイナリのビルド（frontend embed込み）
# 実行順序: frontend build → dist コピー → go build
build:
	cd frontend && npm ci && npm run build
	rm -rf backend/dist
	cp -r frontend/dist backend/dist
	cd backend && go build -o training-app .

# テスト実行（追加依存なし・標準ライブラリのみ）
test:
	cd backend && go test ./internal/... -v -count=1

# バイナリ・distを削除
clean:
	rm -f backend/training-app
	rm -rf backend/dist

# 開発用: フロントエンド（Vite dev server、ポート5173）
dev-front:
	cd frontend && npm run dev

# 開発用: バックエンド単体起動（ポート5000）
# ※ go:embed のため backend/dist が必要。初回は make dev-dist を先に実行
dev-back:
	cd backend && go run main.go

# 開発用: backend/dist のスタブを作成（go run を通すための最低限）
dev-dist:
	mkdir -p backend/dist
	echo '<!DOCTYPE html><html><body>dev mode - open http://localhost:5173</body></html>' > backend/dist/index.html
