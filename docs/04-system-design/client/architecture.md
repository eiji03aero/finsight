# フロントエンドアーキテクチャ

## 概要

Finsight のフロントエンドは、モダンな React エコシステムと Feature-Sliced Design を採用し、スケーラビリティと保守性を重視した設計となっています。

## 技術スタック

### コアライブラリ・フレームワーク

- **UI ライブラリ**: React v19
- **フレームワーク**: TanStack Start
- **UI コンポーネント**: shadcn/ui
- **スタイリング**: Tailwind CSS

### 状態管理・データフェッチング

- **フォーム管理**: @tanstack/react-form
- **データフェッチング**: @tanstack/react-query
- **HTTP クライアント**: Fetch API (Web 標準)

### 開発ツール

- **言語**: TypeScript
- **ビルドツール**: Vite (TanStack Start に内包)
- **リンター**: ESLint
- **フォーマッター**: Prettier

## プロジェクト構成

本プロジェクトはモノレポ構成を採用しており、フロントエンドは `client/` ディレクトリに配置されています。

```
finsight/                    # プロジェクトルート
├── backend/                 # バックエンドコードベース
├── client/                  # フロントエンドコードベース
│   ├── src/                 # ソースコード (Feature-Sliced Design)
│   │   ├── app/             # アプリケーション層
│   │   ├── pages/           # ページ層
│   │   ├── widgets/         # ウィジェット層
│   │   ├── features/        # フィーチャー層
│   │   ├── entities/        # エンティティ層
│   │   └── shared/          # 共有層
│   ├── public/              # 静的ファイル
│   ├── package.json
│   ├── tsconfig.json
│   ├── tailwind.config.ts
│   └── vite.config.ts
└── docs/                    # ドキュメント
```

## アーキテクチャ設計

### Feature-Sliced Design (FSD)

本システムでは、Feature-Sliced Design を採用し、以下の階層構造で実装しています。

#### FSD の基本構造

```
src/
├── app/              # アプリケーション全体の初期化・設定
├── pages/            # ページコンポーネント (ルーティング)
├── widgets/          # 独立した UI ブロック
├── features/         # ユーザーのアクションやビジネス機能
├── entities/         # ビジネスエンティティ
└── shared/           # 再利用可能な共通コード
```

#### 各レイヤーの責務

##### 1. App Layer (アプリケーション層)

**責務**: アプリケーション全体の初期化と設定

**内容**:

- ルーティング設定
- グローバルプロバイダー (React Query, テーマなど)
- グローバルスタイル
- アプリケーションエントリーポイント

**ディレクトリ構成例**:

```
app/
├── routes/           # TanStack Start ルーティング
│   ├── __root.tsx    # ルートレイアウト
│   ├── index.tsx     # トップページ
│   └── dashboard/
│       └── index.tsx
├── providers/        # グローバルプロバイダー
│   ├── query-provider.tsx
│   └── theme-provider.tsx
├── styles/           # グローバルスタイル
│   └── globals.css
└── router.tsx        # ルーター設定
```

##### 2. Pages Layer (ページ層)

**責務**: ルーティングとページレベルの構成

**内容**:

- URL に対応するページコンポーネント
- Widget の組み合わせ
- ページ固有のレイアウト
- SEO メタデータ

**ディレクトリ構成例**:

```
pages/
├── dashboard/
│   ├── ui/
│   │   └── dashboard-page.tsx
│   └── index.ts
├── login/
│   ├── ui/
│   │   └── login-page.tsx
│   └── index.ts
└── report/
    ├── ui/
    │   └── report-page.tsx
    └── index.ts
```

##### 3. Widgets Layer (ウィジェット層)

**責務**: 独立した大きな UI ブロック

**内容**:

- 複数の Feature や Entity を組み合わせた UI
- ページ内の独立したセクション
- ヘッダー、フッター、サイドバーなど

**ディレクトリ構成例**:

```
widgets/
├── header/
│   ├── ui/
│   │   └── header.tsx
│   └── index.ts
├── portfolio-summary/
│   ├── ui/
│   │   └── portfolio-summary.tsx
│   ├── model/
│   │   └── use-portfolio-data.ts
│   └── index.ts
└── transaction-list/
    ├── ui/
    │   └── transaction-list.tsx
    └── index.ts
```

##### 4. Features Layer (フィーチャー層)

**責務**: ユーザーのアクションやビジネス機能

**内容**:

- ユーザーが実行できるアクション (ログイン、データ登録など)
- ビジネスロジックを持つ UI コンポーネント
- フォーム、ボタン、モーダルなどのインタラクティブな要素

**ディレクトリ構成例**:

```
features/
├── auth/
│   ├── login/
│   │   ├── ui/
│   │   │   └── login-form.tsx
│   │   ├── model/
│   │   │   └── use-login.ts
│   │   └── index.ts
│   └── logout/
│       ├── ui/
│       │   └── logout-button.tsx
│       └── index.ts
├── transaction/
│   ├── add-transaction/
│   │   ├── ui/
│   │   │   └── add-transaction-form.tsx
│   │   ├── model/
│   │   │   └── use-add-transaction.ts
│   │   └── index.ts
│   └── edit-transaction/
│       └── ...
└── report/
    ├── generate-report/
    │   ├── ui/
    │   │   └── generate-report-button.tsx
    │   └── index.ts
    └── export-report/
        └── ...
```

##### 5. Entities Layer (エンティティ層)

**責務**: ビジネスエンティティの表現

**内容**:

- データモデルの型定義
- エンティティの表示用 UI コンポーネント
- エンティティに関する API クライアント
- 状態管理 (React Query のクエリ・ミューテーション)

**ディレクトリ構成例**:

```
entities/
├── user/
│   ├── model/
│   │   ├── types.ts
│   │   └── queries.ts       # React Query
│   ├── api/
│   │   └── user-api.ts
│   ├── ui/
│   │   ├── user-card.tsx
│   │   └── user-avatar.tsx
│   └── index.ts
├── transaction/
│   ├── model/
│   │   ├── types.ts
│   │   └── queries.ts
│   ├── api/
│   │   └── transaction-api.ts
│   ├── ui/
│   │   └── transaction-item.tsx
│   └── index.ts
└── portfolio/
    ├── model/
    │   ├── types.ts
    │   └── queries.ts
    ├── api/
    │   └── portfolio-api.ts
    └── index.ts
```

##### 6. Shared Layer (共有層)

**責務**: 再利用可能な汎用コード

**内容**:

- UI コンポーネント (shadcn/ui)
- ユーティリティ関数
- 型定義
- 定数
- API クライアントの基盤

**ディレクトリ構成例**:

```
shared/
├── ui/               # shadcn/ui コンポーネント
│   ├── button.tsx
│   ├── card.tsx
│   ├── dialog.tsx
│   └── ...
├── lib/              # ユーティリティ
│   ├── cn.ts         # classname ユーティリティ
│   ├── format.ts     # フォーマット関数
│   └── date.ts       # 日付処理
├── api/              # API クライアント基盤
│   ├── client.ts     # Fetch API ラッパー
│   └── types.ts      # 共通型
├── config/           # 設定
│   └── constants.ts
└── types/            # 共通型定義
    └── common.ts
```

### 依存関係のルール

Feature-Sliced Design では、**下位レイヤーは上位レイヤーに依存してはいけない**という厳格なルールがあります。

```
app → pages → widgets → features → entities → shared
 ↓      ↓        ↓         ↓          ↓
下位レイヤーへの依存のみ許可
(上位レイヤーへの依存は禁止)
```

**依存関係の例**:

- ✅ `features/auth/login` は `entities/user` を使える
- ✅ `widgets/header` は `features/auth/logout` を使える
- ❌ `entities/user` は `features/auth/login` を使えない
- ❌ `features` は `widgets` を使えない

## 主要コンポーネント

### React v19

**選定理由**:

- 最新の React 機能 (Server Components, Actions など)
- パフォーマンス改善
- 型安全性の向上

**主な機能**:

- React Server Components (RSC)
- use() フック (Promise/Context の読み取り)
- Form Actions
- Optimistic Updates

### TanStack Start

**選定理由**:

- フルスタック React フレームワーク
- ファイルベースルーティング
- サーバーサイドレンダリング (SSR) のサポート
- TanStack エコシステムとの統合

**主な機能**:

- ファイルベースルーティング
- Server Functions (API ルート)
- SSR/SSG サポート
- TypeScript ファーストな設計

### shadcn/ui

**選定理由**:

- アクセシブルで高品質な UI コンポーネント
- Radix UI ベースで堅牢
- カスタマイズ性が高い (コピー&ペーストで使用)
- Tailwind CSS との親和性

**主な特徴**:

- パッケージではなくコードとして管理
- プロジェクトに直接コピーして使用
- 完全にカスタマイズ可能

### Tailwind CSS

**選定理由**:

- ユーティリティファーストの CSS フレームワーク
- 高速な開発体験
- 一貫したデザインシステム
- ビルド時の未使用 CSS の削除

**主な用途**:

- コンポーネントのスタイリング
- レスポンシブデザイン
- ダークモード対応

### @tanstack/react-form

**選定理由**:

- 型安全なフォーム管理
- 柔軟なバリデーション
- パフォーマンス最適化 (部分的な再レンダリング)

**主な機能**:

- フィールドレベルのバリデーション
- 非同期バリデーション
- フォーム状態の管理
- TypeScript との統合

### @tanstack/react-query

**選定理由**:

- 宣言的なデータフェッチング
- 自動キャッシング・再検証
- Optimistic Updates
- デブ向けツールが充実

**主な機能**:

- データキャッシング
- バックグラウンド再検証
- ミューテーション管理
- 楽観的更新

### Fetch API

**選定理由**:

- Web 標準 API
- 追加のライブラリ不要
- React Server Components との親和性

**主な用途**:

- REST API との通信
- サーバーサイドでのデータフェッチ

## 設計方針

### 1. コンポーネント設計

**原則**:

- 単一責任の原則 (SRP) を遵守
- プレゼンテーショナルコンポーネントとコンテナコンポーネントの分離
- Props のインターフェースを明確に定義

**命名規則**:

- コンポーネント: PascalCase (`UserCard`, `LoginForm`)
- ファイル: kebab-case (`user-card.tsx`, `login-form.tsx`)
- フック: use から始まる camelCase (`useUser`, `useLogin`)

### 2. 型安全性

**原則**:

- すべてのコンポーネントで Props に型定義
- API レスポンスの型定義
- `any` の使用を避ける

**例**:

```typescript
// entities/user/model/types.ts
export interface User {
	id: string;
	email: string;
	name: string;
	createdAt: Date;
}

// features/auth/login/ui/login-form.tsx
interface LoginFormProps {
	onSuccess?: (user: User) => void;
	redirectTo?: string;
}

export function LoginForm({ onSuccess, redirectTo }: LoginFormProps) {
	// ...
}
```

### 3. 状態管理

**方針**:

- サーバー状態: React Query で管理
- クライアント状態: React の useState/useReducer で管理
- グローバル状態: Context API または Zustand (必要に応じて)

**状態の分類**:

```typescript
// サーバー状態 (React Query)
const { data: user } = useQuery({
	queryKey: ["user", userId],
	queryFn: () => fetchUser(userId),
});

// クライアント状態 (useState)
const [isOpen, setIsOpen] = useState(false);

// フォーム状態 (TanStack Form)
const form = useForm({
	defaultValues: { email: "", password: "" },
});
```

### 4. コード分割

**方針**:

- ルート単位での自動コード分割 (TanStack Start)
- 重いコンポーネントは React.lazy で遅延ロード
- Dynamic Import の活用

**例**:

```typescript
// 遅延ロード
const HeavyChart = lazy(() => import("@/widgets/heavy-chart"));

function Dashboard() {
	return (
		<Suspense fallback={<ChartSkeleton />}>
			<HeavyChart data={data} />
		</Suspense>
	);
}
```

### 5. エラーハンドリング

**方針**:

- Error Boundary でエラーをキャッチ
- React Query のエラーハンドリング
- ユーザーフレンドリーなエラーメッセージ

**例**:

```typescript
// app/providers/error-boundary.tsx
<ErrorBoundary fallback={<ErrorFallback />}>
	<App />
</ErrorBoundary>;

// features/transaction/add-transaction/model/use-add-transaction.ts
const mutation = useMutation({
	mutationFn: addTransaction,
	onError: (error) => {
		toast.error("取引の追加に失敗しました");
		console.error(error);
	},
});
```

## データフロー

### 一般的なデータフェッチフロー

```
1. ユーザーアクション (ボタンクリックなど)
   ↓
2. Feature コンポーネント (イベントハンドラー)
   ↓
3. React Query Mutation/Query
   ↓
4. API クライアント (Fetch API)
   ↓
5. バックエンド API
   ↓
6. レスポンスの受信
   ↓
7. React Query キャッシュ更新
   ↓
8. UI の自動再レンダリング
```

### Server Component でのデータフェッチ例

```typescript
// app/routes/dashboard/index.tsx (Server Component)
export default async function DashboardPage() {
	// サーバーサイドでデータフェッチ
	const user = await fetchUser();
	const portfolio = await fetchPortfolio(user.id);

	return (
		<div>
			<PortfolioSummary data={portfolio} />
		</div>
	);
}
```

### Client Component でのデータフェッチ例

```typescript
// widgets/portfolio-summary/model/use-portfolio-data.ts
export function usePortfolioData(userId: string) {
	return useQuery({
		queryKey: ["portfolio", userId],
		queryFn: () => fetchPortfolio(userId),
		staleTime: 1000 * 60 * 5, // 5分間キャッシュ
	});
}

// widgets/portfolio-summary/ui/portfolio-summary.tsx
("use client");

export function PortfolioSummary({ userId }: Props) {
	const { data, isLoading, error } = usePortfolioData(userId);

	if (isLoading) return <Skeleton />;
	if (error) return <ErrorMessage error={error} />;

	return <div>{/* UIレンダリング */}</div>;
}
```

## パフォーマンス最適化

### 1. React Query によるキャッシング

- 自動的なキャッシング
- staleTime と cacheTime の適切な設定
- Prefetching の活用

### 2. コンポーネントのメモ化

```typescript
// 重い計算のメモ化
const expensiveValue = useMemo(() => {
	return calculateExpensiveValue(data);
}, [data]);

// コールバック関数のメモ化
const handleClick = useCallback(() => {
	doSomething(value);
}, [value]);
```

### 3. 画像最適化

- Next.js Image コンポーネント (または類似の最適化)
- 遅延ロード
- WebP フォーマットの使用

### 4. バンドルサイズの最適化

- Tree Shaking
- 動的インポート
- ライブラリの選定 (バンドルサイズを考慮)

## アクセシビリティ

### 方針

- WCAG 2.1 AA レベルの準拠を目指す
- shadcn/ui による Radix UI のアクセシビリティ機能を活用
- キーボードナビゲーション対応
- スクリーンリーダー対応

### 実装例

```typescript
// アクセシブルなボタン
<Button
  aria-label="ポートフォリオを追加"
  onClick={handleAdd}
>
  <PlusIcon aria-hidden="true" />
</Button>

// フォームのラベリング
<Label htmlFor="email">メールアドレス</Label>
<Input
  id="email"
  type="email"
  aria-required="true"
  aria-invalid={!!errors.email}
/>
```

## テスト方針

### テス �� 戦略

1. **コンポーネントテスト** (主軸)

   - React Testing Library を使用
   - ユーザーの操作に近い形でテスト
   - 主要な Feature と Widget をカバー

2. **E2E テスト** (クリティカルパス)

   - Playwright または Cypress を使用
   - 主要なユーザーフローをカバー

3. **ユニットテスト** (必要に応じて)
   - ユーティリティ関数
   - 複雑なビジネスロジック

### テスト例

```typescript
// features/auth/login/ui/login-form.test.tsx
import { render, screen, userEvent } from "@testing-library/react";
import { LoginForm } from "./login-form";

describe("LoginForm", () => {
	it("ログインフォームが正しく動作する", async () => {
		const user = userEvent.setup();
		const onSuccess = vi.fn();

		render(<LoginForm onSuccess={onSuccess} />);

		// フォーム入力
		await user.type(
			screen.getByLabelText("メールアドレス"),
			"test@example.com"
		);
		await user.type(screen.getByLabelText("パスワード"), "password123");

		// 送信
		await user.click(screen.getByRole("button", { name: "ログイン" }));

		// 検証
		expect(onSuccess).toHaveBeenCalledWith(
			expect.objectContaining({
				email: "test@example.com",
			})
		);
	});
});
```

## セキュリティ考慮事項

- XSS 対策 (React のデフォルトエスケープ)
- CSRF トークンの管理
- 認証トークンの安全な保存 (httpOnly Cookie)
- 環境変数の適切な管理

## 開発環境

### 推奨エディタ設定

```json
// .vscode/settings.json
{
	"editor.formatOnSave": true,
	"editor.defaultFormatter": "esbenp.prettier-vscode",
	"editor.codeActionsOnSave": {
		"source.fixAll.eslint": true
	}
}
```

### 必須 VSCode 拡張

- ESLint
- Prettier
- Tailwind CSS IntelliSense
- TypeScript and JavaScript Language Features

## 関連ドキュメント

- [UI 仕様](../../03-ui-spec/)
- [API 仕様](../../05-api-reference/)
