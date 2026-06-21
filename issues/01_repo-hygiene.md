## repo 衛生: 空ファイル削除と .gitignore 掃除
id: 01
skill: pr-workflow
branch-slug: repo-hygiene
github_issue: 1
status: open
type: cleanup
対象: context/conventions.md, context/structure.md, .gitignore
内容: repo-guide.md の衛生基準に対する違反を解消する。0バイトの空ファイルを除去し、.gitignore から当スタック（Go + Vue）に無関係な boilerplate を削除する。
確認: 目視確認（git ls-files に成果物が含まれないこと、.gitignore に無関係行が残らないこと）
---

## 背景

`~/dotfiles/docs-agents/repo-guide.md`（リポジトリ衛生基準）に照らし、本 repo に2件の違反がある。

## 対象と作業

### context/conventions.md, context/structure.md
- いずれも 0 バイトの空ファイル。repo-guide「0バイト/プレースホルダだけのファイルを残さない」に違反。
- 対応: **中身を書く予定が無いなら削除**する。harness-guide 上 `context/` は規約・構造を書く場所だが、現状本 repo では未使用のため、空のまま残さず削除する方針。
- 削除後、`CLAUDE.md` 等から両ファイルへの `@import` 参照が無いことを確認する（あれば参照も除去）。

### .gitignore
- WordPress（`wp-content/...`）・Python（`venv/`・`__pycache__/` 等）・Docker（`.db_data/`）の無関係 boilerplate がテンプレ流用のまま残っている。
- 本 repo は Go + Vue + SQLite なので、無関係セクションを削除する。
- 残すべき対象: OS ファイル（`.DS_Store` 等）/ `node_modules/` / `dist/` / `.env` / `instance/`・`backend/instance/` / `backend/dist/` / `backend/training-app`。
- `backend/training-app` の重複行を1つに統合する。

## 確認
- `git ls-files` にバイナリ・DB・dist・node_modules が含まれないこと。
- `.gitignore` に無関係スタックの行・重複行が残らないこと。
- 削除した空ファイルへの参照が repo 内に残っていないこと。
