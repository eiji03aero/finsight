# バックエンドアーキテクチャ

## 概要

Finsight のバックエンドは、保守性と拡張性を重視したクリーンな設計を採用しています。
オニオンアーキテクチャをベースに、ドメイン駆動設計の原則に従った構成となっています。

## 技術スタック

### 言語・フレームワーク

- **言語**: Go (Golang)
- **HTTP サーバー**: Gin Web Framework
- **ORM**: ent (Facebook 製)

### データベース

- **RDBMS**: PostgreSQL

### その他主要ライブラリ

- (必要に応じて追加)

## プロジェクト構成

本プロジェクトはモノレポ構成を採用しており、バックエンドは `backend/` ディレクトリに配置されています。

```
finsight/                    # プロジェクトルート
├── backend/                 # バックエンドコードベース
│   ├── cmd/                 # アプリケーションエントリーポイント
│   │   └── server/
│   │       └── main.go      # メインアプリケーション
│   ├── internal/            # プライベートコード
│   │   ├── domain/          # ドメイン層
│   │   ├── application/     # アプリケーション層
│   │   └── infrastructure/  # インフラストラクチャ層
│   ├── test/                # テストコード
│   │   ├── integration/     # 結合テスト
│   │   └── fixtures/        # テストデータ
│   ├── go.mod
│   ├── go.sum
│   └── docker-compose.test.yml
├── client/                  # フロントエンドコードベース
└── docs/                    # ドキュメント
```

### Go プロジェクトのディレクトリ構成

バックエンドは標準的な Go プロジェクトレイアウトに従っています:

- **`cmd/`**: アプリケーションのエントリーポイント
- **`internal/`**: プライベートなアプリケーションコード (外部からインポート不可)
- **`test/`**: テストコードとテストデータ

## アーキテクチャ設計

### オニオンアーキテクチャ

本システムでは、オニオンアーキテクチャを採用し、以下の層構造で実装しています:

```
┌─────────────────────────────────────┐
│   Infrastructure Layer (外側)        │
│  - HTTPハンドラー (Gin)              │
│  - データベースアクセス (ent)         │
│  - 外部API クライアント              │
│  ┌───────────────────────────────┐  │
│  │   Application Layer           │  │
│  │  - ユースケース               │  │
│  │  - アプリケーションサービス    │  │
│  │  ┌─────────────────────────┐ │  │
│  │  │   Domain Layer (中心)   │ │  │
│  │  │  - ドメインモデル       │ │  │
│  │  │  - ドメインサービス     │ │  │
│  │  │  - リポジトリインターフェース │ │
│  │  └─────────────────────────┘ │  │
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘
```

#### 1. Domain Layer (ドメイン層)

最も内側の層で、ビジネスロジックの中核を担います。

**責務**:

- ドメインモデルの定義
- ビジネスルールの実装
- ドメインサービスの実装
- リポジトリのインターフェース定義

**特徴**:

- 外部の層に依存しない
- フレームワークやライブラリに依存しない純粋な Go コード
- ビジネスロジックの変更に強い

**ディレクトリ構成例**:

```
domain/
├── model/          # ドメインモデル
├── repository/     # リポジトリインターフェース
└── service/        # ドメインサービス
```

#### 2. Application Layer (アプリケーション層)

ユースケースを実装する層です。

**責務**:

- ユースケースの実装
- トランザクション管理
- ドメインオブジェクトの調整

**特徴**:

- ドメイン層のみに依存
- インフラストラクチャ層の詳細を知らない

**ディレクトリ構成例**:

```
application/
├── usecase/        # ユースケース実装
└── dto/            # データ転送オブジェクト
```

#### 3. Infrastructure Layer (インフラストラクチャ層)

外部とのやり取りを担当する層です。

**責務**:

- HTTP リクエスト/レスポンスの処理 (Gin)
- データベースアクセス (ent)
- 外部 API との通信
- リポジトリインターフェースの実装

**ディレクトリ構成例**:

```
infrastructure/
├── http/           # HTTPハンドラー、ルーティング (Gin)
│   ├── handler/
│   ├── middleware/
│   └── router/
├── persistence/    # データベースアクセス (ent)
│   ├── ent/        # entスキーマ・生成コード
│   └── repository/ # リポジトリ実装
└── external/       # 外部API クライアント
```

### 依存関係の方向

```
Infrastructure → Application → Domain
       ↓              ↓
    外側から内側への依存のみ許可
    (内側は外側を知らない)
```

## 主要コンポーネント

### Gin (HTTP サーバー)

**選定理由**:

- 高速な HTTP ルーター
- ミドルウェアのサポート
- 豊富なドキュメントとコミュニティ

**主な用途**:

- API エンドポイントの定義
- リクエストバリデーション
- レスポンスのシリアライゼーション
- 認証・認可ミドルウェア

### ent (ORM)

**選定理由**:

- スキーマファーストな設計
- 型安全なクエリビルダー
- マイグレーション管理
- グラフ構造のデータモデルに対応

**主な用途**:

- データベーススキーマの定義
- エンティティの CRUD 操作
- トランザクション管理
- データベースマイグレーション

### PostgreSQL

**選定理由**:

- ACID 準拠の信頼性
- 豊富な機能セット
- 拡張性とパフォーマンス

## 設計方針

### 1. 依存性逆転の原則 (DIP)

- ドメイン層はインターフェースを定義
- インフラストラクチャ層が具体的な実装を提供
- 依存性注入 (DI) による疎結合

### 2. 単一責任の原則 (SRP)

- 各層、各コンポーネントは明確な責務を持つ
- ビジネスロジックとインフラストラクチャの分離

### 3. テスタビリティ

- インターフェースを通じたモック化が容易
- ドメイン層は外部依存なしでテスト可能
- 各層で独立したユニットテストが可能

### 4. 保守性

- レイヤー間の明確な境界
- フレームワークやライブラリの変更に強い
- ビジネスロジックの可読性と再利用性

## データフロー

### 一般的なリクエストフロー

```
1. HTTPリクエスト
   ↓
2. Ginハンドラー (Infrastructure Layer)
   - リクエストのバリデーション
   - DTOへの変換
   ↓
3. ユースケース (Application Layer)
   - ビジネスロジックの調整
   - トランザクション管理
   ↓
4. ドメインサービス/リポジトリ (Domain Layer)
   - ビジネスルールの実行
   - データ操作の抽象インターフェース
   ↓
5. リポジトリ実装 (Infrastructure Layer)
   - entを使用したデータベースアクセス
   ↓
6. レスポンス生成
   - ドメインモデル → DTO → JSONレスポンス
```

## セキュリティ考慮事項

- (認証・認可の方式)
- (API キー管理)
- (入力バリデーション)
- (SQL インジェクション対策 - ent による型安全性)

## パフォーマンス考慮事項

- (データベース接続プーリング)
- (クエリ最適化)
- (キャッシング戦略)

## エラーハンドリング

- (エラー型の定義)
- (エラーのラッピングと伝播)
- (ログ出力方針)

## テスト方針

本システムでは、**結合テストを中心とした実践的なテスト戦略**を採用します。

### テスト戦略の基本方針

#### 1. 結合テスト (Integration Test) を主軸とする

**対象**: ハンドラーレベルからの E2E フロー
**範囲**: HTTP リクエスト → ハンドラー → ユースケース → リポジトリ → データベース
**目的**: 実際のユーザー操作に近い形で、システム全体の挙動を保証する

**特徴**:
- 実際の PostgreSQL データベースに接続した状態でテストを実行
- すべての層（Infrastructure → Application → Domain）���通過する
- 実環境に近い状態での動作を検証
- 基本的なユースケースはこのレベルのテストで網羅的にカバーする

**テスト対象の例**:
- ユーザー登録フロー全体
- 認証・認可を含むデータ取得
- トランザクションを伴う複数テーブルの更新
- エラーケース（バリデーションエラー、DB 制約違反など）

#### 2. ユニットテスト (Unit Test) で補完する

**対象**: 複雑なビジネスロジックやエッジケース
**範囲**: 個別の関数、メソッド、ドメインサービス
**目的**: 結合テストでカバーしきれない細かい条件分岐や計算ロジックを保証する

**特徴**:
- モック・スタブを活用した高速なテスト
- 外部依存を排除した純粋なロジックのテスト
- 必要に応じて実装（すべてをユニットテストでカバーしない）

**ユニットテストが有効なケース**:
- 複雑な計算ロジック（金融計算、統計処理など）
- 多岐にわたる条件分岐を持つドメインサービス
- エッジケースが多数存在するバリデーションロジック
- 外部 API との連携部分（モック化してテスト）

### テストの優先順位

```
1. 結合テスト (最優先)
   - 主要なユースケースを網羅
   - ハッピーパスとエラーパスの両方をカバー
   ↓
2. ユニットテスト (必要に応じて)
   - 結合テストで不十分な箇所を補完
   - 複雑なロジックの細かいケースをカバー
```

### テスト環境のセットアップ

#### データベース

**結合テスト用 DB**:
- テスト専用の PostgreSQL インスタンスを使用
- Docker Compose などで簡易に立ち上げ可能にする
- テストごとにトランザクションをロールバック、またはテーブルをクリーンアップ

**サンプル構成**:
```yaml
# docker-compose.test.yml
services:
  test-db:
    image: postgres:15
    environment:
      POSTGRES_DB: finsight_test
      POSTGRES_USER: test_user
      POSTGRES_PASSWORD: test_password
    ports:
      - "5433:5432"
```

#### テストヘルパー

```go
// テストごとにDBをクリーンアップするヘルパー
func setupTestDB(t *testing.T) *ent.Client {
    // DB接続
    // マイグレーション実行
    // t.Cleanup() でロールバックまたはクリーンアップ
}

// HTTPサーバーのテストヘルパー
func setupTestServer(t *testing.T) *httptest.Server {
    // Ginルーターのセットアップ
    // 依存性注入
    // テストサーバー起動
}
```

### ディレクトリ構成例

```
project-root/
├── internal/
│   ├── domain/
│   │   └── service/
│   │       ├── calculation.go
│   │       └── calculation_test.go      # ユニットテスト
│   ├── application/
│   │   └── usecase/
│   │       ├── user_usecase.go
│   │       └── user_usecase_test.go     # (必要に応じて)
│   └── infrastructure/
│       └── http/
│           └── handler/
│               ├── user_handler.go
│               └── user_handler_test.go # 結合テスト (メイン)
├── test/
│   ├── integration/           # 結合テスト用のヘルパー・共通処理
│   │   ├── helper.go
│   │   └── testdata/
│   └── fixtures/              # テストデータ
└── docker-compose.test.yml    # テスト用DB
```

### テスト実行例

**結合テストの実行**:
```bash
# テスト用DBを起動
docker-compose -f docker-compose.test.yml up -d

# 結合テストを実行
go test ./internal/infrastructure/http/handler/... -v

# テスト用DBを停止
docker-compose -f docker-compose.test.yml down
```

**ユニットテストの実行**:
```bash
# 特定のドメインサービスのユニットテスト
go test ./internal/domain/service/... -v
```

### テストコード例

#### 結合テストの例
```go
func TestUserHandler_CreateUser(t *testing.T) {
    // DB接続とテストサーバーのセットアップ
    db := setupTestDB(t)
    srv := setupTestServer(t, db)
    defer srv.Close()

    // テストケース
    tests := []struct {
        name           string
        requestBody    map[string]interface{}
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name: "正常系: ユーザー作成成功",
            requestBody: map[string]interface{}{
                "email": "test@example.com",
                "name":  "Test User",
            },
            expectedStatus: http.StatusCreated,
        },
        {
            name: "異常系: メールアドレス重複",
            requestBody: map[string]interface{}{
                "email": "duplicate@example.com",
                "name":  "Duplicate User",
            },
            expectedStatus: http.StatusConflict,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // HTTPリクエストを送信
            // レスポンスを検証
            // DBの状態を検証
        })
    }
}
```

#### ユニットテストの例
```go
func TestCalculationService_CalculateROI(t *testing.T) {
    // 複雑な計算ロジックの細かいケースをテスト
    tests := []struct {
        name          string
        initialAmount float64
        finalAmount   float64
        expectedROI   float64
    }{
        {"正の利益", 1000, 1200, 0.20},
        {"負の利益", 1000, 800, -0.20},
        {"利益ゼロ", 1000, 1000, 0.00},
        {"初期値ゼロのエッジケース", 0, 100, 0.00},
    }

    svc := NewCalculationService()
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            roi := svc.CalculateROI(tt.initialAmount, tt.finalAmount)
            assert.Equal(t, tt.expectedROI, roi)
        })
    }
}
```

### テストカバレッジの目標

- **結合テスト**: 主要なユースケースの 100% カバー
- **全体カバレッジ**: 70% 以上を目標（ただし、カバレッジの数値よりも実用的なテストケースを重視）
- **クリティカルパス**: 金融計算や認証など重要な箇所は必ずテストでカバー

### CI/CD での自動テスト

- すべてのプルリクエストで結合テストを自動実行
- テスト失敗時はマージをブロック
- テストカバレッジレポートを自動生成

## 関連ドキュメント

- [API 仕様](../../05-api-reference/)
- [デプロイメント構成](../../06-operation/deployment.md)
