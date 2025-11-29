.PHONY: help run build clean test migrate-up migrate-down seed dev

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	@echo "🚀 Starting BOBA STREAM backend..."
	@go run cmd/server/main.go

build: ## Build the application
	@echo "🔨 Building application..."
	@go build -o bin/bobastream cmd/server/main.go
	@echo "✅ Build complete: bin/bobastream"

clean: ## Clean build files
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@go clean
	@echo "✅ Clean complete"

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v ./...

dev: ## Run in development mode with hot reload (requires air)
	@echo "🔥 Starting development server with hot reload..."
	@air

tidy: ## Tidy go modules
	@echo "📦 Tidying modules..."
	@go mod tidy
	@echo "✅ Modules tidied"

migrate-up: ## Run database migrations
	@echo "⬆️  Running migrations..."
	@psql $(DB_URL) -f migrations/001_create_users_table.sql
	@psql $(DB_URL) -f migrations/002_create_pcloud_credentials_table.sql
	@psql $(DB_URL) -f migrations/002_5_create_categories_table.sql
	@psql $(DB_URL) -f migrations/003_create_videos_table.sql
	@psql $(DB_URL) -f migrations/004_create_wrapper_links_table.sql
	@psql $(DB_URL) -f migrations/005_create_video_views_table.sql
	@psql $(DB_URL) -f migrations/006_create_video_likes_table.sql
	@psql $(DB_URL) -f migrations/007_create_ads_table.sql
	@psql $(DB_URL) -f migrations/008_create_ad_impressions_table.sql
	@psql $(DB_URL) -f migrations/009_create_daily_stats_table.sql
	@psql $(DB_URL) -f migrations/010_create_indexes.sql
	@echo "✅ Migrations complete"

migrate-down: ## Rollback database migrations
	@echo "⬇️  Rolling back migrations..."
	@psql $(DB_URL) -c "DROP TABLE IF EXISTS ad_impressions CASCADE;"
	@psql $(DB_URL) -c "DROP TABLE IF EXISTS ads CASCADE;"
	@psql $(DB_URL) -c "DROP TABLE IF EXISTS video_likes CASCADE;"
	@psql $(DB_URL) -c "DROP TABLE IF EXISTS video_views CASCADE;"
	@psql $(DB_URL) -c "DROP TABLE IF EXISTS wrapper_links CASCADE;"
	@psql $(DB_URL) -c "DROP TABLE IF EXISTS videos CASCADE;"
	@psql $(DB_URL) -c "DROP TABLE IF EXISTS categories CASCADE;"
	@psql $(DB_URL) -c "DROP TABLE IF EXISTS pcloud_credentials CASCADE;"
	@psql $(DB_URL) -c "DROP TABLE IF EXISTS users CASCADE;"
	@psql $(DB_URL) -c "DROP TABLE IF EXISTS daily_stats CASCADE;"
	@psql $(DB_URL) -c "DROP TYPE IF EXISTS user_role CASCADE;"
	@psql $(DB_URL) -c "DROP TYPE IF EXISTS ad_type CASCADE;"
	@psql $(DB_URL) -c "DROP TYPE IF EXISTS impression_type CASCADE;"
	@echo "✅ Rollback complete"

seed: ## Seed database with sample data
	@echo "🌱 Seeding database..."
	@go run scripts/seed.go
	@echo "✅ Seeding complete"

install-tools: ## Install development tools
	@echo "🔧 Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@echo "✅ Tools installed"

docker-up: ## Start Docker containers
	@echo "🐳 Starting Docker containers..."
	@docker-compose up -d
	@echo "✅ Docker containers started"

docker-down: ## Stop Docker containers
	@echo "🐳 Stopping Docker containers..."
	@docker-compose down
	@echo "✅ Docker containers stopped"

docker-logs: ## Show Docker logs
	@docker-compose logs -f

# Database connection string (override with your own)
DB_URL ?= postgresql://postgres:bobastream123@localhost:5555/boba_stream?sslmode=disable