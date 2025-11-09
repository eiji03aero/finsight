# POST /workspaces

新規ワークスペースを作成するエンドポイント。

## 概要

ログインユーザーがオーナーとなる新しいワークスペースを作成します。

## リクエスト

### エンドポイント

```
POST /workspaces
```

### ヘッダー

```
Authorization: Bearer {token}
Content-Type: application/json
```

### リクエストボディ

| フィールド | 型     | 必須 | 説明               |
| ---------- | ------ | ---- | ------------------ |
| name       | string | Yes  | ワークスペース名   |

**例:**

```json
{
  "name": "新しいプロジェクト"
}
```

## レスポンス

### 成功時（201 Created）

```json
{
  "workspace": {
    "id": 3,
    "name": "新しいプロジェクト",
    "role": "owner",
    "member_count": 1,
    "created_at": "2025-11-16T11:00:00Z"
  }
}
```

**レスポンスフィールド:**

| フィールド              | 型     | 説明                                   |
| ----------------------- | ------ | -------------------------------------- |
| workspace.id            | number | 作成されたワークスペースID             |
| workspace.name          | string | ワークスペース名                       |
| workspace.role          | string | ユーザーの役割（常に "owner"）         |
| workspace.member_count  | number | メンバー数（初期値は1）                |
| workspace.created_at    | string | ワークスペース作成日時（ISO 8601形式） |

### エラーレスポンス

#### 400 Bad Request

リクエストパラメータが不正な場合

```json
{
  "error": "invalid_request",
  "message": "Workspace name is required"
}
```

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

## 処理フロー

1. リクエストパラメータのバリデーション
2. ワークスペースレコードの作成（`workspaces` テーブル）
3. ワークスペースメンバーシップの作成（`workspace_members` テーブル）
   - user_id: リクエストユーザーのID
   - role: "owner"
4. 作成されたワークスペース情報を返却

## 備考

- ワークスペースを作成したユーザーは自動的にオーナーとなる
- 作成直後のメンバー数は常に1（作成者のみ）
- ワークスペース名に長さ制限がある場合は別途バリデーションを実施
