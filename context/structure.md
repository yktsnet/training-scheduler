# Structure

どこに何があるか。コードの書き方（規約）は `conventions.md` を参照。

## トップレベル

```
training-scheduler/
├── backend/          # Go 本体（単一バイナリに frontend を embed）
├── frontend/         # Vue 3 + Vite
├── Makefile          # build / test / dev タスク
├── context/          # Agent 向け共通コンテキスト（本ファイル群）
└── issues/           # ローカル Issue 管理（done/ に完了分と PR 控え）
```

## backend/（Go）

```
backend/
├── main.go                 # エントリポイント。frontend のビルド資産を go:embed で同梱
└── internal/
    ├── handlers/           # HTTP ハンドラ（リソース単位で1ファイル）
    │   ├── users.go / plans.go / reports.go / menus.go
    │   ├── overview.go / admin_handlers.go
    ├── models/             # GORM 構造体
    └── database/           # 接続・seeder
        └── menu_config.json # 研修メニューの原本（起動時に DB へ同期）
```

- DB は SQLite（Pure Go ドライバ）。DB ファイルは生成物として ignore する。
- `backend/dist/` は frontend のビルド資産の置き場（embed 対象、生成物）。

## frontend/（Vue 3）

```
frontend/
└── src/
    ├── views/{Name}View.vue   # 画面
    ├── components/            # 再利用部品（PascalCase）
    ├── router/index.js        # ルーティング集約
    └── assets/
```

## データフロー

```
Vue(View) → fetch → Go handlers → GORM → SQLite
                                     ↑
menu_config.json（原本）→ 起動時に DB へ同期 / 管理画面の編集は JSON へ書き戻し
```

## レイヤー構成

- **配布層**: 単一 Go バイナリ。frontend のビルド資産を `go:embed`（`backend/dist`）で同梱し、外部 DB 不要・インフラコストゼロで配布。
- **表示層**: Vue 3 + Vite。API レスポンスの snake_case をそのまま受ける。
- **API 層**: `internal/handlers`（リソース単位）。
- **永続層**: SQLite（GORM）。原本は `menu_config.json`。

## issues/

- `{NN}_{slug}.md`: 実装対象 Issue。`status: open` のものを Agent が処理。
- `00_template.md`: Issue ひな形。
- `done/`: 完了 Issue と PR 控え（`{id}_{slug}_pr.md`）。
