# POST /auth/signup

ユーザー登録を行うエンドポイント。

## 概要

新規ユーザーを作成し、自動的にデフォルトのワークスペースを作成します。

## リクエスト

### エンドポイント

```
POST /auth/signup
```

### ヘッダー

```
Content-Type: application/json
```

### リクエストボディ

| フィールド | 型     | 必須 | 説明                 |
| ---------- | ------ | ---- | -------------------- |
| email      | string | Yes  | メールアドレス       |
| password   | string | Yes  | パスワード（平文）   |
| name       | string | No   | ユーザー名（任意）   |

**例:**

```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "name": "山田太郎"
}
```

## レスポンス

### 成功時（201 Created）

```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "山田太郎",
    "created_at": "2025-11-16T10:00:00Z"
  },
  "workspace": {
    "id": 1,
    "name": "マイワークスペース",
    "role": "owner",
    "created_at": "2025-11-16T10:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**レスポンスフィールド:**

| フィールド         | 型     | 説明                                       |
| ------------------ | ------ | ------------------------------------------ |
| user.id            | number | ユーザーID                                 |
| user.email         | string | メールアドレス                             |
| user.name          | string | ユーザー名                                 |
| user.created_at    | string | ユーザー作成日時（ISO 8601形式）           |
| workspace.id       | number | 自動作成されたワークスペースID             |
| workspace.name     | string | ワークスペース名（デフォルト名）           |
| workspace.role     | string | ユーザーの役割（常に "owner"）             |
| workspace.created_at| string | ワークスペース作成日時（ISO 8601形式）     |
| token              | string | 認証トークン（JWT）                        |

### エラーレスポンス

#### 400 Bad Request

リクエストパラメータが不正な場合

```json
{
  "error": "invalid_request",
  "message": "Email is required"
}
```

```json
{
  "error": "invalid_request",
  "message": "Password must be at least 8 characters"
}
```

#### 409 Conflict

メールアドレスが既に登録されている場合

```json
{
  "error": "email_already_exists",
  "message": "This email is already registered"
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
2. メールアドレスの重複チェック
3. パスワードのハッシュ化（bcrypt）
4. ユーザーレコードの作成（`users` テーブル）
5. デフォルトワークスペースの作成（`workspaces` テーブル）
   - 名前: "マイワークスペース" などのデフォルト名
6. ワークスペースメンバーシップの作成（`workspace_members` テーブル）
   - role: "owner"
7. 認証トークンの生成
8. ユーザー情報、ワークスペース情報、トークンを返却

## 備考

- パスワードは平文で送信されるため、本番環境では必ずHTTPSを使用すること
- パスワードはサーバー側でbcryptを使用してハッシュ化される
- 作成されたデフォルトワークスペースは、ユーザー自身が唯一のオーナーとなる
- レスポンスに含まれるトークンは、以降のAPI呼び出しで認証に使用される
- デフォルトワークスペースの名前は後から変更可能
