-- Create PDPI Members Table
-- This table caches member data from PDPI API for performance and offline capability

CREATE TABLE IF NOT EXISTS pdpi_members (
    id VARCHAR(36) PRIMARY KEY COMMENT 'UUID from PDPI API',
    npa VARCHAR(20) UNIQUE NOT NULL COMMENT 'Nomor Peserta Anggota',
    nama VARCHAR(255) NOT NULL,
    gelar VARCHAR(50),
    gelar2 VARCHAR(50),
    email VARCHAR(255) UNIQUE,
    no_hp VARCHAR(20),
    nik VARCHAR(20),
    jenis_kelamin VARCHAR(20),
    tempat_lahir VARCHAR(100),
    tgl_lahir DATE,
    alamat_rumah TEXT,
    cabang VARCHAR(100),
    provinsi VARCHAR(100),
    kota_kabupaten VARCHAR(100),
    status VARCHAR(20) COMMENT 'Aktif, Non-Aktif',
    alumni VARCHAR(255),
    thn_lulus INT,
    tempat_tugas VARCHAR(255),
    tempat_praktek_1 VARCHAR(255),
    tempat_praktek_2 VARCHAR(255),
    subspesialis VARCHAR(100),
    no_str VARCHAR(50),
    str_berlaku_sampai DATE,
    no_sip VARCHAR(50),
    sip_berlaku_sampai DATE,
    
    user_id BIGINT COMMENT 'Foreign key to users table',
    synced_at TIMESTAMP NULL COMMENT 'Last sync time from PDPI API',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Foreign key constraint
    CONSTRAINT fk_pdpi_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    
    -- Indexes for better query performance
    INDEX idx_npa (npa),
    INDEX idx_email (email),
    INDEX idx_nik (nik),
    INDEX idx_cabang (cabang),
    INDEX idx_provinsi (provinsi),
    INDEX idx_status (status),
    INDEX idx_user_id (user_id),
    INDEX idx_synced_at (synced_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
