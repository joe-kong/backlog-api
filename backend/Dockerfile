# ビルドステージ
FROM golang:1.22-alpine AS builder

# 必要なパッケージのインストール
RUN apk add --no-cache git

# ワーキングディレクトリの設定
WORKDIR /app

# 依存関係のコピーとダウンロード
# COPYコマンドを分割して、ファイルが存在しない場合のエラーを回避
COPY go.mod* ./
COPY go.sum* ./
RUN go mod download || go mod tidy

# ソースコードのコピー
COPY . .

# アプリケーションのビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# 実行ステージ
FROM alpine:latest

# タイムゾーンの設定
RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone && \
    apk del tzdata

# 作業ディレクトリの作成
WORKDIR /app

# ビルドステージからバイナリをコピー
COPY --from=builder /app/server .

# .envファイルをコピー
COPY .env .

# 実行ユーザーの設定
RUN adduser -D appuser
USER appuser

# ポートの公開
EXPOSE 8081

# アプリケーションの実行
CMD ["./server"] 