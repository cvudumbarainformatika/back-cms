-- Alter Agenda Table to match requirements
-- 1. Change date and end_date to DATETIME to support time
-- 2. Change fee to VARCHAR to support text descriptions
-- 3. Add published_at for consistency

ALTER TABLE agenda
MODIFY COLUMN date DATETIME NOT NULL,
MODIFY COLUMN end_date DATETIME NULL,
MODIFY COLUMN fee VARCHAR(255) COMMENT 'Fee description e.g., Gratis, Rp 50.000',
ADD COLUMN published_at TIMESTAMP NULL AFTER status;

-- Add index for published_at
CREATE INDEX idx_agenda_published_at ON agenda(published_at);
