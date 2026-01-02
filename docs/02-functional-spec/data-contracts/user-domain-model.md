---
id: DC-001
title: User Domain Model
status: implemented
related:
  features: [F-0001]
  tables: [users, workspaces]
updated: 2026-01-02
---

# User Domain Model

## 概要

このドキュメントは、ユーザー認証とワークスペース管理に関連するドメインモデルを定義します。User（ユーザー）エンティティとWorkspace（ワークスペース）エンティティ、およびそれらの関係性を含みます。

## エンティティ定義

### 1. User（ユーザー）

**目的**: システムにアカウントを持つユーザーを表現します。

**属性:**

| 属性名         | 型        | 制約                | 説明                                  |
| -------------- | --------- | ------------------- | ------------------------------------- |
| `id`           | `int`     | Primary Key, Auto   | ユーザーの一意識別子                  |
| `email`        | `string`  | Unique, Not Empty   | ログインに使用するメールアドレス      |
| `password_hash`| `string`  | Not Empty, Sensitive| bcrypt でハッシュ化されたパスワード   |
| `created_at`   | `time`    | Auto-set on create  | アカウント作成日時                    |
| `updated_at`   | `time`    | Auto-update         | 最終更新日時                          |

**ビジネスルール:**

- メールアドレスは**大文字小文字を区別せず一意**である必要があります
  - 例: `test@example.com` = `Test@Example.com` = `TEST@EXAMPLE.COM`
- メールアドレスはすべて**小文字に正規化**して保存されます
- パスワードは**決して平文で保存されません**（bcrypt ハッシュのみ）
- パスワードは最低8文字である必要があります（アプリケーション層で検証）
- メールアドレス形式は RFC 5322 に準拠する必要があります

**リレーション:**

- `workspaces` (Many-to-Many): このユーザーが所属するワークスペース

**検証:**

- メールアドレス形式検証（正規表現: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`）
- パスワードハッシュは空であってはならない

**フック:**

- `BeforeCreate` / `BeforeUpdate`: メールアドレスを小文字に正規化

**セキュリティ考慮事項:**

- `password_hash` フィールドは `Sensitive()` としてマークされ、ログに出力されません
- パスワードのハッシュ化には bcrypt (cost factor 10以上) を使用
- メールアドレスの一意性はデータベース制約とアプリケーション層の両方で強制

---

### 2. Workspace（ワークスペース）

**目的**: ユーザーまたはチームの隔離された作業環境を表現します。

**属性:**

| 属性名       | 型     | 制約               | 説明                     |
| ------------ | ------ | ------------------ | ------------------------ |
| `id`         | `int`  | Primary Key, Auto  | ワークスペースの一意識別子 |
| `name`       | `string` | Not Empty        | ワークスペース名（ユーザー提供） |
| `created_at` | `time` | Auto-set on create | ワークスペース作成日時     |
| `updated_at` | `time` | Auto-update        | 最終更新日時               |

**ビジネスルール:**

- ワークスペース名は**必須**です
- ワークスペース名には**任意の文字**を含めることができます（特殊文字、絵文字など）
- ワークスペース名に**一意性制約はありません**（複数のワークスペースが同じ名前を持つことが可能）
- サインアップ時、ユーザーは1つのワークスペースを取得します

**リレーション:**

- `users` (Many-to-Many): このワークスペースにアクセスできるユーザー

**検証:**

- 名前は空であってはならない

---

### 3. User-Workspace Relationship（ユーザー・ワークスペース関係）

**タイプ**: Many-to-Many

**実装**: Ent ORM が自動的に中間テーブルを作成します。

**中間テーブル詳細**（Ent により自動生成）:

| カラム          | 型    | 制約                              | 説明                   |
| --------------- | ----- | --------------------------------- | ---------------------- |
| `user_id`       | `int` | Foreign Key (users.id), Not Null  | ユーザーへの参照       |
| `workspace_id`  | `int` | Foreign Key (workspaces.id), Not Null | ワークスペースへの参照 |

**Primary Key**: (`user_id`, `workspace_id`) の複合キー

**インデックス:**

- `user_id` にインデックス（効率的な user → workspaces クエリ）
- `workspace_id` にインデックス（効率的な workspace → users クエリ）

**ビジネスルール:**

- サインアップフロー: 1人のユーザーは初期状態で1つのワークスペースにリンクされます
- 将来: ワークスペースごとに複数のユーザー、ユーザーごとに複数のワークスペースをサポート可能

---

## リレーションシップと多重度

```
User ◄──────────► Workspace
     (many-to-many)

- 1人のユーザーは複数のワークスペースを持つことができます
- 1つのワークスペースは複数のユーザーを持つことができます
- サインアップ時: 1人のユーザーは初期状態で1つのワークスペースにリンクされます
```

**シナリオ例:**

1. **サインアップフロー**:
   - ユーザーがサインアップ → 1 User + 1 Workspace + 1 リンク を作成

2. **将来のマルチワークスペース** (スコープ外):
   - ユーザーが追加のワークスペースを作成 → 1 Workspace + 1 リンク を作成
   - ユーザーがワークスペースに招待される → 1 リンク を作成（新しい User や Workspace は作成しない）

---

## 状態遷移

### User エンティティの状態

```
[存在しない]
    → (サインアップ) → [アクティブユーザー]
```

**備考:**

- MVP では、アカウント有効化、停止、削除の状態はありません
- 将来: `active`, `suspended`, `deleted` の値を持つ `status` フィールドを追加可能

### Workspace エンティティの状態

```
[存在しない]
    → (ユーザーサインアップ) → [アクティブワークスペース]
```

**備考:**

- ワークスペースはサインアップ時にユーザーと同時に作成されます
- MVP では、ワークスペースのアーカイブや削除はありません
- 将来: `status` フィールドまたは `deleted_at` でソフトデリートを追加可能

---

## データベーススキーマ（生成されるSQL）

### Users テーブル

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- メールアドレス一意性のインデックス（UNIQUE 制約により自動作成）
CREATE UNIQUE INDEX users_email_key ON users (email);
```

### Workspaces テーブル

```sql
CREATE TABLE workspaces (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### User-Workspace 中間テーブル

```sql
CREATE TABLE user_workspaces (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workspace_id INTEGER NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, workspace_id)
);

-- 効率的な検索のためのインデックス
CREATE INDEX user_workspaces_user_id_idx ON user_workspaces (user_id);
CREATE INDEX user_workspaces_workspace_id_idx ON user_workspaces (workspace_id);
```

---

## データ検証まとめ

### バックエンド（Ent）

| フィールド           | 検証                   | 強制方法                |
| -------------------- | ---------------------- | ----------------------- |
| `users.email`        | 形式（正規表現）       | Ent バリデーター        |
| `users.email`        | 一意性                 | データベース制約 + Ent  |
| `users.email`        | 小文字化               | Ent フック              |
| `users.password_hash`| 空でない               | Ent バリデーター        |
| `workspaces.name`    | 空でない               | Ent バリデーター        |

### フロントエンド（Zod）

| フィールド       | 検証               | エラーメッセージ                           |
| ---------------- | ------------------ | ------------------------------------------ |
| `email`          | 形式（email）      | "有効なメールアドレスを入力してください"   |
| `password`       | 最小長（8文字）    | "パスワードは8文字以上である必要があります"|
| `workspaceName`  | 空でない           | "ワークスペース名を入力してください"       |

---

## TypeScript 型定義（フロントエンド）

**場所**: `client/src/entities/user/types.ts`

```typescript
export interface User {
  id: number
  email: string
  createdAt: string // ISO 8601 タイムスタンプ
  updatedAt: string
}

export interface Workspace {
  id: number
  name: string
  createdAt: string
  updatedAt: string
}

export interface SignupRequest {
  email: string
  password: string
  workspaceName: string
}

export interface SignupResponse {
  user: User
  workspace: Workspace
  message: string
}
```

**備考:**

- `password_hash` はフロントエンドに公開されません
- タイムスタンプは JSON で ISO 8601 文字列としてシリアライズされます

---

## クエリパターン（Ent）

### よく使われるクエリ

**1. ワークスペース付きユーザー作成（サインアップ）:**

```go
tx, _ := client.Tx(ctx)

user, _ := tx.User.Create().
    SetEmail("user@example.com").
    SetPasswordHash(hashedPassword).
    Save(ctx)

workspace, _ := tx.Workspace.Create().
    SetName("My Workspace").
    AddUsers(user).
    Save(ctx)

tx.Commit()
```

**2. メールアドレスの存在チェック:**

```go
exists, _ := client.User.Query().
    Where(user.Email(strings.ToLower(email))).
    Exist(ctx)
```

**3. ワークスペース付きユーザー取得:**

```go
u, _ := client.User.Query().
    Where(user.ID(userID)).
    WithWorkspaces().
    Only(ctx)
```

---

## セキュリティ考慮事項

### 1. パスワード保存

- 平文パスワードは**決して保存しない**
- cost factor ≥ 10 で bcrypt を使用
- `password_hash` フィールドは Ent で `Sensitive()` としてマーク（ログに出力されない）

### 2. メールアドレス一意性

- データベース制約が重複を防止
- フックによる大文字小文字非区別（すべてのメールアドレスを小文字で保存）
- レースコンディションを防止

### 3. 外部キー制約

- `ON DELETE CASCADE` によりユーザー/ワークスペース削除時にクリーンアップを保証
- 参照整合性を維持

### 4. 入力検証

- フロントエンド（Zod）とバックエンド（Ent）の両方で検証
- 多層防御アプローチ

---

## 関連ドキュメント

- **機能仕様**: [F-0001: ユーザーサインアップ](../features/F-0001-user-signup.md)
- **API仕様**: [POST /auth/signup](../../05-api-reference/rest/auth/POST-signup.md)
- **UI仕様**: [S-100: サインアップページ](../../03-ui-spec/screens/S-100-signup.md)
- **データベーステーブル**: [テーブル定義](../../04-system-design/backend/database/tables.md)
