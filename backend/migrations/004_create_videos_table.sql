CREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(500) NOT NULL,
    description TEXT,
    thumbnail_url TEXT,
    source_url TEXT NOT NULL,
    source_url_expires_at TIMESTAMP,
    duration_seconds INTEGER,
    file_size_mb DECIMAL(10,2),
    pcloud_file_id VARCHAR(255),
    pcloud_credential_id UUID NOT NULL REFERENCES pcloud_credentials(id) ON DELETE RESTRICT,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    tags TEXT[],
    view_count INTEGER DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    is_published BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_videos_published ON videos(is_published, published_at DESC);
CREATE INDEX idx_videos_category ON videos(category_id);
CREATE INDEX idx_videos_pcloud_credential ON videos(pcloud_credential_id);
CREATE INDEX idx_videos_view_count ON videos(view_count DESC);
CREATE INDEX idx_videos_like_count ON videos(like_count DESC);
CREATE INDEX idx_videos_tags ON videos USING GIN(tags);
CREATE INDEX idx_videos_deleted_at ON videos(deleted_at);

-- Comment for clarity
COMMENT ON COLUMN videos.pcloud_credential_id IS 'References pcloud account used to store this video';
COMMENT ON COLUMN videos.category_id IS 'Video category for filtering and organization';

-- Trigger for updated_at
CREATE TRIGGER update_videos_updated_at
    BEFORE UPDATE ON videos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();