# Guarantee Ledger

## Guarantees

### 1. `backend/internal/handlers/users_test.go` — UsersHandler (GetUsers / CreateUser / DeleteUser)

- `GetUsers` はユーザーが1人もいなければ空配列を返す
- `CreateUser` は emoji と initial からユーザーを作成し、initial は大文字化して保存する
- `CreateUser` は作成したユーザーに非ゼロの ID を採番する
- `CreateUser` は同じ emoji が既に登録されていると 400 を返す
- `CreateUser` は initial が空文字・空白のみ・4文字以上のいずれかだと 400 を返し、3文字以内なら 200 を返す
- `DeleteUser` はユーザーを削除すると、そのユーザーに紐づく Plan・Report・Progress も連鎖して削除する
- `CreateUser` は emoji が無い（未指定・空文字）と 400 を返す
- `DeleteUser` は存在しない ID を指定しても 200 を返す（no-op で成功扱いになる）

| 保証（要約） | 対応テスト |
|---|---|
| 空一覧 | `TestGetUsers_Empty` |
| 作成・initial大文字化・ID採番 | `TestCreateUser_Success` |
| emoji重複拒否 | `TestCreateUser_DuplicateEmoji` |
| initialバリデーション | `TestCreateUser_InitialValidation` |
| 削除の連鎖 | `TestDeleteUser_CascadesRelatedData` |
| emoji必須 | `TestCreateUser_EmojiRequired` |
| 存在しないIDの削除はno-op | `TestDeleteUser_NonexistentID_NoOp` |

*（区切り内は `UsersHandler` の3メソッドが混在するので、行ごとに主語を明示する）*

### 2. `backend/internal/handlers/reports_test.go` — ReportsHandler (SaveReport / GetReports)

- `SaveReport` は日報を作成し、指定した date・content で保存する
- `SaveReport` は同一 user + date の組に対して upsert する（2回目の送信で新規行を作らず、既存行の content を更新する）
- `SaveReport` は `X-User-Id` ヘッダーが無いと 401 を返す
- `GetReports` は認証ヘッダーが無いと 200 で空配列を返す
- `GetReports` は date の降順で返す

| 保証（要約） | 対応テスト |
|---|---|
| 作成 | `TestSaveReport_Creates` |
| upsert | `TestSaveReport_Upserts` |
| 未認証は401 | `TestSaveReport_NoAuth_Returns401` |
| 未認証は空配列 | `TestGetReports_NoAuth_ReturnsEmpty` |
| 日付降順 | `TestGetReports_OrderedByDateDesc` |

*（区切り内に `SaveReport` と `GetReports` の2メソッドが混在するので、行ごとに主語を明示する）*

### 3. `backend/internal/handlers/overview_test.go` — OverviewHandler (GetOverview / UpdateOverview)

- `GetOverview` は対象メニューに Progress が無い場合、target_days を Menu.Days にフォールバックする
- `UpdateOverview` は同一 user + menu の組に対して upsert する（2回目の送信で新規行を作らず、既存行の offset_days を更新する）
- `UpdateOverview` は `X-User-Id` ヘッダーが無いと 401 を返す
- `GetOverview` は対象メニューに Progress が既に存在する場合、target_days は Menu.Days ではなく Progress 側の値を使う

| 保証（要約） | 対応テスト |
|---|---|
| target_daysのフォールバック | `TestGetOverview_FallsBackToMenuDays` |
| upsert | `TestUpdateOverview_UpsertsProgress` |
| 未認証は401 | `TestUpdateOverview_NoAuth_Returns401` |
| Progress優先のtarget_days | `TestGetOverview_UsesProgressTargetDaysWhenPresent` |

*（区切り内に `GetOverview` と `UpdateOverview` の2メソッドが混在するので、行ごとに主語を明示する）*

### 4. `backend/internal/handlers/menus_test.go` — MenusHandler (GetMenus / SaveSelection)

- `GetMenus` は登録済みの全メニューを返す
- `SaveSelection` は選択したメニューごとに新規プランを作成し、`content` の先頭行を「【メニュー名 研修計画（計N日間）】」ヘッダーに、続けて「1日目：」〜「N日目：」の日付テンプレートにする
- `SaveSelection` は既存プランに対しては先頭のヘッダー行だけを更新し、ユーザーが書き込んだ本文（2行目以降）は保持する
- `SaveSelection` は選択解除されたメニューに対応する既存プランを削除する
- `SaveSelection` は `X-User-Id` ヘッダーが無いと 401 を返す
- `SaveSelection` は存在しない menu_id を指定してもエラーにせず無視する（プランは作られず 200 を返す）

| 保証（要約） | 対応テスト |
|---|---|
| 全件取得 | `TestGetMenus_ReturnsAll` |
| 新規プラン作成 | `TestSaveSelection_CreatesNewPlans` |
| ヘッダーのみ更新 | `TestSaveSelection_UpdatesHeaderOnly` |
| 選択解除で削除 | `TestSaveSelection_DeletesDeselectedPlans` |
| 未認証は401 | `TestSaveSelection_NoAuth_Returns401` |
| 存在しないmenu_idはno-op | `TestSaveSelection_NonexistentMenuID_NoOp` |

*（区切り内に `GetMenus` と `SaveSelection` の2メソッドが混在するので、行ごとに主語を明示する）*

### 5. `backend/internal/handlers/plans_test.go` — PlansHandler (GetPlans / UpdatePlan)

- `GetPlans` はプラン一覧をメニュー名付きで返す
- `GetPlans` は認証ヘッダーが無いと 200 で空配列を返す
- `UpdatePlan` は所有者以外のユーザーが更新しようとすると 404 を返し、content は変更しない
- `UpdatePlan` は存在しない plan id を指定しても、所有者不一致と同じ 404 を返す

| 保証（要約） | 対応テスト |
|---|---|
| メニュー名付き取得 | `TestGetPlans_ReturnsWithMenuName` |
| 未認証は空配列 | `TestGetPlans_NoAuth_ReturnsEmpty` |
| 所有者以外は404かつ不変更 | `TestUpdatePlan_OnlyOwnerCanUpdate` |
| 存在しないIDも404 | `TestUpdatePlan_NonexistentID_Returns404` |

*（区切り内に `GetPlans` と `UpdatePlan` の2メソッドが混在するので、行ごとに主語を明示する）*

### 6. `backend/internal/handlers/admin_handlers_test.go` — AdminHandler (Login / CreateMenu / UpdateMenu / DeleteMenu) / AdminAuthMiddleware

- `AdminHandler.Login` は正しいパスワードで 200 とトークンを返す
- `AdminHandler.Login` は誤ったパスワードで 401 を返す
- `AdminAuthMiddleware` は `X-Admin-Token` ヘッダーが無いと 401 を返し、正しいトークンなら後続処理に到達させる
- `AdminHandler.CreateMenu` はメニューを DB に保存する
- `AdminHandler.UpdateMenu` はメニューの name・days を更新する
- `AdminHandler.DeleteMenu` はメニューを削除すると、そのメニューに紐づく Plan・Progress も連鎖して削除する
- `AdminHandler.Login` は `ADMIN_PASSWORD` 環境変数が未設定の場合、デフォルトパスワード `admin123` でログインできる

| 保証（要約） | 対応テスト |
|---|---|
| ログイン成功 | `TestAdminLogin_Success` |
| ログイン失敗 | `TestAdminLogin_Failure` |
| 認証ミドルウェア | `TestAdminAuthMiddleware` |
| メニュー作成 | `TestAdminCreateMenu_Success` |
| メニュー更新 | `TestAdminUpdateMenu_Success` |
| メニュー削除の連鎖 | `TestAdminDeleteMenu_Success` |
| 未設定時のデフォルトパスワード | `TestAdminLogin_DefaultPasswordWhenEnvUnset` |

*（区切り内に複数のハンドラメソッドとミドルウェアが混在するので、行ごとに主語を明示する）*

## Gaps

以下は保証すべきと思われるが、対応するテストが無い。

- `AdminHandler.CreateMenu` / `UpdateMenu` は不正な入力（days が 0 以下など）に対するバリデーションが実装されていない（テストを追加すると red になるため見送り。実装修正を `issues/03_admin-menu-days-validation.md` として起票済み）

## About

対象は `backend/internal/handlers/` 配下の各ハンドラが公開する HTTP エンドポイントとミドルウェアの、リクエスト〜レスポンス・DB 状態変化として外部から観測可能な振る舞い。対象外は frontend（未着手）と、admin トークン生成方法などハンドラ内部の実装詳細。**ここに載っていない振る舞いは約束ではなく、予告なく変わりうる。** 地位は design-decisions.md 相当のドキュメントと同格。
