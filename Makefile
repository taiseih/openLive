.PHONY: help up down logs infra server frontend all clean

help: ## ヘルプを表示
	@echo "openLive モノレポ 分離版Docker Compose"
	@echo ""
	@echo "使用可能なコマンド:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "起動順序:"
	@echo "  1. make infra     - 共通インフラを起動"
	@echo "  2. make server    - Serverを起動"
	@echo "  3. make frontend  - Frontendを起動"
	@echo "または:"
	@echo "  make all          - 全サービスを一括起動"

# 共通インフラのみ起動（PostgreSQL）
infra: ## 共通インフラのみ起動
	@echo "Starting infrastructure services (PostgreSQL)..."
	docker-compose up -d

# Serverのみ起動
server: ## Serverのみ起動（要: infra起動済み）
	@echo "Starting Server..."
	docker-compose -f docker-compose.server.yml up -d

# Frontendのみ起動
frontend: ## Frontendのみ起動（要: infra, server起動済み）
	@echo "Starting Frontend..."
	docker-compose -f docker-compose.frontend.yml up -d

# 全サービスを一括起動
all: ## 全サービスを一括起動
	@echo "Starting all services..."
	docker-compose up -d
	docker-compose -f docker-compose.server.yml up -d
	docker-compose -f docker-compose.frontend.yml up -d
	@echo ""
	@echo "✅ All services started!"
	@echo "  - Frontend: http://localhost:3000"
	@echo "  - Backend API: http://localhost:8080"

# 全サービスを停止
down: ## 全サービスを停止
	@echo "Stopping all services..."
	docker-compose -f docker-compose.frontend.yml down
	docker-compose -f docker-compose.server.yml down
	docker-compose down

# インフラのみ停止
infra-down: ## 共通インフラのみ停止
	docker-compose down

# Serverのみ停止
server-down: ## Serverのみ停止
	docker-compose -f docker-compose.server.yml down

# Frontendのみ停止
frontend-down: ## Frontendのみ停止
	docker-compose -f docker-compose.frontend.yml down

# ログを表示
logs: ## 全サービスのログを表示
	docker-compose logs -f

# インフラログを表示
logs-infra: ## 共通インフラのログを表示
	docker-compose logs -f

# Serverログを表示
logs-server: ## Serverのログを表示
	docker-compose -f docker-compose.server.yml logs -f

# Frontendログを表示
logs-frontend: ## Frontendのログを表示
	docker-compose -f docker-compose.frontend.yml logs -f

# Serverを再起動
restart-server: ## Serverを再起動
	docker-compose -f docker-compose.server.yml restart

# Frontendを再起動
restart-frontend: ## Frontendを再起動
	docker-compose -f docker-compose.frontend.yml restart

# 全サービスを再起動
restart: ## 全サービスを再起動
	@echo "Restarting all services..."
	docker-compose restart
	docker-compose -f docker-compose.server.yml restart
	docker-compose -f docker-compose.frontend.yml restart

# ボリューム含めて完全削除
clean: ## 全コンテナとボリュームを削除
	@echo "Cleaning up all containers and volumes..."
	docker-compose -f docker-compose.frontend.yml down -v
	docker-compose -f docker-compose.server.yml down -v
	docker-compose down -v
	@echo "✅ Cleanup complete!"

# リビルド
rebuild: ## 全サービスをリビルドして起動
	@echo "Rebuilding all services..."
	docker-compose up -d --build
	docker-compose -f docker-compose.server.yml up -d --build
	docker-compose -f docker-compose.frontend.yml up -d --build

# Serverのみリビルド
rebuild-server: ## Serverのみリビルドして起動
	docker-compose -f docker-compose.server.yml up -d --build

# Frontendのみリビルド
rebuild-frontend: ## Frontendのみリビルドして起動
	docker-compose -f docker-compose.frontend.yml up -d --build

# ステータス確認
status: ## 全サービスのステータスを確認
	@echo "=== Infrastructure Status ==="
	docker-compose ps
	@echo ""
	@echo "=== Server Status ==="
	docker-compose -f docker-compose.server.yml ps
	@echo ""
	@echo "=== Frontend Status ==="
	docker-compose -f docker-compose.frontend.yml ps

# DBマイグレーション
migrate: ## DBマイグレーションを実行
	@echo "Running database migrations..."
	docker-compose exec postgres psql -U openlive -d openlive_db -f /docker-entrypoint-initdb.d/001_init.sql
	@echo "✅ Migration complete!"

# DB接続
db: ## PostgreSQLに接続
	docker-compose exec postgres psql -U openlive -d openlive_db

# 開発用：Serverシェル
shell-server: ## Serverコンテナのシェルに接続
	docker-compose -f docker-compose.server.yml exec server sh

# 開発用：Frontendシェル
shell-frontend: ## Frontendコンテナのシェルに接続
	docker-compose -f docker-compose.frontend.yml exec frontend sh

