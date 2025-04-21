# ビルドステージ
FROM node:16-alpine AS builder

# 作業ディレクトリの設定
WORKDIR /app

# 依存関係ファイルのコピー
COPY package.json package-lock.json ./

# 依存関係のインストール
RUN npm ci

# ソースコードのコピー
COPY . .

# アプリケーションのビルド
RUN npm run build

# 開発環境での実行ステージ
FROM node:16-alpine AS dev

# 作業ディレクトリの設定
WORKDIR /app

# 依存関係ファイルのコピー
COPY package.json package-lock.json ./

# 開発用依存関係のインストール
RUN npm install

# ソースコードのコピー
COPY . .

# ポートの公開
EXPOSE 3000

# アプリケーションの実行
CMD ["npm", "start"] 