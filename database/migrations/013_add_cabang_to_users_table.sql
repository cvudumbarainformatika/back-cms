-- Add cabang field to users table
-- This field stores the branch/office location of the user

ALTER TABLE users ADD COLUMN cabang VARCHAR(255) NULL COMMENT 'Branch/office location' AFTER status;

-- Add index for better query performance
ALTER TABLE users ADD INDEX idx_cabang (cabang);
