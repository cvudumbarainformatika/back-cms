-- Create Menus Table
CREATE TABLE IF NOT EXISTS menus (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    label VARCHAR(255) NOT NULL,
    slug VARCHAR(255),
    to VARCHAR(255),
    icon VARCHAR(100),
    parent_id BIGINT,
    position VARCHAR(50) COMMENT 'header, sidebar, footer',
    order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    is_fixed BOOLEAN DEFAULT FALSE,
    roles JSON COMMENT 'Array of roles: ["public", "member", "admin_cabang", etc.]',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_menus_parent_id 
        FOREIGN KEY (parent_id) REFERENCES menus(id) ON DELETE CASCADE,
    INDEX idx_position (position),
    INDEX idx_parent_id (parent_id),
    INDEX idx_order (order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
