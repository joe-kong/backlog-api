#!/bin/bash
set -e

# 設定変数
S3_BUCKET="your-s3-bucket-name"
CLOUDFRONT_DISTRIBUTION_ID="your-cloudfront-distribution-id"
REGION="ap-northeast-1"

# ビルド
echo "Reactアプリケーションをビルドしています..."
cd frontend
npm ci
npm run build

# S3にアップロード
echo "S3バケットにビルド成果物をアップロードしています..."
aws s3 sync build/ s3://$S3_BUCKET/ --delete --region $REGION

# CloudFrontキャッシュを無効化
echo "CloudFrontキャッシュを無効化しています..."
aws cloudfront create-invalidation \
  --distribution-id $CLOUDFRONT_DISTRIBUTION_ID \
  --paths "/*" \
  --region $REGION

echo "フロントエンドのデプロイが完了しました!" 