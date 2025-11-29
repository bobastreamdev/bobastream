CREATE TABLE video_views (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    viewer_ip VARCHAR(45) NOT NULL,
    user_agent TEXT,
    watch_duration_seconds INTEGER DEFAULT 0,
    watched_percentage DECIMAL(5,2) DEFAULT 0,
    session_id VARCHAR(255) NOT NULL,
    viewed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    date DATE DEFAULT CURRENT_DATE,
    UNIQUE(session_id, video_id)
);

-- Indexes
CREATE INDEX idx_video_views_video_id ON video_views(video_id);
CREATE INDEX idx_video_views_user_id ON video_views(user_id);
CREATE INDEX idx_video_views_date ON video_views(date, video_id);
CREATE INDEX idx_video_views_session ON video_views(session_id);
CREATE INDEX idx_video_views_watched_pct ON video_views(watched_percentage);