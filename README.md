# Training Scheduler

[![CI](https://github.com/yktsnet/training-scheduler/actions/workflows/ci.yml/badge.svg)](https://github.com/yktsnet/training-scheduler/actions/workflows/ci.yml)

新入社員の自律性を促すことを目的とした研修支援ツールです。「システムによる自動管理」と「手書き感覚のアナログ操作」を融合させ、ガチガチの進捗管理ではなく、新人の「主観的な手応え」をベースにメンターが静かに見守るためのアプリケーションです。

---

## Quick Start

### Prerequisites
- [Go 1.25+](https://go.dev/)
- [Node.js 20+](https://nodejs.org/)

### Setup
リポジトリをクローンしてビルドを実行し、Webサーバーを起動します。

```bash
# プロジェクトのビルド（フロントエンドのビルドとGoへの埋め込みを一括実行）
make build

# デモモード（30分ごとの自動リセット有効）かつ初期パスワードを指定して起動
DEMO_MODE=true ADMIN_PASSWORD=admin123 ./backend/training-app
```

- アプリ起動URL: http://localhost:5000
  - ※起動ポートは環境変数 `PORT` を指定することで変更可能です（例: `PORT=8080 ./backend/training-app`）。
- 管理者ログイン用の初期パスワード: `admin123`

---

## Overview

本ツールは、チーム内の信頼関係を前提とした小規模チーム向けの研修プランナーです。一般的なガントチャート型の厳格な進捗管理ツールとは異なり、新人の主体的な内省と、メンターのゆるやかな見守りをサポートすることに特化しています。

- **主体的プランニング**: システムは枠組みだけを提示し、具体的な計画は新人が自身の言葉で記述します。
- **内省の可視化**: 機械的な進捗率（％）の計算ではなく、本人の「主観的なズレ（手応え）」をマネージャと共有します。
- **非干渉の監視**: マネージャは新人の自律を妨げず、ダッシュボードから状況を静かに見守り、必要な時だけサポートに入ります。
- **ゆるやかな識別（アニマル・ログイン）**: パスワード等による厳格な認証ではなく、動物の絵文字を選ぶだけのシンプルなログインを採用。チーム内の信頼関係を前提とした、遊び心のあるアカウント管理です。

### Demo Mode & Admin Panel
本番やデモ用にURLを一般公開する際の安全性を担保するため、以下の機能が備わっています。

- **管理者機能 (Admin Panel)**:
  - アニマル選択画面の最下部「🐾 管理者としてログイン」から、管理者用パスワードでログインできます。
  - カリキュラム（メニュー）の追加・編集・削除をブラウザ上で動的に行えます。変更はDBと `menu_config.json` に即座に同期されます。
- **デモモード (Demo Mode)**:
  - 環境変数 `DEMO_MODE=true` で起動すると有効化されます。
  - 30分ごとにDBの変更差分を検知し、データが書き換えられている場合に初期ダミーデータ（🐶 ユーザー、Git計画、2日分の日報、進捗）へ自動リセットします。
- **デモモードの停止方法**:
  - 起動時に環境変数 `DEMO_MODE` を指定しない、または `DEMO_MODE=false` を設定することで、30分自動リセット機能を無効化できます（通常の本番運用では無効にしてください）。
- **パスワードの変更方法**:
  - 起動時に環境変数 `ADMIN_PASSWORD=新しいパスワード` を指定します。未指定の場合はデフォルト値の `admin123` が使用されます。

---

## User Interface

### User (Animal Login)

<img src="src/animals.png" width="500" alt="menu-pic">

- **役割**: アプリを利用する個人（新人・メンター）の識別。
- **項目**: `emoji` (🦁や🐰などのユニークな絵文字)。

### Menu (Curriculum)

<img src="src/menu.png" width="500" alt="menu-pic">

- **役割**: 研修カリキュラムのマスターデータ（全ユーザー共通）。
- **項目**: 名称、目安日数、概要、参考URL。
- ※ `internal/database/menu_config.json` をマスターとして起動時に自動同期します。

### Plan (Training Plan)

<img src="src/plan.png" width="500" alt="plan-pic">

- **役割**: 各メニューに対する具体的な学習計画。
- **項目**: `content` (自由記述のテキスト)、`user_id`。

### Report (Daily Log)

<img src="src/daily.png" width="500" alt="daily-pic">

- **役割**: 日付単位の事実と内省の記録。
- **項目**: `date` (YYYY-MM-DD)、`content` (日報内容)、`user_id`。

### Progress (Status & Condition)

<img src="src/overview.png" width="500" alt="overview-pic">

- **役割**: ダッシュボード表示用のメタ情報。
- **項目**: 開始日、目標日数、`offset_days` (主観ズレ値 1〜5)、ステータスメモ。

---

## Architecture

```mermaid
graph TD
    %% 1. 関係者を配置
    Admin["管理者"]
    Newcomer["新人 (アニマルログイン)"]
    Mentor["メンター (見守り手)"]

    %% 2. 下部のデータライフサイクル (左から右への直感的な流れ)
    Menu["① 研修メニュー <br>(共通カリキュラム)"] -->|選択して計画化| Plan["② 個人の計画 <br>(自由記述の目標)"]
    Plan -->|日々の実行と振り返り| Report["③ 日報 <br>(事実と内省の記録)"]
    Report -->|主観による自己評価| Progress["④ 手応え・ズレ <br>(1〜5 の自己評価)"]
    Progress -->|進捗の自動集約| Dashboard["全体ダッシュボード <br>(Overview)"]

    %% 3. 関係者からデータフローへのアプローチ (上から下への矢印)
    Admin -->|カリキュラムの登録/編集| Menu
    
    Newcomer -->|目標を記入| Plan
    Newcomer -->|日々の出来事を記入| Report
    Newcomer -->|手応えを評価| Progress
    
    Mentor -->|ダッシュボードで状況を俯瞰| Dashboard
```

---

## Tech Stack

| Layer | Technology | Reason |
|---|---|---|
| **Frontend** | Vue 3, Vite, Vue Router | リアクティブなUI構築と、シングルページアプリケーション（SPA）のルーティングをシンプルに統合するため。 |
| **Backend** | Go (Gin), GORM | 高パフォーマンスかつ静的なGoの型安全性を活かし、Web APIを軽量かつ高速に提供するため。 |
| **Database** | SQLite (Pure Go driver) | 外部データベースサーバーの設定や運用管理コストをゼロにし、単一ファイルのみで動作を完結させるため。 |
| **Embedding** | go:embed | フロントエンドのビルド資産（HTML/JS/CSS）をGoのバイナリ自体に埋め込み、単一バイナリだけで配布・起動できるようにするため。 |

---

## Design Decisions

- **アニマルログイン（ゆるやかな識別）**: 
  パスワードによる厳格な認証をあえて排し、動物の絵文字を選ぶだけのカジュアルなアカウント管理を採用しています。これはチーム内の信頼関係を前提に、誰がどの作業をしているかを気楽に共有するための設計です。
- **SQLite + JSONのハイブリッド同期**:
  管理画面からの編集はSQLiteへ直接書き込まれますが、開発環境やデプロイ時の整合性を担保するため、同時に `menu_config.json` にも書き出されます。これにより、メニュー設定がGit管理可能になります。
- **デモモードでの動的日付 Seeder**:
  デモ起動時に日付を `現在日時 - 3日` などと相対計算し、日報や進捗が常に「直近数日間」のものとしてリアルに再現されるように設計しています。

---

## Scope

### In Scope
- 各自のアニマル（絵文字）による簡易ログイン
- 研修カリキュラムに応じた学習計画（Plan）の作成と自己編集
- 1日単位のシンプルな日報（Daily Log）入力
- 新人の主観的なズレ（手応え 1〜5）とメモを共有するダッシュボード（Overview）
- 管理者画面からの研修メニューのCRUD操作およびJSONファイルの自動保存

### Out of Scope
- パスワードを用いた一般ユーザー認証（アニマルログインのみ）
- メンターや管理者による一般ユーザーデータの直接編集（読み取りのみ可）

---

## Deploy

`main` ブランチへのプッシュにより、GitHub Actionsでテストおよび自動ビルドが行われ、本番サーバーへデプロイされます。フロントエンドの静的ファイルがGoバイナリに埋め込まれているため、生成された単一の実行ファイルをサーバーに配置して起動するだけでデプロイが完了します。

---

## Development

### Local Run
開発用にフロントエンドとバックエンドをそれぞれホットリロード有効で起動します。

```bash
# 初回のみ: go:embed 用のスタブ作成
make dev-dist

# ターミナル1: バックエンド (port 5000)
make dev-back

# ターミナル2: フロントエンド (port 5173, HMR有効)
make dev-front
```

### Running Tests

```bash
make test
```

外部依存なし（in-memory SQLite）。Go の標準ライブラリのみで動作します。

---

## 🚀 CI/CD

### CI (Continuous Integration)

`main` / `go-dev` への push および pull request 時に GitHub Actions で自動実行されます。

1. Go テスト (`go test ./internal/...`)
2. フロントエンドビルド + Go バイナリビルドの疎通確認

### CD (Continuous Deployment)

`main` ブランチへのプッシュ時に、Tailscale VPN 経由で対象サーバーへ自動デプロイします。

#### デプロイ設計と技術選定の判断基準 (Tailscale経由プッシュ型自動デプロイ)

本プロジェクトのデプロイは、クラウド上の CI サービス（GitHub Actions）から実環境へ直接リリースを行う **「プッシュ（Push）型自動デプロイ方式」** を採用しています。

* **構成の前提**: テスト・ビルド・デプロイの一連を GitHub Actions 上で完結させ、ビルドログと成否を GitHub に一元管理します。
* **VPN（Tailscale）を用いたセキュアな一時接続**:
  オンプレミス環境やプライベートネットワーク内のサーバーに対してデプロイを行う際、インターネット上に SSH ポート（22）を直接開放することはセキュリティ上の大きな脅威となります。本システムでは、デプロイ実行時のみ GitHub Actions ランナーから一時的に Tailscale（閉域網 VPN）へ接続させ、認証されたセキュアな経路を通じてのみ SSH デプロイを行うことで、高い開発効率とインフラセキュリティを両立しています。
* **物理サーバーのビルド負荷軽減**:
  Go のビルドや Node.js のビルド処理はすべて GitHub Actions（クラウド環境）側で行われ、最適化された単一バイナリのみがサーバーに転送されます。デプロイ時のオンプレミスサーバーの CPU やメモリ負荷をほぼゼロに抑え、本番稼働中のサービスへの影響を最小限に抑えています。

#### Initial Setup

**1. GitHub Secrets の登録**

リポジトリの `Settings > Secrets and variables > Actions` に以下を登録：

| Secret 名 | 内容 |
|---|---|
| `DEPLOY_HOST` | デプロイ先サーバーのホスト名または IP |
| `DEPLOY_USER` | デプロイ用 SSH ログインユーザー名 |
| `SSH_PRIVATE_KEY` | SSH 秘密鍵（`~/.ssh/id_ed25519` 等の中身） |
| `TS_OAUTH_CLIENT_ID` | Tailscale OAuth Client ID |
| `TS_OAUTH_SECRET` | Tailscale OAuth Client Secret |

**2. デプロイ先サーバー側の sudoers 設定**

デプロイユーザーがサービス再起動やバイナリの配置をパスワードなしで実行できるよう設定します。

```bash
sudo visudo -f /etc/sudoers.d/training-scheduler
```

以下を追記します（ユーザー名や配置パスは環境に合わせて調整してください）：

```
YOUR_USER ALL=(ALL) NOPASSWD: \
  /usr/bin/systemctl restart training-scheduler, \
  /usr/bin/mv /tmp/training-app /opt/training-scheduler/training-app, \
  /usr/bin/chmod +x /opt/training-scheduler/training-app
```

**3. systemd サービスファイルの配置（未設置の場合）**

```ini
# /etc/systemd/system/training-scheduler.service
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
sudo systemctl enable training-scheduler
```


