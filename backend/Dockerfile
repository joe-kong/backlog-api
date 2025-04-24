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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server/main.go

# 実行ステージ
FROM alpine:latest

# タイムゾーンとSSL証明書の設定
RUN apk --no-cache add tzdata ca-certificates && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone

# 作業ディレクトリの作成
WORKDIR /app

# ビルドステージからバイナリをコピー
COPY --from=builder /app/server .
COPY --from=builder /app/.env* ./

# 権限の設定
RUN chmod +x ./server && \
    adduser -D appuser && \
    chown -R appuser:appuser /app

# 実行ユーザーの設定
USER appuser

# ポートの公開
EXPOSE 8081

# アプリケーションの実行
CMD ["./server"] 