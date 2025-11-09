# プロジェクト構成

## 概要

Finsight はモノレポ構成を採用したフルスタックアプリケーションです。
バックエンド (Go) とフロントエンド (React) が同一リポジトリで管理され、効率的な開発とデプロイを実現しています。

## ディレクトリ構造

```
finsight/                           # プロジェクトルート
├── backend/                        # バックエンドコードベース
│   ├── cmd/                        # アプリケーションエントリーポイント
│   │   └── server/
│   │       └── main.go             # メインアプリケーション
│   ├── internal/                   # プライベートコード (Go の internal パッケージ)
│   │   ├── domain/                 # ドメイン層 (オニオンアーキテクチャ)
│   │   │   ├── model/              # ドメインモデル
│   │   │   ├── repository/         # リポジトリインターフェース
│   │   │   └── service/            # ドメインサービス
│   │   ├── application/            # アプリケーション層
│   │   │   ├── usecase/            # ユースケース
│   │   │   └── dto/                # データ転送オブジェクト
│   │   └── infrastructure/         # インフラストラクチャ層
│   │       ├── http/               # HTTP レイヤー (Gin)
│   │       │   ├── handler/        # HTTP ハンドラー
│   │       │   ├── middleware/     # ミドルウェア
│   │       │   └── router/         # ルーティング
│   │       ├── persistence/        # データベースアクセス
│   │       │   ├── ent/            # ent スキーマ・生成コード
│   │       │   └── repository/     # リポジトリ実装
│   │       └── external/           # 外部 API クライアント
│   ├── test/                       # テストコード
│   │   ├── integration/            # 結合テスト用ヘルパー
│   │   │   ├── helper.go
│   │   │   └── testdata/
│   │   └── fixtures/               # テストデータ
│   ├── go.mod                      # Go モジュール定義
│   ├── go.sum                      # Go 依存関係ロック
│   └── docker-compose.test.yml     # テスト用 DB 構成
│
├── client/                         # フロントエンドコードベース
│   ├── src/                        # ソースコード (Feature-Sliced Design)
│   │   ├── app/                    # アプリケーション層 (FSD)
│   │   │   ├── routes/             # TanStack Start ルーティング
│   │   │   │   ├── __root.tsx
│   │   │   │   ├── index.tsx
│   │   │   │   └── dashboard/
│   │   │   ├── providers/          # グローバルプロバイダー
│   │   │   │   ├── query-provider.tsx
│   │   │   │   └── theme-provider.tsx
│   │   │   └── styles/             # グローバルスタイル
│   │   │       └── globals.css
│   │   ├── pages/                  # ページ層 (FSD)
│   │   │   ├── dashboard/
│   │   │   ├── login/
│   │   │   └── report/
│   │   ├── widgets/                # ウィジェット層 (FSD)
│   │   │   ├── header/
│   │   │   ├── portfolio-summary/
│   │   │   └── transaction-list/
│   │   ├── features/               # フィーチャー層 (FSD)
│   │   │   ├── auth/
│   │   │   │   ├── login/
│   │   │   │   └── logout/
│   │   │   ├── transaction/
│   │   │   │   ├── add-transaction/
│   │   │   │   └── edit-transaction/
│   │   │   └── report/
│   │   ├── entities/               # エンティティ層 (FSD)
│   │   │   ├── user/
│   │   │   │   ├── model/
│   │   │   │   ├── api/
│   │   │   │   └── ui/
│   │   │   ├── transaction/
│   │   │   └── portfolio/
│   │   └── shared/                 # 共有層 (FSD)
│   │       ├── ui/                 # shadcn/ui コンポーネント
│   │       ├── lib/                # ユーティリティ
│   │       ├── api/                # API クライアント基盤
│   │       ├── config/             # 設定
│   │       └── types/              # 共通型定義
│   ├── public/                     # 静的ファイル
│   ├── package.json                # npm パッケージ定義
│   ├── tsconfig.json               # TypeScript 設定
│   ├── tailwind.config.ts          # Tailwind CSS 設定
│   ├── vite.config.ts              # Vite 設定
│   └── .env.example                # 環境変数サンプル
│
├── docs/                           # ドキュメント
│   ├── 00-overview/                # プロジェクト概要
│   ├── 01-requirements/            # 要件定義
│   ├── 03-ui-spec/                 # UI 仕様
│   ├── 04-system-design/           # システム設計
│   │   ├── backend/
│   │   │   └── architecture.md
│   │   ├── client/
│   │   │   └── architecture.md
│   │   ├── database-design.md
│   │   └── project-structure.md    # このファイル
│   ├── 05-api-reference/           # API リファレンス
│   ├── 06-operation/               # 運用
│   ├── 07-testing/                 # テスト
│   └── 08-release-notes/           # リリースノート
│
├── .github/                        # GitHub 設定
│   └── workflows/                  # CI/CD ワークフロー
├── .gitignore                      # Git 除外設定
├── README.md                       # プロジェクト README
└── TASKS.md                        # タスク管理
```

## 各ディレクトリの詳細

### `/backend` - バックエンド

**技術スタック**:

- 言語: Go
- フレームワーク: Gin
- ORM: ent
- データベース: PostgreSQL

**アーキテクチャ**: オニオンアーキテクチャ

詳細は [バックエンドアーキテクチャ](backend/architecture.md) を参照してください。

#### `/backend/cmd`

アプリケーションのエントリーポイントを配置します。

- `cmd/server/main.go`: メインアプリケーション
- 将来的に CLI ツールなどを追加する場合も `cmd/` 配下に配置

#### `/backend/internal`

プライベートなアプリケーションコード。Go の `internal` パッケージ規約により、外部からインポートできません。

- **domain/**: ビジネスロジックの中核
- **application/**: ユースケースの実装
- **infrastructure/**: 外部とのやり取り (HTTP, DB, 外部 API)

#### `/backend/test`

テストコードとテストデータを配置します。

- **integration/**: 結合テスト用のヘルパー関数
- **fixtures/**: テストデータ (JSON, SQL など)

### `/client` - フロントエンド

**技術スタック**:

- UI ライブラリ: React v19
- フレームワーク: TanStack Start
- UI コンポーネント: shadcn/ui
- スタイリング: Tailwind CSS
- 状態管理: React Query, TanStack Form

**アーキテクチャ**: Feature-Sliced Design (FSD)

詳細は [フロントエンドアーキテクチャ](client/architecture.md) を参照してください。

#### `/client/src`

ソースコードのルート。Feature-Sliced Design に従った階層構造です。

**FSD の階層** (上位 → 下位):

1. **app/**: アプリケーション初期化、ルーティング、グローバル設定
2. **pages/**: ページコンポーネント
3. **widgets/**: 独立した UI ブロック
4. **features/**: ユーザーアクションやビジネス機能
5. **entities/**: ビジネスエンティティ
6. **shared/**: 汎用的な共有コード

**依存関係のルール**: 下位レイヤーは上位レイヤーに依存できません。

### `/docs` - ドキュメント

プロジェクトの全ドキュメントを一元管理します。

- **00-overview/**: プロジェクトの概要
- **01-requirements/**: 要件定義書
- **03-ui-spec/**: UI/UX 仕様
- **04-system-design/**: システム設計書
- **05-api-reference/**: API 仕様書
- **06-operation/**: デプロイ・運用手順
- **07-testing/**: テスト計画・手順
- **08-release-notes/**: リリースノート

## モノレポ構成の利点

### 1. コードの一元管理

- バックエンドとフロントエンドの変更を同一 PR で管理
- 型定義の共有が容易 (OpenAPI スキーマから生成など)
- リファクタリングの影響範囲が明確

### 2. 開発効率の向上

- 単一リポジトリでのクローン・セットアップ
- 統一された開発環境
- フルスタックの変更を一括でテスト可能

### 3. デプロイの柔軟性

- バックエンドとフロントエンドを独立してデプロイ可能
- Docker Compose で簡単にローカル環境を構築
- CI/CD パイプラインの統合管理

## 開発フロー

### ローカル開発

#### バックエンド

```bash
cd backend
go run cmd/server/main.go
```

#### フロントエンド

```bash
cd client
npm install
npm run dev
```

### テスト実行

#### バックエンド

```bash
cd backend
# テスト用 DB 起動
docker-compose -f docker-compose.test.yml up -d

# テスト実行
go test ./... -v

# テスト用 DB 停止
docker-compose -f docker-compose.test.yml down
```

#### フロントエンド

```bash
cd client
# ユニットテスト・コンポーネントテスト
npm run test

# E2E テスト
npm run test:e2e
```

## 命名規則

### バックエンド (Go)

- **パッケージ名**: 小文字、単数形 (`user`, `transaction`)
- **ファイル名**: スネークケース (`user_handler.go`, `calculate_roi.go`)
- **関数/メソッド**: PascalCase (公開) または camelCase (非公開)
- **テストファイル**: `*_test.go`

### フロントエンド (TypeScript/React)

- **コンポーネント**: PascalCase (`UserCard`, `LoginForm`)
- **ファイル名**: kebab-case (`user-card.tsx`, `login-form.tsx`)
- **フック**: `use` + PascalCase (`useUser`, `useLogin`)
- **ユーティリティ関数**: camelCase (`formatCurrency`, `parseDate`)
- **型定義**: PascalCase (`User`, `Transaction`)

## 環境変数管理

### バックエンド

- `.env.example` をコピーして `.env` を作成
- 本番環境では環境変数として直接設定

### フロントエンド

- `.env.example` をコピーして `.env` を作成
- `VITE_` プレフィックスで始まる変数のみクライアントに公開
- 機密情報は含めない (API キーなどはバックエンド経由で取得)

## バージョン管理

### Git ブランチ戦略

- `main`: 本番環境デプロイ用
- `develop`: 開発環境 (オプション)
- `feature/*`: 機能開発
- `bugfix/*`: バグ修正

### コミットメッセージ

```
<type>(<scope>): <subject>

例:
feat(backend): add user authentication
fix(client): resolve form validation issue
docs: update architecture documentation
```

## 関連ドキュメント

- [プロジェクト概要](../00-overview/project-overview.md)
- [バックエンドアーキテクチャ](backend/architecture.md)
- [フロントエンドアーキテクチャ](client/architecture.md)
- [データベース設計](database-design.md)
- [API リファレンス](../05-api-reference/)
