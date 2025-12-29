-- Add Profile Fields to Users Table
-- This migration adds phone, address, bio, and avatar fields to support user profile updates

ALTER TABLE users ADD COLUMN phone VARCHAR(20) NULL AFTER status;
ALTER TABLE users ADD COLUMN address TEXT NULL AFTER phone;
ALTER TABLE users ADD COLUMN bio TEXT NULL AFTER address;
ALTER TABLE users ADD COLUMN avatar VARCHAR(255) NULL AFTER bio;
