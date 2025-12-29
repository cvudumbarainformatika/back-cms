-- Create Homepage Table (Single row table for homepage configuration)
CREATE TABLE IF NOT EXISTS homepage (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    hero JSON COMMENT 'Hero section data',
    stats JSON COMMENT 'Array of statistics',
    features JSON COMMENT 'Array of features',
    seo JSON COMMENT 'SEO metadata',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Organization Profile Table
CREATE TABLE IF NOT EXISTS profil_organisasi (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    visi_misi TEXT,
    sejarah TEXT,
    ad_art TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
