# reports

| 保証 | 対応テスト |
|---|---|
| `POST /api/reports` は送った date・content で日報を保存する | `TestSaveReport_Creates` |
| `POST /api/reports` は同じユーザー・同じ date の2回目の送信で、行を増やさず content を書き換える | `TestSaveReport_Upserts` |
| `POST /api/reports` は `X-User-Id` ヘッダーが無いと 401 を返す | `TestSaveReport_NoAuth_Returns401` |
| `GET /api/reports` は日報を date の新しい順に返す | `TestGetReports_OrderedByDateDesc` |
| `GET /api/reports` は `X-User-Id` ヘッダーが無いと 200 で空配列を返す | `TestGetReports_NoAuth_ReturnsEmpty` |
| `GET /api/reports` は他のユーザーの日報を含めない | `TestGetReports_ExcludesOtherUsers` |
