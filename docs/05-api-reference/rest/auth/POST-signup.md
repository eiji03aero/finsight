# POST /auth/signup

ユーザー登録を行うエンドポイント。

## 概要

新規ユーザーを作成し、自動的にユーザー指定の名前でワークスペースを作成します。
作成後は自動的にログイン状態となり、セッションクッキーが発行されます。

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

| フィールド     | 型     | 必須 | 説明                           |
| -------------- | ------ | ---- | ------------------------------ |
| email          | string | Yes  | メールアドレス                 |
| password       | string | Yes  | パスワード（平文、最小8文字）  |
| workspaceName  | string | Yes  | ワークスペース名               |

**例:**

```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "workspaceName": "My Workspace"
}
```

## レスポンス

### 成功時（201 Created）

**ヘッダー:**

```
Set-Cookie: finsight_session=abc123...; Path=/; HttpOnly; Secure; SameSite=Lax
```

**ボディ:**

```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "createdAt": "2025-12-27T10:00:00Z",
    "updatedAt": "2025-12-27T10:00:00Z"
  },
  "workspace": {
    "id": 1,
    "name": "My Workspace",
    "createdAt": "2025-12-27T10:00:00Z",
    "updatedAt": "2025-12-27T10:00:00Z"
  },
  "message": "Account created successfully"
}
```

**レスポンスフィールド:**

| フィールド          | 型     | 説明                                       |
| ------------------- | ------ | ------------------------------------------ |
| user.id             | number | ユーザーID                                 |
| user.email          | string | メールアドレス                             |
| user.createdAt      | string | ユーザー作成日時（ISO 8601形式）           |
| user.updatedAt      | string | ユーザー更新日時（ISO 8601形式）           |
| workspace.id        | number | 自動作成されたワークスペースID             |
| workspace.name      | string | ワークスペース名（ユーザー指定）           |
| workspace.createdAt | string | ワークスペース作成日時（ISO 8601形式）     |
| workspace.updatedAt | string | ワークスペース更新日時（ISO 8601形式）     |
| message             | string | 成功メッセージ                             |

**認証:**

- 認証はセッションクッキー (`finsight_session`) で管理されます
- クッキーは HttpOnly, Secure 属性が設定され、XSS攻撃から保護されます
- 以降のAPIリクエストでは、このクッキーが自動的に送信されます

### エラーレスポンス

#### 400 Bad Request

リクエストパラメータが不正な場合

**メールアドレス形式が不正:**

```json
{
  "error": "Invalid email format",
  "code": "VALIDATION_ERROR",
  "details": {
    "field": "email",
    "message": "Email must be a valid email address"
  }
}
```

**パスワードが短すぎる:**

```json
{
  "error": "Password must be at least 8 characters",
  "code": "VALIDATION_ERROR",
  "details": {
    "field": "password",
    "message": "Password must be at least 8 characters"
  }
}
```

**ワークスペース名が空:**

```json
{
  "error": "Workspace name is required",
  "code": "VALIDATION_ERROR",
  "details": {
    "field": "workspaceName",
    "message": "Workspace name cannot be empty"
  }
}
```

#### 409 Conflict

メールアドレスが既に登録されている場合

```json
{
  "error": "Email already registered",
  "code": "EMAIL_EXISTS",
  "details": {
    "field": "email",
    "message": "An account with this email already exists"
  }
}
```

#### 500 Internal Server Error

サーバー内部エラー

```json
{
  "error": "Internal server error",
  "code": "INTERNAL_ERROR",
  "details": {
    "message": "An unexpected error occurred. Please try again later."
  }
}
```

## 処理フロー

1. リクエストパラメータのバリデーション
   - メールアドレス形式チェック
   - パスワード長チェック（最小8文字）
   - ワークスペース名が空でないことをチェック
2. メールアドレスの重複チェック（大文字小文字を区別しない）
3. パスワードのハッシュ化（bcrypt）
4. トランザクション開始
5. ユーザーレコードの作成（`users` テーブル）
   - メールアドレスは小文字に正規化して保存
6. ワークスペースの作成（`workspaces` テーブル）
   - ユーザー指定の `workspaceName` を使用
7. ユーザー・ワークスペース関連付け（many-to-many リレーション）
8. トランザクションコミット
9. セッションの作成とクッキーの発行
10. ユーザー情報、ワークスペース情報、成功メッセージを返却

## バリデーション

### メールアドレス

- 形式: RFC 5322 準拠の標準的なメール形式
- 一意性: 大文字小文字を区別せず一意（test@example.com = Test@Example.com）
- データベース制約: UNIQUE 制約により重複を防止

### パスワード

- 最小長: 8文字
- ハッシュ化: bcrypt (cost factor 10以上)
- 保存: ハッシュ化された値のみをデータベースに保存

### ワークスペース名

- 必須項目
- 特殊文字・絵文字使用可能
- 一意性制約なし（複数のワークスペースが同じ名前を持つことが可能）

## セキュリティ

- **HTTPS必須**: パスワードは平文で送信されるため、本番環境では必ずHTTPSを使用すること
- **パスワードハッシュ化**: bcryptを使用してサーバー側でハッシュ化
- **セッション管理**: HttpOnly, Secure, SameSite=Lax 属性付きクッキーで管理
- **メール正規化**: メールアドレスは小文字に正規化して保存し、重複を防止
- **データベース制約**: メールアドレスの一意性はデータベースレベルでも強制

## 備考

- 作成されたワークスペースは、ユーザー自身が唯一のメンバーとなります
- ワークスペース名は後から変更可能です
- セッションの有効期限はデフォルトで7日間です
- メールアドレス確認（Email Verification）は現時点では実装されていません
