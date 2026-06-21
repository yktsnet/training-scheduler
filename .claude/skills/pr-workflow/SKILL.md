---
name: pr-workflow
description: Issue に基づきブランチ作成 → 実装 → 検証 → PR 作成までを行う。issue() から起動される。
---

# pr-workflow

`issues/` の Issue 1件を実装し、PR を出すまでのワークフロー。

## 手順

1. **ブランチ作成**: `claude/{id}-{branch-slug}` を `main` から切る。
2. **実装**: Issue の `対象` / `内容` / 詳細セクションに従う。Issue ファイルの `status:` は変更しない（issue-finish が処理する）。
3. **検証**（提出前に必ず実行）:
   - backend: `make test`（in-memory SQLite、外部依存なし）
   - frontend: `cd frontend && npm ci && npm run build`（ビルドが通ること＝型/参照の静的確認）
   - Issue の `確認` 欄に書かれた静的チェックを併せて実施する。
4. **PR 作成**: `gh pr create`。本文の `## 検証手順` には Agent 側で完結しない確認（ブラウザ目視・本番動作）のみを書き、user に委ねる。

## 制約

- `main` への直接 push・force push はしない（settings.json の deny に従う）。
- 人間が読む説明文（コミット・PR・コメント）に固有の接続情報を直書きしない（`~/dotfiles/secrets-agents/` の辞書に従う）。
- コミットは conventional commits（`context/conventions.md`）。
