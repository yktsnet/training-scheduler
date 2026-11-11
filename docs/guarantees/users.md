# users

## ユーザーの作成

| 保証 | 対応テスト |
|---|---|
| `POST /api/users` は emoji と initial からユーザーを作る | `TestCreateUser_Success` |
| `POST /api/users` は作成したユーザーを 200 で返す | `TestCreateUser_Success` |
| `POST /api/users` は initial を大文字にして保存する | `TestCreateUser_Success` |
| `POST /api/users` は作成したユーザーに 0 でない ID を付ける | `TestCreateUser_Success` |

## 作成時の入力の検証

| 保証 | 対応テスト |
|---|---|
| `POST /api/users` は同じ emoji が登録済みなら 400 を返す | `TestCreateUser_DuplicateEmoji` |
| `POST /api/users` は emoji が無いと 400 を返す | `TestCreateUser_EmojiRequired` |
| `POST /api/users` は initial が空文字か空白だけなら 400 を返す | `TestCreateUser_InitialValidation` |
| `POST /api/users` は initial が4文字以上なら 400 を返す | `TestCreateUser_InitialValidation` |
| `POST /api/users` は3文字の initial を受け付ける | `TestCreateUser_InitialValidation` |

## ユーザーの一覧と削除

| 保証 | 対応テスト |
|---|---|
| `GET /api/users` はユーザーが1人もいなければ 200 で空配列を返す | `TestGetUsers_Empty` |
| `DELETE /api/users/:id` はそのユーザーの Plan・Report・Progress も削除する | `TestDeleteUser_CascadesRelatedData` |
| `DELETE /api/users/:id` は存在しない ID でも 200 を返す | `TestDeleteUser_NonexistentID_NoOp` |
