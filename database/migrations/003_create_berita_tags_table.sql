-- Create Berita Tags Table
CREATE TABLE IF NOT EXISTS berita_tags (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Berita Tag Map (Many-to-Many) Table
CREATE TABLE IF NOT EXISTS berita_tag_map (
    berita_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (berita_id, tag_id),
    CONSTRAINT fk_berita_tag_map_berita_id 
        FOREIGN KEY (berita_id) REFERENCES berita(id) ON DELETE CASCADE,
    CONSTRAINT fk_berita_tag_map_tag_id 
        FOREIGN KEY (tag_id) REFERENCES berita_tags(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
