# 研修スケジュール管理 Webアプリ (Training Scheduler)

[![CI](https://github.com/yktsnet/training-scheduler/actions/workflows/ci.yml/badge.svg)](https://github.com/yktsnet/training-scheduler/actions/workflows/ci.yml)

新入社員の**自律性を促す**ことを目的とした研修支援ツールです。「システムによる自動管理」と「手書き感覚のアナログ操作」を融合させ、ガチガチの進捗管理ではなく、新人の「主観的な手応え」をベースにメンターが静かに見守るためのアプリケーションです。

---

## 💡 コンセプト

- **主体的プランニング**: システムは枠組みだけを提示し、具体的な計画は新人が自身の言葉で記述します。
- **内省の可視化**: 機械的な進捗率（％）の計算ではなく、本人の「主観的なズレ（手応え）」をマネージャと共有します。
- **非干渉の監視**: マネージャは新人の自律を妨げず、ダッシュボードから状況を静かに見守り、必要な時だけサポートに入ります。
- **ゆるやかな識別（アニマル・ログイン）**: パスワード等による厳格な認証ではなく、動物の絵文字を選ぶだけのシンプルなログインを採用。チーム内の信頼関係を前提とした、遊び心のあるアカウント管理です。

---

## 🔒 セキュリティについて

本アプリはチーム内の**信頼関係を前提とした小規模利用**を想定して設計されています。

- ログインは絵文字の選択のみで、パスワード認証はありません
- 他ユーザーのデータへの書き込みはサーバー側で防止していますが、読み取りは制限していません
- **インターネットに公開する場合は、Cloudflare Access 等による IP 制限・アクセス制御を別途設けることを強く推奨します**

---

## 🗄 データ構造 (Models)

### 0. User (アニマルユーザー)

<img src="src/animals.png" width="500" alt="menu-pic">

- **役割**: アプリを利用する個人（新人・メンター）の識別。
- **項目**: `emoji` (🦁や🐰などのユニークな絵文字)。

### 1. Menu (研修メニュー)

<img src="src/menu.png" width="500" alt="menu-pic">

- **役割**: 研修カリキュラムのマスターデータ（全ユーザー共通）。
- **項目**: 名称、目安日数、概要、参考URL。
- ※ `internal/database/menu_config.json` をマスターとして起動時に自動同期します。

### 2. Plan (研修計画)

<img src="src/plan.png" width="500" alt="plan-pic">

- **役割**: 各メニューに対する具体的な学習計画。
- **項目**: `content` (自由記述のテキスト)、`user_id`。

### 3. Report (日報)

<img src="src/daily.png" width="500" alt="daily-pic">

- **役割**: 日付単位の事実と内省の記録。
- **項目**: `date` (YYYY-MM-DD)、`content` (日報内容)、`user_id`。

### 4. Progress (進捗状況)

<img src="src/overview.png" width="500" alt="overview-pic">

- **役割**: ダッシュボード表示用のメタ情報。
- **項目**: 開始日、目標日数、`offset_days` (主観ズレ値 1〜5)、ステータスメモ。

---

## 🛠 技術スタック

**Frontend**
- Vue 3 (Composition API)
- Vite / Vue Router
- Axios, date-fns

**Backend**
- Go 1.25+
- Gin (Web Framework)
- GORM (ORM) / SQLite (Pure Go driver)
- go:embed (フロントエンド資産をバイナリに内包)

---

## 📦 セットアップ

### 前提条件

- Go 1.25+
- Node.js 20+

### 単一バイナリのビルドと起動

フロントエンドのビルド → Go バイナリへの embed → 起動を一連で行います。

```bash
# ビルド（frontend + backend を単一バイナリに）
make build

# 起動
./backend/training-app
```

ブラウザで http://localhost:5000 にアクセス。

### 開発モード（フロント・バック分離）

```bash
# 初回のみ: go:embed 用のスタブ作成
make dev-dist

# ターミナル1: バックエンド (port 5000)
make dev-back

# ターミナル2: フロントエンド (port 5173, HMR有効)
make dev-front
```

### テスト実行

```bash
make test
```

外部依存なし（in-memory SQLite）。Go の標準ライブラリのみで動作します。

---

## 🚀 CI/CD

### CI（自動テスト・ビルド確認）

`main` / `go-dev` への push および pull request 時に GitHub Actions で自動実行されます。

1. Go テスト (`go test ./internal/...`)
2. フロントエンドビルド + Go バイナリビルドの疎通確認

### CD（自動デプロイ）

`main` への push 時に VPS へ自動デプロイします。

#### 初回セットアップ

**1. GitHub Secrets の登録**

リポジトリの `Settings > Secrets and variables > Actions` に以下を登録：

| Secret 名 | 内容 |
|---|---|
| `VPS_HOST` | VPS のホスト名または IP |
| `VPS_USER` | SSH ログインユーザー名 |
| `SSH_PRIVATE_KEY` | SSH 秘密鍵（`~/.ssh/id_ed25519` 等の中身） |

**2. VPS 側の sudoers 設定**

デプロイユーザーが `systemctl` をパスワードなしで実行できるよう設定：

```bash
sudo visudo -f /etc/sudoers.d/training-app
```

以下を追記：

```
YOUR_USER ALL=(ALL) NOPASSWD: \
  /bin/systemctl stop training-app, \
  /bin/systemctl start training-app, \
  /bin/mv /tmp/training-app /opt/training-scheduler/training-app, \
  /bin/chmod +x /opt/training-scheduler/training-app
```

**3. systemd サービスファイルの配置（未設置の場合）**

```ini
# /etc/systemd/system/training-app.service
[Unit]
Description=Training Scheduler App
After=network.target

[Service]
Type=simple
User=YOUR_USER
WorkingDirectory=/opt/training-scheduler
ExecStart=/opt/training-scheduler/training-app
Restart=always
MemoryMax=150M

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable training-app
```


