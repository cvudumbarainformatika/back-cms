-- Migration Template
-- File naming convention: NNN_description_in_snake_case.sql
-- Example: 001_create_users_table.sql, 002_create_posts_table.sql, etc.

-- ============================================================================
-- EXAMPLE: Create Examples Table
-- ============================================================================
-- This template shows how to create a basic table with common columns

CREATE TABLE IF NOT EXISTS examples (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Indexes for better query performance
    INDEX idx_email (email),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- NOTES for your migrations:
-- ============================================================================

-- MySQL Data Types:
-- INT, BIGINT - Integer values
-- VARCHAR(255) - String with max length
-- TEXT, LONGTEXT - Large text
-- DATE - Date only (YYYY-MM-DD)
-- DATETIME - Date and time (YYYY-MM-DD HH:MM:SS)
-- TIMESTAMP - Automatically set to current time
-- BOOLEAN - True/False (0/1)
-- DECIMAL(10,2) - Decimal for prices/money
-- JSON - JSON data type
-- ENUM('value1', 'value2') - Enumerated values

-- Best Practices:
-- 1. Always use AUTO_INCREMENT PRIMARY KEY for id
-- 2. Use TIMESTAMP DEFAULT CURRENT_TIMESTAMP for created_at
-- 3. Use TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP for updated_at
-- 4. Use proper collation: utf8mb4_unicode_ci (supports emoji)
-- 5. Add indexes on columns used in WHERE, JOIN, ORDER BY clauses
-- 6. Use NOT NULL for required fields, NULL for optional
-- 7. Add meaningful comments to complex tables
-- 8. Consider foreign key constraints for relationships

-- Foreign Key Example:
-- ALTER TABLE posts ADD CONSTRAINT fk_posts_user_id 
-- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- ============================================================================
-- TEMPLATE: One-to-Many Relationship
-- ============================================================================

-- Parent table
CREATE TABLE IF NOT EXISTS categories (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_slug (slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Child table with foreign key
CREATE TABLE IF NOT EXISTS products (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    category_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_category_id (category_id),
    CONSTRAINT fk_products_category_id 
        FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- TEMPLATE: Many-to-Many Relationship (Junction Table)
-- ============================================================================

-- User table
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Role table
CREATE TABLE IF NOT EXISTS roles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Junction table (Many-to-Many)
CREATE TABLE IF NOT EXISTS user_roles (
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_user_roles_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role_id 
        FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- How to Run Migrations
-- ============================================================================
-- 1. Create .sql files in database/migrations/ directory
-- 2. Name them sequentially: 001_*, 002_*, etc.
-- 3. Execute them in order against your database
-- 4. Keep all migrations and never delete them (for history)

-- Manual execution:
-- mysql -u username -p database_name < database/migrations/001_create_users_table.sql
-- Or use MySQL client in your IDE/tool

-- ============================================================================
-- Important Notes
-- ============================================================================
-- - Always set utf8mb4_unicode_ci for emoji support
-- - Use InnoDB engine for foreign key support
-- - Add timestamps (created_at, updated_at) to track changes
-- - Use meaningful index names
-- - Document complex migrations with comments
-- - Test migrations on staging before production
-- - Rollback plan: Keep backup of old schema
