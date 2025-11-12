# ビルドステージ
FROM node:18-alpine AS builder

WORKDIR /app

# 依存関係をコピーしてインストール
COPY frontend/package*.json ./
RUN npm ci --only=production

# ソースコードをコピー
COPY frontend/ ./

# Next.jsアプリをビルド
RUN npm run build

# 実行ステージ
FROM node:18-alpine

WORKDIR /app

ENV NODE_ENV production

# 必要なファイルをコピー
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

EXPOSE 3000

ENV PORT 3000

CMD ["node", "server.js"]

