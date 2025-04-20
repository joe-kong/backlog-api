# Backlog 更新情報検索アプリケーション

このアプリケーションは、Backlog APIを利用して更新情報を検索し、お気に入り登録できるWebアプリケーションです。

## 技術スタック

### バックエンド
- Go 1.20
- GraphQL (gqlgen)
- クリーンアーキテクチャ
- OAuth 2.0認証

### フロントエンド
- React
- TypeScript
- Apollo Client

## プロジェクト構造

```
.
├── backend/                # Goバックエンド
│   ├── cmd/                # エントリーポイント
│   │   └── server/         # サーバー実行コード
│   └── internal/           # 内部パッケージ
│       ├── domain/         # ドメインモデル
│       ├── usecase/        # ユースケース
│       ├── interface/      # インターフェース
│       └── infrastructure/ # インフラストラクチャ
└── frontend/              # Reactフロントエンド
    ├── public/            # 静的ファイル
    └── src/               # ソースコード
        ├── components/    # UIコンポーネント
        ├── containers/    # コンテナコンポーネント
        ├── pages/         # ページコンポーネント
        ├── graphql/       # GraphQLクエリ・ミューテーション
        ├── hooks/         # カスタムフック
        ├── utils/         # ユーティリティ関数
        ├── context/       # Reactコンテキスト
        └── types/         # 型定義
```

## 開発における簡略化ポイント

- ユニットテスト: 一部のコアロジックのみに限定
- エラーハンドリング: 基本的なものに限定
- UIデザイン: 最小限の実装に限定
- ドキュメンテーション: 主要なコンポーネントとAPIのみ

## セットアップ

### 前提条件
- Go 1.20以上
- Node.js 16以上
- npm 7以上

### バックエンド開発

```bash
cd backend
go mod download
go run cmd/server/main.go
```

### BEのプロセスを終了して再起動する場合
```bash
pkill -f "go run cmd/server/main.go"
```

### フロントエンド開発

```bash
cd frontend
npm install
npm start
```

### FEのプロセスを終了するコマンド
```bash
pkill -f "react-scripts start" || true
```

## OAuth 2.0認証フロー

このアプリケーションはOAuth 2.0の認可コードフローを使用してBacklog APIにアクセスします。

1. アプリケーションがBacklogの認可エンドポイントにユーザーをリダイレクト
2. ユーザーがBacklogでアプリケーションのアクセスを許可
3. Backlogが認可コードをリダイレクトURIに返す
4. アプリケーションがこの認可コードを使用してアクセストークンを取得
5. アクセストークンを使用してBacklog APIにアクセス

## 機能

- Backlog更新情報の検索
- キーワードによるフィルタリング
- お気に入り登録機能
- OAuth 2.0によるBacklog認証 