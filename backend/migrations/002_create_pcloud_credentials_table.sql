CREATE TABLE pcloud_credentials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_name VARCHAR(100) NOT NULL,
    api_token TEXT NOT NULL,
    access_token TEXT,
    refresh_token TEXT,
    token_expires_at TIMESTAMP,
    storage_used_gb DECIMAL(10,2) DEFAULT 0,
    storage_limit_gb DECIMAL(10,2) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_pcloud_is_active ON pcloud_credentials(is_active);

-- Trigger for updated_at
CREATE TRIGGER update_pcloud_credentials_updated_at
    BEFORE UPDATE ON pcloud_credentials
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();