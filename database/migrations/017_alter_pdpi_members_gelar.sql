-- Alter pdpi_members table to increase gelar and gelar2 column sizes
-- Fix: Data too long for column 'gelar2' error during PDPI sync

ALTER TABLE pdpi_members 
MODIFY COLUMN gelar VARCHAR(100) COMMENT 'Academic title prefix, increased from VARCHAR(50)';

ALTER TABLE pdpi_members 
MODIFY COLUMN gelar2 VARCHAR(100) COMMENT 'Academic title suffix, increased from VARCHAR(50)';
