# POST /workspaces/:id/members

ワークスペースにメンバーを追加するエンドポイント。

## 概要

指定したワークスペースに新しいメンバーを招待・追加します。オーナー権限が必要です。

## リクエスト

### エンドポイント

```
POST /workspaces/:id/members
```

**パスパラメータ:**

| パラメータ | 型     | 説明               |
| ---------- | ------ | ------------------ |
| id         | number | ワークスペースID   |

### ヘッダー

```
Authorization: Bearer {token}
Content-Type: application/json
```

### リクエストボディ

| フィールド | 型     | 必須 | 説明                           |
| ---------- | ------ | ---- | ------------------------------ |
| user_id    | number | Yes  | 招待するユーザーのID           |
| role       | string | No   | 役割（デフォルト: "member"）   |

**例:**

```json
{
  "user_id": 5,
  "role": "member"
}
```

## レスポンス

### 成功時（201 Created）

```json
{
  "member": {
    "id": 5,
    "email": "newmember@example.com",
    "name": "新メンバー",
    "role": "member",
    "joined_at": "2025-11-16T12:00:00Z"
  }
}
```

**レスポンスフィールド:**

| フィールド      | 型     | 説明                                   |
| --------------- | ------ | -------------------------------------- |
| member.id       | number | ユーザーID                             |
| member.email    | string | メールアドレス                         |
| member.name     | string | ユーザー名                             |
| member.role     | string | 役割（"owner" / "member"）             |
| member.joined_at | string | メンバー追加日時（ISO 8601形式）       |

### エラーレスポンス

#### 400 Bad Request

リクエストパラメータが不正な場合

```json
{
  "error": "invalid_request",
  "message": "User ID is required"
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

#### 403 Forbidden

オーナー権限がない場合

```json
{
  "error": "forbidden",
  "message": "Only workspace owners can add members"
}
```

#### 404 Not Found

ワークスペースまたはユーザーが存在しない場合

```json
{
  "error": "not_found",
  "message": "Workspace or user not found"
}
```

#### 409 Conflict

ユーザーが既にワークスペースのメンバーである場合

```json
{
  "error": "already_member",
  "message": "User is already a member of this workspace"
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

1. リクエストユーザーがワークスペースのオーナーであることを確認
2. 追加するユーザーが存在することを確認（`users` テーブル）
3. 既に同じユーザーがメンバーでないことを確認（`workspace_members` テーブル）
4. ワークスペースメンバーシップの作成（`workspace_members` テーブル）
5. 追加されたメンバーの情報を返却

## 備考

- メンバーの追加はワークスペースオーナーのみが実行可能
- 同じユーザーを複数回追加することはできない（UNIQUE制約）
- role パラメータを省略した場合、デフォルトで "member" が設定される
- オーナー権限の付与も可能だが、慎重に行うこと
