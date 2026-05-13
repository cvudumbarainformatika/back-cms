-- Add missing fields to pdpi_members table for complete PDPI sync
-- These fields are in struct but missing from initial migration (migration 015)

ALTER TABLE pdpi_members
ADD COLUMN npa_numeric BIGINT COMMENT 'Numeric NPA for sorting/filtering' AFTER npa,
ADD COLUMN foto TEXT COMMENT 'Photo URL' AFTER npa_numeric,
ADD COLUMN kota_kabupaten_kantor VARCHAR(100) COMMENT 'Office city/regency' AFTER kota_kabupaten,
ADD COLUMN provinsi_kantor VARCHAR(100) COMMENT 'Office province' AFTER kota_kabupaten_kantor,
ADD COLUMN tempat_praktek_1_tipe VARCHAR(255) COMMENT 'Practice 1 type' AFTER tempat_praktek_1,
ADD COLUMN tempat_praktek_1_tipe_2 VARCHAR(255) COMMENT 'Practice 1 type 2' AFTER tempat_praktek_1_tipe,
ADD COLUMN tempat_praktek_1_alkes VARCHAR(255) COMMENT 'Practice 1 medical equipment' AFTER tempat_praktek_1_tipe_2,
ADD COLUMN tempat_praktek_1_alkes_2 VARCHAR(255) COMMENT 'Practice 1 medical equipment 2' AFTER tempat_praktek_1_alkes,
ADD COLUMN tempat_praktek_2_tipe VARCHAR(255) COMMENT 'Practice 2 type' AFTER tempat_praktek_2,
ADD COLUMN tempat_praktek_2_tipe_2 VARCHAR(255) COMMENT 'Practice 2 type 2' AFTER tempat_praktek_2_tipe,
ADD COLUMN tempat_praktek_2_alkes VARCHAR(255) COMMENT 'Practice 2 medical equipment' AFTER tempat_praktek_2_tipe_2,
ADD COLUMN tempat_praktek_2_alkes_2 VARCHAR(255) COMMENT 'Practice 2 medical equipment 2' AFTER tempat_praktek_2_alkes,
ADD COLUMN kota_kabupaten_praktek_2 VARCHAR(100) COMMENT 'Practice 2 city/regency' AFTER tempat_praktek_2_alkes_2,
ADD COLUMN provinsi_praktek_2 VARCHAR(100) COMMENT 'Practice 2 province' AFTER kota_kabupaten_praktek_2,
ADD COLUMN tempat_praktek_3 VARCHAR(255) COMMENT 'Practice 3 location' AFTER provinsi_praktek_2,
ADD COLUMN tempat_praktek_3_tipe VARCHAR(255) COMMENT 'Practice 3 type' AFTER tempat_praktek_3,
ADD COLUMN tempat_praktek_3_tipe_2 VARCHAR(255) COMMENT 'Practice 3 type 2' AFTER tempat_praktek_3_tipe,
ADD COLUMN tempat_praktek_3_alkes VARCHAR(255) COMMENT 'Practice 3 medical equipment' AFTER tempat_praktek_3_tipe_2,
ADD COLUMN tempat_praktek_3_alkes_2 VARCHAR(255) COMMENT 'Practice 3 medical equipment 2' AFTER tempat_praktek_3_alkes,
ADD COLUMN kota_kabupaten_praktek_3 VARCHAR(100) COMMENT 'Practice 3 city/regency' AFTER tempat_praktek_3_alkes_2,
ADD COLUMN provinsi_praktek_3 VARCHAR(100) COMMENT 'Practice 3 province' AFTER kota_kabupaten_praktek_3,
ADD COLUMN gelar_fisr VARCHAR(100) COMMENT 'FISR degree' AFTER subspesialis;
