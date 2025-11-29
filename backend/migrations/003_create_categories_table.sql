CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_active ON categories(is_active);
CREATE INDEX idx_categories_order ON categories(display_order);

-- Trigger for updated_at
CREATE TRIGGER update_categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Seed default categories for adult content platform
INSERT INTO categories (name, slug, description, display_order, icon) VALUES
('Mahasiswi', 'mahasiswi', 'Konten mahasiswi Indonesia', 1, '🎓'),
('SMA', 'sma', 'Konten anak SMA', 2, '📚'),
('ABG', 'abg', 'Anak baru gede', 3, '👧'),
('Tante', 'tante', 'Wanita dewasa', 4, '👩'),
('Jilbab', 'jilbab', 'Berjilbab', 5, '🧕'),
('Indo', 'indo', 'Indonesia asli', 6, '🇮🇩'),
('Colmek', 'colmek', 'Coli memek', 7, '💦'),
('Live', 'live', 'Live show', 8, '🔴'),
('Viral', 'viral', 'Video viral terbaru', 9, '🔥'),
('Premium', 'premium', 'Konten premium eksklusif', 10, '⭐');