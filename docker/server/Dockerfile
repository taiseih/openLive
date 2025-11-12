# ビルドステージ
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 依存関係のキャッシュ
COPY server/go.mod server/go.sum ./
RUN go mod download

# ソースコードをコピー
COPY server/ ./

# バイナリをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server cmd/server/main.go

# 実行ステージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# ビルドしたバイナリをコピー
COPY --from=builder /app/bin/server .

# マイグレーションファイルをコピー
COPY server/migrations ./migrations

EXPOSE 8080

CMD ["./server"]

