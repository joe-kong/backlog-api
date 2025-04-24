#!/bin/bash
set -e

# 設定変数
ECR_REPOSITORY="your-ecr-repo"
ECS_CLUSTER="your-ecs-cluster"
ECS_SERVICE="your-ecs-service"
AWS_REGION="ap-northeast-1"
IMAGE_TAG=$(date +%Y%m%d-%H%M%S)

# AWS ECR認証
echo "ECRログイン中..."
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR_REPOSITORY

# イメージのビルドとプッシュ
echo "Dockerイメージをビルドしています..."
cd backend
docker build -t $ECR_REPOSITORY:$IMAGE_TAG -t $ECR_REPOSITORY:latest .

echo "ECRにイメージをプッシュしています..."
docker push $ECR_REPOSITORY:$IMAGE_TAG
docker push $ECR_REPOSITORY:latest

# ECSサービスの更新
echo "ECSサービスを更新しています..."
aws ecs update-service \
  --cluster $ECS_CLUSTER \
  --service $ECS_SERVICE \
  --force-new-deployment \
  --region $AWS_REGION

echo "デプロイ開始！サービスが安定するまで待機してください..."
aws ecs wait services-stable \
  --cluster $ECS_CLUSTER \
  --services $ECS_SERVICE \
  --region $AWS_REGION

echo "バックエンドのデプロイが完了しました!" 