.PHONY: help build run run-local test clean up down logs migrate

help: ## ヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Goバイナリをビルド
	go build -o bin/server cmd/server/main.go

run: ## Docker環境でサーバーを起動
	docker-compose up

run-local: ## ローカル環境でサーバーを起動
	go run cmd/server/main.go

test: ## テストを実行
	go test -v ./...

test-coverage: ## テストカバレッジを取得
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## ビルド成果物を削除
	rm -rf bin/
	rm -rf tmp/
	rm -f coverage.out coverage.html

up: ## Dockerコンテナを起動
	cd .. && docker-compose up -d

down: ## Dockerコンテナを停止
	cd .. && docker-compose down

logs: ## Dockerログを表示
	cd .. && docker-compose logs -f server

migrate: ## マイグレーションを実行
	@echo "Running migrations..."
	PGPASSWORD=openlive psql -h localhost -U openlive -d openlive_db -f migrations/001_init.sql

deps: ## 依存関係をインストール
	go mod download
	go mod tidy

lint: ## リントを実行
	golangci-lint run

fmt: ## コードをフォーマット
	go fmt ./...

