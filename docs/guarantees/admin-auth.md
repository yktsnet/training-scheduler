# admin-auth

| 保証 | 対応テスト |
|---|---|
| `POST /api/admin/login` は `ADMIN_PASSWORD` と一致するパスワードに 200 を返す | `TestAdminLogin_Success` |
| `POST /api/admin/login` は `ADMIN_PASSWORD` と一致するパスワードに token を返す | `TestAdminLogin_Success` |
| `POST /api/admin/login` は誤ったパスワードに 401 を返す | `TestAdminLogin_Failure` |
| 管理者 API は `X-Admin-Token` ヘッダーが無いと 401 を返す | `TestAdminAuthMiddleware` |
| 管理者 API は正しい `X-Admin-Token` があればリクエストを通す | `TestAdminAuthMiddleware` |
| 管理者 API は誤った `X-Admin-Token` に 401 を返す | `TestAdminAuthMiddleware_WrongToken` |
