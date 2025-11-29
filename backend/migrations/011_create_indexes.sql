-- Composite index for feed query optimization
CREATE INDEX idx_videos_feed_score ON videos(is_published, published_at DESC, view_count DESC, like_count DESC) 
WHERE is_published = true;

-- Composite index for video views aggregation
CREATE INDEX idx_video_views_aggregate ON video_views(video_id, watched_percentage, session_id);

-- Composite index for analytics queries
CREATE INDEX idx_ad_impressions_analytics ON ad_impressions(date, ad_id, impression_type);

-- Composite index for videos by pCloud account (useful for admin dashboard)
CREATE INDEX idx_videos_by_pcloud_account ON videos(pcloud_credential_id, created_at DESC);