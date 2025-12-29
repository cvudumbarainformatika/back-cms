-- Create Dynamic Contents Table
CREATE TABLE IF NOT EXISTS dynamic_contents (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    body LONGTEXT COMMENT 'Markdown content',
    html LONGTEXT COMMENT 'Rendered HTML',
    date DATE,
    image JSON COMMENT 'Image data (URL, alt, etc.)',
    authors JSON COMMENT 'Array of author names or IDs',
    badge JSON COMMENT 'Badge data',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_slug (slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
