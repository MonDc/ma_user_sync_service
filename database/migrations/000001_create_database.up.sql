-- =====================================================
-- CREATE DATABASE: ma_db
-- =====================================================
-- This is the primary database for the MA User Sync Service.
-- It stores synced user data and service-specific tables.
-- =====================================================

CREATE DATABASE IF NOT EXISTS ma_db
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

-- =====================================================
-- After creation, switch to the database and verify:
-- USE ma_db;
-- SHOW TABLES;
-- =====================================================