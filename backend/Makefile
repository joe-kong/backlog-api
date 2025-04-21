# バックエンドサービスのMakefile

# 変数定義
BINARY_NAME=server
GO_FILES=$(shell find . -name "*.go" -type f)
BUILD_DIR=build
DOCKER_IMAGE=backlog-backend

.PHONY: all build run clean test docker-build docker-run lint help init deps

# デフォルトターゲット
all: build

# ヘルプコマンド
help:
	@echo "利用可能なコマンド:"
	@echo "make build         - バックエンドをビルド"
	@echo "make run           - ローカルでバックエンドを実行"
	@echo "make clean         - ビルド成果物を削除"
	@echo "make test          - テストを実行"
	@echo "make lint          - コードの静的解析を実行"
	@echo "make docker-build  - Dockerイメージをビルド"
	@echo "make docker-run    - Dockerコンテナを実行"
	@echo "make init          - Goモジュールの初期化"
	@echo "make deps          - 依存関係の更新"
	@echo "make dev           - 開発モードで実行（ホットリロード）"
	@echo "make all           - ビルドを実行（デフォルト）"

# Goモジュールの初期化
init:
	@echo "Goモジュールを初期化中..."
	@if [ ! -f go.mod ]; then \
		go mod init github.com/yourname/backlog; \
	fi
	@make deps

# 依存関係をインストール
deps:
	@echo "依存関係を更新中..."
	@go mod tidy
	@go mod download

# バックエンドをビルド
build: deps
	@echo "バックエンドをビルド中..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

# ローカルでバックエンドを実行
run: deps
	@echo "バックエンドを実行中..."
	@go run cmd/server/main.go

# ビルド成果物を削除
clean:
	@echo "ビルド成果物を削除中..."
	@rm -rf $(BUILD_DIR)
	@rm -rf tmp

# テストを実行
test: deps
	@echo "テストを実行中..."
	@go test -v ./...

# コードの静的解析を実行
lint: deps
	@echo "コードの静的解析を実行中..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lintが見つかりません。インストールするには：go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

# Dockerイメージをビルド
docker-build: init
	@echo "go.modとgo.sumの存在を確認..."
	@if [ ! -f go.mod ] || [ ! -f go.sum ]; then \
		echo "go.modまたはgo.sumが見つかりません。make initを実行してください"; \
		exit 1; \
	fi
	@echo "Dockerイメージをビルド中..."
	@docker build -t $(DOCKER_IMAGE) -f Dockerfile ..

# Dockerコンテナを実行
docker-run:
	@echo "Dockerコンテナを実行中..."
	@docker run -p 8081:8081 --env-file .env $(DOCKER_IMAGE)

# 開発用のホットリロード実行（air必須）
dev: deps
	@echo "開発モードで実行中（ホットリロード）..."
	@if command -v air > /dev/null; then \
		air -c .air.toml; \
	else \
		echo "airが見つかりません。インストールするには：go install github.com/cosmtrek/air@latest"; \
		go run cmd/server/main.go; \
	fi 