# Backlog情報検索アプリケーション

### アプリケーションの実行

1. リポジトリをクローン
2. プロジェクトディレクトリに移動
3. 以下のコマンドでアプリケーションを起動：

ローカル開発の場合：

　①バックエンドの起動：

```bash
cd backend
go mod download
make run または　go run cmd/server/main.go
```

　②フロントエンドの起動：
```bash
cd frontend
npm install
npm start
```
Dockerを開始するコマンド：
```bash
docker-compose up --build
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
- `OPENAI_API_KEY`: OpenAI APIキー（AI分析機能に必要）

これらの変数は `backend/.env` ファイルで設定できます。

### OpenAI API設定

バックログアイテムのAI分析機能を使用するには、以下の手順でOpenAI APIキーを設定してください：
1. [OpenAIのウェブサイト](https://platform.openai.com/)でアカウントを作成し、APIキーを取得
2. `.env`ファイルの`OPENAI_API_KEY`変数にAPIキーを設定：
   ```
   OPENAI_API_KEY=your_openai_api_key_here
   ```
3. バックエンドサーバーを再起動

注意: OpenAI APIキーが設定されていない場合、AI分析機能はモックデータを使用します。

docker-composeでDockerを停止するコマンド
```bash
# コンテナを停止するだけ（コンテナは残ります）
docker-compose stop

# コンテナを停止して削除する（次回起動を早くしたい場合はこちら）
docker-compose down

# 全てのリソース（ボリュームを含む）を削除する場合
docker-compose down --volumes
```


## アーキテクチャ

- フロントエンドはReactで構築されます
- バックエンドはGoで構築され、Ginフレームワークを使用しています
- 認証はBacklogのOAuthを使用して処理されます
- コンテナ化とデプロイにはDockerを使用しています

## 技術スタック

### バックエンド
- Go 1.20
- GraphQL (gqlgen)
- クリーンアーキテクチャ
- OAuth 2.0認証

バックエンドのコントローラのロジックを直接main.goに記述し、
リポジトリの実装はmemoryパッケージに直接実装することで、シンプルに保っているようです。
将来的にデータ永続化のため、DBの実装および機能追加により、これらの空のディレクトリを活用して、
適切に分離された構造に移行する。

### インフラ
- Docker
- Docker Compose

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


## OAuth 2.0認証フロー

OAuth 2.0の認可コードフローを使用してBacklog APIにアクセスしています。

1. アプリケーションがBacklogの認可エンドポイントにユーザーをリダイレクト
2. ユーザーがBacklogでアプリケーションのアクセスを許可
3. Backlogが認可コードをリダイレクトURIに返す
4. アプリケーションがこの認可コードを使用してアクセストークンを取得
5. アクセストークンを使用してBacklog APIにアクセス

## 機能

- Backlog更新情報の検索
- キーワードによるフィルタリング
- お気に入り登録と解除機能　
　✳インメモ保存なので、BEサーバーを再起動するとお気に入り登録項目が失われます。
- OAuth 2.0によるBacklog認証
- OpenAI APIを使用したアイテムのAI分析機能 