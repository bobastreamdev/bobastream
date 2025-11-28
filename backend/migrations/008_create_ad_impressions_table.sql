CREATE TYPE impression_type AS ENUM ('view', 'click', 'skip');

CREATE TABLE ad_impressions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ad_id UUID NOT NULL REFERENCES ads(id) ON DELETE CASCADE,
    video_id UUID REFERENCES videos(id) ON DELETE SET NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    viewer_ip VARCHAR(45) NOT NULL,
    impression_type impression_type NOT NULL,
    watched_duration INTEGER DEFAULT 0,
    session_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    date DATE DEFAULT CURRENT_DATE
);

-- Indexes
CREATE INDEX idx_ad_impressions_ad_id ON ad_impressions(ad_id);
CREATE INDEX idx_ad_impressions_video_id ON ad_impressions(video_id);
CREATE INDEX idx_ad_impressions_date ON ad_impressions(date, ad_id);
CREATE INDEX idx_ad_impressions_type ON ad_impressions(impression_type);