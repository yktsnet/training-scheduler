## Admin メニュー作成・更新の days バリデーション
id: 03
skill: pr-workflow
branch-slug: admin-menu-days-validation
github_issue:
status: open
type: fix
対象: backend/internal/handlers/admin_handlers.go
内容: `AdminHandler.CreateMenu` / `UpdateMenu` は `models.Menu` に直接 `ShouldBindJSON` しており、`days` に対する `binding` タグが無いため 0 以下の値もそのまま DB に保存できてしまう。保証台帳の棚卸し（guarantee-audit）で発覚した欠落バリデーション。
確認: `make test` が通ること。既存の `TestAdminCreateMenu_Success` / `TestAdminUpdateMenu_Success` を壊さないこと。
---

## 背景

`docs/guarantees.md` の棚卸しで、`days` の妥当性検証テストを追加しようとしたところ、現状の実装は days が 0 や負数でもバリデーションなしで受理してしまうことが判明した（テストを書くと red になるため、テスト追加ではなく本Issueとして起票）。

## 設計方針（叩き台）

- `CreateMenu` / `UpdateMenu` それぞれで `days <= 0` の場合に 400 を返すガードを追加する
- リクエスト用の別 struct にバインドして `binding:"required,min=1"` を使うか、bind 後に手動チェックするかは実装者判断
- 実装後、`docs/guarantees.md` の `AdminHandler.CreateMenu / UpdateMenu は不正な入力（days が 0 以下など）に対するバリデーションが未保証` という Gaps 行を本体の保証に昇格させ、対応するテストを追加する
