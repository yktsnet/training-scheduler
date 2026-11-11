# admin-menus

## メニューの一覧と作成

| 保証 | 対応テスト |
|---|---|
| `GET /api/admin/menus` は登録済みの全メニューを返す | `TestAdminGetMenus_ReturnsAll` |
| `POST /api/admin/menus` はメニューを保存する | `TestAdminCreateMenu_Success` |
| `POST /api/admin/menus` は ID の付いたメニューを 200 で返す | `TestAdminCreateMenu_Success` |
| `POST /api/admin/menus` は days が 0 以下なら 400 を返す | `TestAdminCreateMenu_InvalidDays` |
| `POST /api/admin/menus` は days が 0 以下のメニューを保存しない | `TestAdminCreateMenu_InvalidDays` |
| `POST /api/admin/menus` は保存後の全メニューを `menu_config.json` へ書き出す | `TestAdminCreateMenu_WritesMenuConfigJSON` |

## メニューの更新

| 保証 | 対応テスト |
|---|---|
| `PUT /api/admin/menus/:id` はメニューの name・days を書き換える | `TestAdminUpdateMenu_Success` |
| `PUT /api/admin/menus/:id` は days が 0 以下なら 400 を返す | `TestAdminUpdateMenu_InvalidDays` |
| `PUT /api/admin/menus/:id` は days が 0 以下の送信でメニューを変えない | `TestAdminUpdateMenu_InvalidDays` |
| `PUT /api/admin/menus/:id` は存在しない ID に 404 を返す | `TestAdminUpdateMenu_NonexistentID_Returns404` |

## メニューの削除

| 保証 | 対応テスト |
|---|---|
| `DELETE /api/admin/menus/:id` はメニューを削除する | `TestAdminDeleteMenu_Success` |
| `DELETE /api/admin/menus/:id` はそのメニューの Plan・Progress も削除する | `TestAdminDeleteMenu_Success` |
| `DELETE /api/admin/menus/:id` は存在しない ID に 404 を返す | `TestAdminDeleteMenu_NonexistentID_Returns404` |
