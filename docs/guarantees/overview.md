# overview

| 保証 | 対応テスト |
|---|---|
| `GET /api/overview` は Progress の無いメニューの target_days をメニューの日数にする | `TestGetOverview_FallsBackToMenuDays` |
| `GET /api/overview` は Progress のあるメニューの target_days を Progress の値にする | `TestGetOverview_UsesProgressTargetDaysWhenPresent` |
| `POST /api/overview` は同じユーザー・同じメニューの2回目の送信で、行を増やさず offset_days を書き換える | `TestUpdateOverview_UpsertsProgress` |
| `POST /api/overview` は `X-User-Id` ヘッダーが無いと 401 を返す | `TestUpdateOverview_NoAuth_Returns401` |
| `GET /api/overview` は `X-User-Id` ヘッダーが無いと 200 で空配列を返す | `TestGetOverview_NoAuth_ReturnsEmpty` |
