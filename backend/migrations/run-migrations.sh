#!/bin/bash
set -e

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing migrations"

# Run migrations in order
echo "Running migration 001: Create users table..."
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f /migrations/001_create_users_table.sql

echo "Running migration 002: Create pcloud credentials table..."
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f /migrations/002_create_pcloud_credentials_table.sql

echo "Running migration 002.5: Create categories table..."
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f /migrations/002_5_create_categories_table.sql

echo "Running migration 004: Create videos table..."
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f /migrations/004_create_videos_table.sql

echo "Running migration 005: Create wrapper links table..."
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f /migrations/005_create_wrapper_links_table.sql

echo "Running migration 006: Create video views table..."
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f /migrations/006_create_video_views_table.sql

echo "Running migration 007: Create video likes table..."
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f /migrations/007_create_video_likes_table.sql

echo "Running migration 008: Create ads table..."
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f /migrations/008_create_ads_table.sql

echo "Running migration 009: Create ad impressions table..."
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f /migrations/009_create_ad_impressions_table.sql

echo "Running migration 010: Create daily stats table..."
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f /migrations/010_create_daily_stats_table.sql

echo "Running migration 011: Create indexes..."
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f /migrations/011_create_indexes.sql

echo "✅ All migrations completed successfully!"