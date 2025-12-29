-- Create Agenda Table
CREATE TABLE IF NOT EXISTS agenda (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) COMMENT 'webinar, workshop, seminar, kongres',
    date DATE NOT NULL,
    end_date DATE,
    is_online BOOLEAN DEFAULT FALSE,
    location VARCHAR(255),
    skp DECIMAL(10,2),
    quota INT,
    registration_url VARCHAR(255),
    image_url VARCHAR(255),
    fee DECIMAL(10,2),
    status VARCHAR(50) DEFAULT 'draft' COMMENT 'draft, published',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_slug (slug),
    INDEX idx_type (type),
    INDEX idx_date (date),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Agenda Registrations Table
CREATE TABLE IF NOT EXISTS agenda_registrations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    agenda_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    status VARCHAR(50) DEFAULT 'pending' COMMENT 'pending, confirmed, cancelled',
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_agenda_user (agenda_id, user_id),
    CONSTRAINT fk_agenda_registrations_agenda_id 
        FOREIGN KEY (agenda_id) REFERENCES agenda(id) ON DELETE CASCADE,
    CONSTRAINT fk_agenda_registrations_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_agenda_id (agenda_id),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
