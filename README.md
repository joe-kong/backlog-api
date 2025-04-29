# Backlog情報検索アプリケーション

### アプリケーションの実行

1. リポジトリをクローン
2. プロジェクトディレクトリに移動
3. 環境変数は適切に設定
4. 以下のコマンドでアプリケーションを起動：

ローカル開発の場合：

　①バックエンドの起動：

```bash
cd backend
go mod download
make run または　go run cmd/server/main.go
```

ユニットテスト実行コマンド
```
make test
```

　②フロントエンドの起動：
```bash
cd frontend
npm install
npm start
```

4. http://localhost:3000 でアプリケーションにアクセス



### 環境変数

アプリケーションは設定のために以下の環境変数を使用します：

- `BACKLOG_SPACE_URL`: BacklogスペースのURL
- `BACKLOG_CLIENT_ID`: BacklogのOAuthクライアントID
- `BACKLOG_CLIENT_SECRET`: BacklogのOAuthクライアントシークレット
- `BACKLOG_AUTH_URL`: Backlog認証URL
- `BACKLOG_TOKEN_URL`: BacklogトークンURL
- `OAUTH_REDIRECT_URI`: OAuthリダイレクトURI（デフォルト: http://localhost:8081/api/auth/callback）
- `PORT`: バックエンドサーバーのポート（デフォルト: 8081）
- `FRONTEND_URL`: フロントエンドアプリケーションのURL（デフォルト: http://localhost:3000）
- `REACT_APP_API_URL`: バックエンドAPIのURL（デフォルト: http://localhost:8081）- フロントエンド用
- `OPENAI_API_KEY`: OpenAI APIキー（AI分析機能に必要）

#### 環境変数の設定方法

**バックエンド**:
環境ごとに設定ファイルを分けていません、デプロイ先より必要な変数が値を変更してビルド/デプロイしてください。
ファイルパス：`backend/.env`
ローカルテストの場合の設定例：
```
PORT=8081
FRONTEND_URL=http://localhost:3000
BACKLOG_SPACE_URL=your_backlog_space_url
BACKLOG_CLIENT_ID=your_client_id
BACKLOG_CLIENT_SECRET=your_client_secret
OAUTH_REDIRECT_URI=http://localhost:8081/api/auth/callback　
BACKLOG_AUTH_URL=your_auth_url
BACKLOG_TOKEN_URL=your_token_url
USE_DYNAMODB=false　# ローカルテストでも永続化したい場合、trueにしてください。
DYNAMODB_REGION=ap-northeast-1
OPENAI_API_KEY=your_openai_api_key
```


**フロントエンド**:
フロントエンド用の環境変数は以下のファイルに設定します：
- 開発環境: `frontend/.env.development`
- 本番環境: `frontend/.env.production`

例:
```
REACT_APP_API_URL=http://localhost:8081  # 開発環境
REACT_APP_API_URL=https://api.example.com  # 本番環境
```

### OpenAI API設定

バックログアイテムのAI分析機能を使用するには、OpenAI APIキーを設定してください：
1. [OpenAIのウェブサイト](https://platform.openai.com/)でアカウントを作成し、APIキーを取得
2. `.env`ファイルの`OPENAI_API_KEY`変数にAPIキーを設定：
   ```
   OPENAI_API_KEY=your_openai_api_key
   ```
⇨OpenAI APIキーが設定されていない場合、AI分析機能はモックデータを使用します。

## アーキテクチャ

- フロントエンドはReactで構築されます
- バックエンドはGoで構築され、Ginフレームワークとクリーンアーキテクチャを使用しています
- 認証はBacklogのOAuthを使用して処理されます
- Dockerを使用しています

## 技術スタック

### バックエンド
- Go 1.20
- GraphQL (gqlgen)
- クリーンアーキテクチャ
- OAuth 2.0認証
- DynamoDB（お気に入り永続化）

バックエンドのコントローラのロジックを直接main.goに記述し、
リポジトリの実装はmemoryパッケージに直接実装することで、シンプルに保っているようです。
将来的に業務追加により、クリーンアーキテクチャをしっかりする場合、空のディレクトリを活用して、
適切に分離された構造に移行する。

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
│           └── persistence/# 永続化層
│               ├── memory/ # インメモリ実装
│               └── dynamodb/# DynamoDB実装
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


## 実装機能

- Backlog更新情報の検索
- キーワードによるフィルタリング
- お気に入り登録と解除機能　
　✳インメモ保存モードの場合、BEサーバーを再起動するとお気に入り登録項目が失われます。DynamoDBモードを有効にすることで永続化できます。
- OAuth 2.0によるBacklog認証
- OpenAI APIを使用したアイテムのAI分析機能 


## インフラ構成

以下のAWSサービスを使用してクラウドにリリースしています：

### フロントエンド（React SPA）

- **ホスティング**: Amazon S3 + CloudFront
- **ドメイン**: backlog.51incloud.com

### バックエンド（Go API）

- **ホスティング**: ECS Fargate
- **ロードバランサー**: Application Load Balancer (HTTPS:443)
- **ドメイン**: backlog-app-api.51incloud.com
- **コンテナ化**: ECRにDockerイメージを保存

### ネットワーク

- **VPC**: 10.0.0.0/16
- **サブネット**:
  - パブリック: 2つのAZ (ap-northeast-1a, ap-northeast-1c)
  - プライベート: 2つのAZ (ap-northeast-1a, ap-northeast-1c)
- **インターネット接続**:
  - パブリックサブネット: Internet Gateway
  - プライベートサブネット: NAT Gateway

### データストア

- **お気に入り情報**: DynamoDB

### セキュリティ

- **CloudFront**: TLS接続 (HTTPS)
- **ALB**: TLS接続 (HTTPS), セキュリティグループによるアクセス制限
- **ECS**: プライベートサブネットでの実行、ALBからのみアクセス可能
- **IAM**: 最小権限の原則に基づくロール設定

### CI/CD

- デプロイスクリプト:
  - `deploy-frontend.sh`: フロントエンドのビルドとS3へのデプロイ
  - `deploy-backend.sh`: バックエンドのDockerビルドとECRへのプッシュ、ECS更新

## リソースサイズと最適化

- **ECS Fargate**: 最小構成で運用
- **S3+CloudFront**: 静的コンテンツのキャッシュによる高速化
- **DynamoDB**: オンデマンドキャパシティモードで自動スケーリング 