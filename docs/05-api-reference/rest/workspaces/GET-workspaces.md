# GET /workspaces

ログインユーザーが所属するワークスペース一覧を取得するエンドポイント。

## 概要

現在ログインしているユーザーが所属する全てのワークスペースを取得します。

## リクエスト

### エンドポイント

```
GET /workspaces
```

### ヘッダー

```
Authorization: Bearer {token}
```

### クエリパラメータ

なし

## レスポンス

### 成功時（200 OK）

```json
{
  "workspaces": [
    {
      "id": 1,
      "name": "マイワークスペース",
      "role": "owner",
      "member_count": 1,
      "created_at": "2025-11-16T10:00:00Z"
    },
    {
      "id": 2,
      "name": "家族の家計簿",
      "role": "member",
      "member_count": 3,
      "created_at": "2025-11-15T14:30:00Z"
    }
  ]
}
```

**レスポンスフィールド:**

| フィールド            | 型     | 説明                                   |
| --------------------- | ------ | -------------------------------------- |
| workspaces            | array  | ワークスペース一覧                     |
| workspaces[].id       | number | ワークスペースID                       |
| workspaces[].name     | string | ワークスペース名                       |
| workspaces[].role     | string | ユーザーの役割（"owner" / "member"）   |
| workspaces[].member_count | number | ワークスペースのメンバー数         |
| workspaces[].created_at | string | ワークスペース作成日時（ISO 8601形式） |

### エラーレスポンス

#### 401 Unauthorized

認証トークンが無効または欠落している場合

```json
{
  "error": "unauthorized",
  "message": "Authentication required"
}
```

#### 500 Internal Server Error

サーバー内部エラー

```json
{
  "error": "internal_error",
  "message": "An unexpected error occurred"
}
```

## 備考

- 返却されるワークスペースは、ユーザーが `workspace_members` テーブルに登録されているものに限られる
- `role` は "owner" または "member" のいずれか
- ワークスペースは作成日時の降順で返却される
