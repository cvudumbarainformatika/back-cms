-- Add author_id and rejection fields to berita table

ALTER TABLE `berita` ADD COLUMN `author_id` BIGINT DEFAULT NULL AFTER `author`;
ALTER TABLE `berita` ADD INDEX `idx_author_id` (`author_id`);

ALTER TABLE `berita` ADD COLUMN `rejection_reason` TEXT DEFAULT NULL AFTER `status`;
ALTER TABLE `berita` ADD COLUMN `rejected_at` TIMESTAMP NULL DEFAULT NULL AFTER `rejection_reason`;
ALTER TABLE `berita` ADD COLUMN `rejected_by` VARCHAR(100) DEFAULT NULL AFTER `rejected_at`;
