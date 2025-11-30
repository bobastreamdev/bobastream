CREATE TYPE ad_type AS ENUM ('preroll', 'banner', 'popup');

CREATE TABLE ads (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    ad_type ad_type NOT NULL,
    content_url TEXT,
    redirect_url TEXT,
    duration_seconds INTEGER,
    display_frequency INTEGER DEFAULT 1,
    is_active BOOLEAN DEFAULT true,
    priority INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_ads_type ON ads(ad_type);
CREATE INDEX idx_ads_active ON ads(is_active);
CREATE INDEX idx_ads_priority ON ads(priority DESC);
CREATE INDEX idx_ads_deleted_at ON ads(deleted_at);

-- Trigger for updated_at
CREATE TRIGGER update_ads_updated_at
    BEFORE UPDATE ON ads
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();