CREATE TABLE wrapper_links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    wrapper_token VARCHAR(255) UNIQUE NOT NULL,
    slug VARCHAR(500) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_wrapper_token ON wrapper_links(wrapper_token);
CREATE INDEX idx_wrapper_slug ON wrapper_links(slug);
CREATE INDEX idx_wrapper_video_id ON wrapper_links(video_id);