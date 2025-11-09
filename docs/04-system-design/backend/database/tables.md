# テーブル定義

## 1. users（ユーザー）

ユーザーのアカウント情報を管理するテーブル

| カラム名      | データ型  | NULL | デフォルト | 制約                     | 説明                              |
| ------------- | --------- | ---- | ---------- | ------------------------ | --------------------------------- |
| id            | BIGINT    | NO   | AUTO       | PRIMARY KEY              | ユーザーID                        |
| email         | VARCHAR   | NO   | -          | UNIQUE                   | メールアドレス                    |
| password_hash | VARCHAR   | NO   | -          | -                        | パスワードハッシュ (bcrypt等)     |
| created_at    | TIMESTAMP | NO   | NOW()      | -                        | 作成日時                          |
| updated_at    | TIMESTAMP | NO   | NOW()      | -                        | 更新日時                          |

**インデックス:**
- PRIMARY KEY: `id`
- UNIQUE INDEX: `email`

**備考:**
- パスワードは必ずハッシュ化して保存すること
- セキュリティ要件に基づき、適切なパスワード強度チェックを実施

---

## 2. workspaces（ワークスペース）

家計簿データを管理する単位となるワークスペース情報を管理するテーブル

| カラム名   | データ型  | NULL | デフォルト | 制約        | 説明             |
| ---------- | --------- | ---- | ---------- | ----------- | ---------------- |
| id         | BIGINT    | NO   | AUTO       | PRIMARY KEY | ワークスペースID |
| name       | VARCHAR   | NO   | -          | -           | ワークスペース名 |
| created_at | TIMESTAMP | NO   | NOW()      | -           | 作成日時         |
| updated_at | TIMESTAMP | NO   | NOW()      | -           | 更新日時         |

**インデックス:**
- PRIMARY KEY: `id`

**備考:**
- ワークスペースは複数のユーザーがメンバーとして所属可能
- すべての家計簿データはワークスペース単位で管理される

---

## 3. workspace_members（ワークスペースメンバー）

ユーザーとワークスペースの所属関係を管理する中間テーブル

| カラム名     | データ型    | NULL | デフォルト | 制約                        | 説明                                   |
| ------------ | ----------- | ---- | ---------- | --------------------------- | -------------------------------------- |
| id           | BIGINT      | NO   | AUTO       | PRIMARY KEY                 | ID                                     |
| workspace_id | BIGINT      | NO   | -          | FOREIGN KEY (workspaces.id) | ワークスペースID                       |
| user_id      | BIGINT      | NO   | -          | FOREIGN KEY (users.id)      | ユーザーID                             |
| role         | VARCHAR(20) | NO   | 'member'   | CHECK (role IN ('owner', 'member')) | 役割 (owner: オーナー, member: メンバー) |
| created_at   | TIMESTAMP   | NO   | NOW()      | -                           | 作成日時                               |

**インデックス:**
- PRIMARY KEY: `id`
- UNIQUE INDEX: `workspace_id, user_id`
- INDEX: `user_id`
- INDEX: `workspace_id`

**備考:**
- 1ユーザーは複数のワークスペースに所属可能
- 1ワークスペースには複数のユーザーが所属可能（多対多）
- ワークスペース作成時に、作成者は自動的に 'owner' 役割で登録される
- 役割による権限制御はアプリケーション層で実装

---

## 4. transactions（収支データ）

個別の収入・支出記録を管理するテーブル

| カラム名          | データ型       | NULL | デフォルト | 制約                        | 説明                                      |
| ----------------- | -------------- | ---- | ---------- | --------------------------- | ----------------------------------------- |
| id                | BIGINT         | NO   | AUTO       | PRIMARY KEY                 | 収支データID                              |
| workspace_id      | BIGINT         | NO   | -          | FOREIGN KEY (workspaces.id) | ワークスペースID                          |
| transaction_date  | DATE           | NO   | -          | -                           | 取引日                                    |
| amount            | DECIMAL(15, 2) | NO   | -          | CHECK (amount >= 0)         | 金額（正の数値）                          |
| type              | VARCHAR(10)    | NO   | -          | CHECK (type IN ('income', 'expense')) | 収支区分 (income: 収入, expense: 支出)    |
| category_id       | BIGINT         | YES  | NULL       | FOREIGN KEY (categories.id) | カテゴリID（任意）                        |
| memo              | TEXT           | YES  | NULL       | -                           | メモ（任意）                              |
| created_at        | TIMESTAMP      | NO   | NOW()      | -                           | 作成日時                                  |
| updated_at        | TIMESTAMP      | NO   | NOW()      | -                           | 更新日時                                  |

**インデックス:**
- PRIMARY KEY: `id`
- INDEX: `workspace_id`
- INDEX: `workspace_id, transaction_date`
- INDEX: `workspace_id, type`
- INDEX: `category_id`

**備考:**
- データの隔離性を保つため、すべてのクエリで workspace_id による権限チェックを実施
- 削除時は論理削除ではなく物理削除を想定

---

## 5. categories（カテゴリマスタ）

収支の分類を管理するマスタテーブル

| カラム名     | データ型    | NULL | デフォルト | 制約                        | 説明                                           |
| ------------ | ----------- | ---- | ---------- | --------------------------- | ---------------------------------------------- |
| id           | BIGINT      | NO   | AUTO       | PRIMARY KEY                 | カテゴリID                                     |
| workspace_id | BIGINT      | NO   | -          | FOREIGN KEY (workspaces.id) | ワークスペースID                               |
| name         | VARCHAR     | NO   | -          | -                           | カテゴリ名                                     |
| type         | VARCHAR(10) | NO   | -          | CHECK (type IN ('income', 'expense')) | カテゴリ種別 (income: 収入用, expense: 支出用) |
| created_at   | TIMESTAMP   | NO   | NOW()      | -                           | 作成日時                                       |
| updated_at   | TIMESTAMP   | NO   | NOW()      | -                           | 更新日時                                       |

**インデックス:**
- PRIMARY KEY: `id`
- INDEX: `workspace_id`
- UNIQUE INDEX: `workspace_id, name, type`

**備考:**
- 同一ワークスペース内で、同じタイプ(income/expense)のカテゴリ名は重複不可
- ワークスペース単位でカテゴリを管理し、メンバー間で共有される

---

## 6. repeated_transactions（繰り返しトランザクション）

定期的に発生する収支データを管理するテーブル

| カラム名          | データ型       | NULL | デフォルト | 制約                        | 説明                                      |
| ----------------- | -------------- | ---- | ---------- | --------------------------- | ----------------------------------------- |
| id                | BIGINT         | NO   | AUTO       | PRIMARY KEY                 | 繰り返しトランザクションID                |
| workspace_id      | BIGINT         | NO   | -          | FOREIGN KEY (workspaces.id) | ワークスペースID                          |
| started_at        | DATE           | NO   | -          | -                           | 開始日（必須）                            |
| finished_at       | DATE           | NO   | '2999-12-31' | -                         | 終了日（デフォルトは2999-12-31で無期限を表現）|
| transaction_date  | DATE           | NO   | -          | -                           | 取引日                                    |
| amount            | DECIMAL(15, 2) | NO   | -          | CHECK (amount >= 0)         | 金額（正の数値）                          |
| type              | VARCHAR(10)    | NO   | -          | CHECK (type IN ('income', 'expense')) | 収支区分 (income: 収入, expense: 支出)    |
| category_id       | BIGINT         | YES  | NULL       | FOREIGN KEY (categories.id) | カテゴリID（任意）                        |
| memo              | TEXT           | YES  | NULL       | -                           | メモ（任意）                              |
| created_at        | TIMESTAMP      | NO   | NOW()      | -                           | 作成日時                                  |
| updated_at        | TIMESTAMP      | NO   | NOW()      | -                           | 更新日時                                  |

**インデックス:**
- PRIMARY KEY: `id`
- INDEX: `workspace_id`
- INDEX: `workspace_id, started_at, finished_at`
- INDEX: `category_id`

**備考:**
- started_at: 繰り返しトランザクションの開始日（必須）
- finished_at: 繰り返しトランザクションの終了日（アプリ的にはオプショナルだが、DBレベルでは NOT NULL で、デフォルト値 '2999-12-31' を設定して無期限を表現）
- transaction_date: 実際の取引日（毎月の何日に発生するかを表現）
- アプリケーション側では、finished_at が '2999-12-31' の場合は「終了日なし（無期限）」として扱う
- 繰り返しトランザクションは transactions テーブルとは独立して管理される

---

## 7. csv_templates（CSV テンプレート）

CSV ファイルのマッピング設定を管理するテーブル

| カラム名        | データ型  | NULL | デフォルト | 制約                        | 説明                         |
| --------------- | --------- | ---- | ---------- | --------------------------- | ---------------------------- |
| id              | BIGINT    | NO   | AUTO       | PRIMARY KEY                 | テンプレートID               |
| workspace_id    | BIGINT    | NO   | -          | FOREIGN KEY (workspaces.id) | ワークスペースID             |
| template_name   | VARCHAR   | NO   | -          | -                           | テンプレート名               |
| column_mappings | JSONB     | NO   | -          | -                           | 列マッピング設定（JSON形式） |
| created_at      | TIMESTAMP | NO   | NOW()      | -                           | 作成日時                     |
| updated_at      | TIMESTAMP | NO   | NOW()      | -                           | 更新日時                     |

**インデックス:**
- PRIMARY KEY: `id`
- INDEX: `workspace_id`
- UNIQUE INDEX: `workspace_id, template_name`

**column_mappings JSON スキーマ例:**
```json
{
  "dateColumn": {
    "index": 0,
    "format": "YYYY-MM-DD"
  },
  "amountColumn": {
    "index": 1
  },
  "typeColumn": {
    "index": 2,
    "mapping": {
      "入金": "income",
      "出金": "expense"
    }
  },
  "categoryColumn": {
    "index": 3,
    "defaultValue": null
  },
  "memoColumn": {
    "index": 4
  }
}
```

**備考:**
- column_mappings には柔軟なマッピング設定を JSON 形式で保存
- 同一ワークスペース内でテンプレート名は重複不可
- ワークスペース単位でテンプレートを管理し、メンバー間で共有される

---

## 8. reports（レポート設定）

保存されたレポート設定を管理するテーブル

| カラム名      | データ型  | NULL | デフォルト | 制約                        | 説明                     |
| ------------- | --------- | ---- | ---------- | --------------------------- | ------------------------ |
| id            | BIGINT    | NO   | AUTO       | PRIMARY KEY                 | レポートID               |
| workspace_id  | BIGINT    | NO   | -          | FOREIGN KEY (workspaces.id) | ワークスペースID         |
| report_name   | VARCHAR   | NO   | -          | -                           | レポート名               |
| report_config | JSONB     | NO   | -          | -                           | レポート設定（JSON形式） |
| created_at    | TIMESTAMP | NO   | NOW()      | -                           | 作成日時                 |
| updated_at    | TIMESTAMP | NO   | NOW()      | -                           | 更新日時                 |

**インデックス:**
- PRIMARY KEY: `id`
- INDEX: `workspace_id`
- UNIQUE INDEX: `workspace_id, report_name`

**report_config JSON スキーマ例:**
```json
{
  "period": {
    "startYearMonth": "2025-01",
    "endYearMonth": "2025-12"
  },
  "displayItems": {
    "showIncome": true,
    "showExpense": true,
    "groupByCategory": true,
    "groupByAttribute": false,
    "separateRepeatedVariable": true
  },
  "chartType": "line",
  "aggregationPeriod": "monthly"
}
```

**備考:**
- report_config には表示項目、期間、グラフタイプなどの設定を JSON 形式で保存
- 同一ワークスペース内でレポート名は重複不可
- ワークスペース単位でレポート設定を管理し、メンバー間で共有される

---

## データ整合性・制約

### 外部キー制約

すべての外部キーに対して、以下の制約を設定:
- **ON DELETE**:
  - users テーブルの削除時: CASCADE（ユーザーに紐づく workspace_members を削除）
  - workspaces テーブルの削除時: CASCADE（ワークスペースに紐づくすべてのデータを削除）
  - categories の削除時: SET NULL または RESTRICT（参照されている場合は削除不可）

### CHECK 制約

- **transactions.amount**: `amount >= 0` (負の数値不可)
- **transactions.type**: `type IN ('income', 'expense')`
- **categories.type**: `type IN ('income', 'expense')`
- **repeated_transactions.amount**: `amount >= 0`
- **repeated_transactions.type**: `type IN ('income', 'expense')`
- **workspace_members.role**: `role IN ('owner', 'member')`

### UNIQUE 制約

- **users.email**: メールアドレスは一意
- **workspace_members**: `(workspace_id, user_id)` の組み合わせは一意
- **categories**: `(workspace_id, name, type)` の組み合わせは一意
- **csv_templates**: `(workspace_id, template_name)` の組み合わせは一意
- **reports**: `(workspace_id, report_name)` の組み合わせは一意

---

## セキュリティ考慮事項

### データの隔離性

**重要**: すべてのクエリにおいて、以下の2段階の権限チェックを必須とする

1. **ワークスペースメンバーシップ確認**:
   - ユーザーが操作対象のワークスペースのメンバーであることを workspace_members テーブルで確認
   - メンバーでない場合はアクセスを拒否

2. **workspace_id によるフィルタリング**:
   - すべてのデータアクセスで、workspace_id によるフィルタリングを実施
   - 他のワークスペースのデータへのアクセスを完全に防止
   - ORM (Ent) のフックやミドルウェアを活用して、自動的に workspace_id フィルタを適用

3. **実装方針**:
   - API リクエストごとに、セッションユーザーと workspace_id の組み合わせをチェック
   - ワークスペース選択は明示的に行われ、現在のワークスペースコンテキストをセッションまたはリクエストヘッダーで管理

### パスワード管理

- パスワードは必ず bcrypt 等の安全なハッシュアルゴリズムを使用
- 最小パスワード長、複雑性の要件を設定

### SQL インジェクション対策

- ORM (Ent) のパラメータバインディング機能を使用
- 動的 SQL の直接実行を避ける

---

## パフォーマンス考慮事項

### インデックス設計

- よく検索される列にインデックスを設定（workspace_id, user_id, transaction_date, type など）
- 複合インデックスを活用して、複数条件の検索を最適化
- JSONB 型のカラム（column_mappings, report_config）には必要に応じて GIN インデックスを検討

### N+1 問題の回避

- Ent の Eager Loading 機能を使用して、必要なリレーションを事前にロード
- 集計クエリでは適切に JOIN を使用

### ページネーション

- 大量データ表示時は LIMIT/OFFSET を使用
- カーソルベースのページネーションも検討
