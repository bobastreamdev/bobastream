-- Drop old index if exists
DROP INDEX IF EXISTS idx_videos_feed_score;

-- ✅ Composite index yang lebih optimal untuk feed query
CREATE INDEX idx_videos_feed_optimized ON videos(
    is_published, 
    published_at DESC, 
    view_count DESC, 
    like_count DESC,
    id  -- ✅ Include PK untuk covering index
) WHERE is_published = true;

-- ✅ Partial index untuk recent videos (last 30 days)
CREATE INDEX idx_videos_recent_feed ON videos(
    published_at DESC,
    view_count DESC
) WHERE is_published = true AND published_at > CURRENT_DATE - INTERVAL '30 days';

-- Add comment for documentation
COMMENT ON INDEX idx_videos_feed_optimized IS 'Optimized composite index for video feed queries with covering index';
COMMENT ON INDEX idx_videos_recent_feed IS 'Partial index for recent videos (last 30 days) to improve hot data access';