# AWSデプロイガイド - Backlogアプリケーション
BacklogアプリケーションをAWS環境にデプロイする手順を説明します。

## アーキテクチャ概要

- **フロントエンド**: S3 + CloudFront
- **バックエンド**: ECS (Fargate) + ALB
- **データベース**: DynamoDB

## 前提条件

- AWSアカウント
- AWS CLIのインストールと設定
- Terraformのインストール
- Dockerのインストール

## 1. 環境変数の設定

まず、必要な環境変数を設定します。プロジェクトルートに`.env.production`ファイルを作成します。

```bash
# AWS認証情報
AWS_REGION=ap-northeast-1
AWS_PROFILE=your-aws-profile

# アプリケーション設定
APP_NAME=backlog-app
ENVIRONMENT=production

# Backlog設定
BACKLOG_SPACE_URL=https://your-space.backlog.jp
BACKLOG_CLIENT_ID=your-client-id
BACKLOG_CLIENT_SECRET=your-client-secret
BACKLOG_AUTH_URL=https://your-space.backlog.jp/OAuth2/authorize
BACKLOG_TOKEN_URL=https://your-space.backlog.jp/OAuth2/token

# その他の設定
OPENAI_API_KEY=your-openai-api-key
```

## 2. インフラストラクチャのデプロイ

Terraformを使用してAWSリソースをプロビジョニングします。

```bash
# Terraformディレクトリに移動
cd terraform

# Terraformの初期化
terraform init

# 実行計画の確認
terraform plan -var-file=..terraform/terraform.tfvars

# インフラのデプロイ
terraform apply -var-file=..terraform/terraform.tfvars
```

デプロイが完了すると、以下の情報が出力されます：
- CloudFrontドメイン（フロントエンドURL）
- ALBエンドポイント（バックエンドAPI URL）
- ECRリポジトリURL

## 3. フロントエンドのデプロイ

フロントエンドをビルドしてS3にデプロイします。

```bash
# フロントエンド環境変数ファイルの作成
cat > frontend/.env.production << EOF
REACT_APP_API_URL=https://api.your-domain.com
EOF

# デプロイスクリプトを実行
chmod +x deploy-frontend.sh
./deploy-frontend.sh
```

## 4. バックエンドのデプロイ

バックエンドをビルドしてECRにプッシュし、ECSにデプロイします。

```bash
# バックエンド環境変数ファイルの作成
cat > backend/.env.production << EOF
PORT=8081
FRONTEND_URL=https://your-cloudfront-domain.cloudfront.net
APP_ENV=production
USE_DYNAMODB=true
DYNAMODB_REGION=ap-northeast-1
BACKLOG_SPACE_URL=${BACKLOG_SPACE_URL}
BACKLOG_CLIENT_ID=${BACKLOG_CLIENT_ID}
BACKLOG_CLIENT_SECRET=${BACKLOG_CLIENT_SECRET}
OAUTH_REDIRECT_URI=https://api.your-domain.com/api/auth/callback
BACKLOG_AUTH_URL=${BACKLOG_AUTH_URL}
BACKLOG_TOKEN_URL=${BACKLOG_TOKEN_URL}
OPENAI_API_KEY=${OPENAI_API_KEY}
EOF

# デプロイスクリプトを実行
./deploy-backend.sh
```

## 5. カスタムドメインの設定（オプション）

本番環境では、独自ドメインを使用することをお勧めします。

### フロントエンド（CloudFront）
1. Route 53でドメインを設定
2. ACM証明書を作成
3. CloudFrontディストリビューションを更新

### バックエンド（ALB）
1. Route 53でサブドメイン（api.your-domain.com）を設定
2. ACM証明書を作成
3. ALBリスナーをHTTPSに更新

## 6. 確認とテスト

デプロイ完了後、以下の点をテストします：

1. フロントエンドにアクセスしてページが表示されるか確認
2. ログイン機能が正常に動作するか確認
3. API呼び出しが成功するか確認
4. DynamoDBへのお気に入り保存が機能するか確認







## トラブルシューティング

### CloudFrontの更新が反映されない
キャッシュ無効化を実行してください：
```bash
aws cloudfront create-invalidation --distribution-id YOUR_DISTRIBUTION_ID --paths "/*"
```

### ECSタスクが失敗する
CloudWatchログを確認して原因を特定してください：
```bash
aws logs get-log-events --log-group-name /ecs/backlog-app --log-stream-name YOUR_LOG_STREAM
```

### CORS問題
バックエンドのCORS設定で、CloudFrontドメインからのリクエストを許可してください。

## 運用とメンテナンス

### ログの確認
```bash
# ECSのログを確認
aws logs get-log-events --log-group-name /ecs/backlog-app --log-stream-name YOUR_LOG_STREAM

# ALBのアクセスログを確認（S3バケットに保存された場合）
aws s3 ls s3://your-alb-logs-bucket/
```

### スケーリング
ECSサービスのタスク数を調整することで、バックエンドのスケーリングが可能です：
```bash
aws ecs update-service --cluster backlog-app-cluster --service backlog-app-service --desired-count 4
```

### データバックアップ
DynamoDBのバックアップを定期的に取得することをお勧めします：
```bash
aws dynamodb create-backup --table-name favorites --backup-name favorites-backup
```

---