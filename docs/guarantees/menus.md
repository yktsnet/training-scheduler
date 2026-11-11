# menus

| 保証 | 対応テスト |
|---|---|
| `GET /api/menus` は登録済みの全メニューを返す | `TestGetMenus_ReturnsAll` |
| `POST /api/menus/select` は選んだメニューごとに、そのユーザーのプランを作る | `TestSaveSelection_CreatesNewPlans` |
| `POST /api/menus/select` が作るプランの content は「【メニュー名 研修計画（計N日間）】」の行で始まる | `TestSaveSelection_CreatesNewPlans` |
| `POST /api/menus/select` が作るプランの content は「1日目：」から「N日目：」までの行を持つ | `TestSaveSelection_CreatesNewPlans` |
| `POST /api/menus/select` は既存のプランについて、先頭のヘッダー行を現在のメニュー名・日数で書き換える | `TestSaveSelection_UpdatesHeaderOnly` |
| `POST /api/menus/select` は既存のプランの2行目以降を変えない | `TestSaveSelection_UpdatesHeaderOnly` |
| `POST /api/menus/select` は選ばれなかったメニューのプランを削除する | `TestSaveSelection_DeletesDeselectedPlans` |
| `POST /api/menus/select` は存在しない menu_id を無視して 200 を返す | `TestSaveSelection_NonexistentMenuID_NoOp` |
| `POST /api/menus/select` は `X-User-Id` ヘッダーが無いと 401 を返す | `TestSaveSelection_NoAuth_Returns401` |
