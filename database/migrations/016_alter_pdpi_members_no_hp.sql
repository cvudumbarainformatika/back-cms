-- Alter pdpi_members table to increase no_hp column size
-- Fix: Data too long for column 'no_hp' error during PDPI sync

ALTER TABLE pdpi_members 
MODIFY COLUMN no_hp VARCHAR(100) COMMENT 'Phone number, increased from VARCHAR(20) to accommodate longer formats';
