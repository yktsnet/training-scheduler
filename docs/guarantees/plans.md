# plans

| 保証 | 対応テスト |
|---|---|
| `GET /api/plans` はプランを menu_name 付きで返す | `TestGetPlans_ReturnsWithMenuName` |
| `GET /api/plans` は `X-User-Id` ヘッダーが無いと 200 で空配列を返す | `TestGetPlans_NoAuth_ReturnsEmpty` |
| `POST /api/plans/:id` は所有者以外が送ると 404 を返す | `TestUpdatePlan_OnlyOwnerCanUpdate` |
| `POST /api/plans/:id` は所有者以外が送っても content を変えない | `TestUpdatePlan_OnlyOwnerCanUpdate` |
| `POST /api/plans/:id` は存在しない ID に 404 を返す | `TestUpdatePlan_NonexistentID_Returns404` |
| `POST /api/plans/:id` は所有者が送ると 200 を返す | `TestUpdatePlan_OwnerUpdatesContent` |
| `POST /api/plans/:id` は所有者が送った content でプランを置き換える | `TestUpdatePlan_OwnerUpdatesContent` |
| `GET /api/plans` は他のユーザーのプランを含めない | `TestGetPlans_ExcludesOtherUsers` |
