variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "ap-northeast-1"
}

# 変数
variable "app_name" {
  description = "アプリケーション名"
  default     = "backlog-app"
}


variable "environment" {
  description = "Environment (dev, staging, prod)"
  type        = string
  default     = "production"
}
