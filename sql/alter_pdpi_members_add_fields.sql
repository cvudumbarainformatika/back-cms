-- Add new columns to pdpi_members table to match Supabase schema

ALTER TABLE pdpi_members
ADD COLUMN npa_numeric BIGINT NULL AFTER npa,
ADD COLUMN foto TEXT NULL AFTER nama,
ADD COLUMN kota_kabupaten_kantor VARCHAR(255) NULL AFTER kota_kabupaten,
ADD COLUMN provinsi_kantor VARCHAR(255) NULL AFTER provinsi,
ADD COLUMN tempat_praktek_1_tipe VARCHAR(100) NULL AFTER tempat_praktek_1,
ADD COLUMN tempat_praktek_1_tipe_2 VARCHAR(100) NULL AFTER tempat_praktek_1_tipe,
ADD COLUMN tempat_praktek_1_alkes VARCHAR(100) NULL AFTER tempat_praktek_1_tipe_2,
ADD COLUMN tempat_praktek_1_alkes_2 VARCHAR(100) NULL AFTER tempat_praktek_1_alkes,
ADD COLUMN tempat_praktek_2_tipe VARCHAR(100) NULL AFTER tempat_praktek_2,
ADD COLUMN tempat_praktek_2_tipe_2 VARCHAR(100) NULL AFTER tempat_praktek_2_tipe,
ADD COLUMN tempat_praktek_2_alkes VARCHAR(100) NULL AFTER tempat_praktek_2_tipe_2,
ADD COLUMN tempat_praktek_2_alkes_2 VARCHAR(100) NULL AFTER tempat_praktek_2_alkes,
ADD COLUMN kota_kabupaten_praktek_2 VARCHAR(255) NULL AFTER tempat_praktek_2_alkes_2,
ADD COLUMN provinsi_praktek_2 VARCHAR(255) NULL AFTER kota_kabupaten_praktek_2,
ADD COLUMN tempat_praktek_3 VARCHAR(255) NULL AFTER provinsi_praktek_2,
ADD COLUMN tempat_praktek_3_tipe VARCHAR(100) NULL AFTER tempat_praktek_3,
ADD COLUMN tempat_praktek_3_tipe_2 VARCHAR(100) NULL AFTER tempat_praktek_3_tipe,
ADD COLUMN tempat_praktek_3_alkes VARCHAR(100) NULL AFTER tempat_praktek_3_tipe_2,
ADD COLUMN tempat_praktek_3_alkes_2 VARCHAR(100) NULL AFTER tempat_praktek_3_alkes,
ADD COLUMN kota_kabupaten_praktek_3 VARCHAR(255) NULL AFTER tempat_praktek_3_alkes_2,
ADD COLUMN provinsi_praktek_3 VARCHAR(255) NULL AFTER kota_kabupaten_praktek_3,
ADD COLUMN gelar_fisr VARCHAR(50) NULL AFTER subspesialis;

-- Add index for sorting performance
CREATE INDEX idx_pdpi_members_npa_numeric ON pdpi_members(npa_numeric);
