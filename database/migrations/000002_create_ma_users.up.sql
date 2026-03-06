-- =====================================================
-- MA_USERS: Local copy of designated users
-- This table mirrors mi_users structure but contains
-- ONLY users with special designations (merchants, etc.)
-- =====================================================

CREATE TABLE `ma_users` (
    -- Primary identifiers (mirror mi_users)
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY, -- Matches mi_users.id (no auto_increment)
    `public_id` CHAR(36) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
    
    -- Core identity
    `username` VARCHAR(50),
    `email` VARCHAR(254) NOT NULL,
    `first_name` VARCHAR(100) NOT NULL,
    `last_name` VARCHAR(100) NOT NULL,
    
    -- Mithaq-specific
    `mithaq_name` VARCHAR(100),
    `mithaq_email` VARCHAR(254),
    
    -- Contact & verification
    `mobile` VARCHAR(20),
    `mobile_verified_at` TIMESTAMP NULL,
    `email_verified_at` TIMESTAMP NULL,
    
    -- Personal info
    `date_of_birth` DATE NOT NULL,
    `address` VARCHAR(255),
    `avatar_url` TEXT,
    
    -- Status & classification
    `is_verified` BOOLEAN DEFAULT FALSE,
    `status` ENUM('active','suspended','banned','deleted','pending') DEFAULT 'active',
    `user_type` VARCHAR(50),
    `role` ENUM('user','admin','support') DEFAULT 'user',
    
    -- Compliance
    `kyc_level` ENUM('none','basic','verified','enhanced') DEFAULT 'none',
    `verification_number` VARCHAR(50),
    
    -- Relationships
    `referrer_id` BIGINT UNSIGNED,
    
    -- Localization & auth
    `locale` VARCHAR(10) DEFAULT 'en',
    `timezone` VARCHAR(50) DEFAULT 'UTC',
    `identity_provider` ENUM('local','google','apple','facebook','github','microsoft') DEFAULT 'local',
    `metadata` JSON,
    `consents` JSON,
    
    -- Original timestamps & audit (preserved from mi_users)
    `original_created_at` TIMESTAMP NOT NULL,
    `original_updated_at` TIMESTAMP NOT NULL,
    `original_last_login_at` TIMESTAMP NULL,
    `original_created_by` VARCHAR(255),
    `original_updated_by` VARCHAR(255),
    `original_created_ip` VARCHAR(45),
    
    -- Deletion tracking (mirror mi_users)
    `archived_at` TIMESTAMP NULL,
    `deleted_at` TIMESTAMP NULL,
    `deleted` TINYINT(1) DEFAULT 0,
    
    -- =================================================
    -- SYNC-SPECIFIC FIELDS (YOUR control)
    -- =================================================
    `first_synced_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- When first added to ma_users
    `last_synced_at` TIMESTAMP NULL, -- Last successful sync
    `sync_status` ENUM('pending','synced','failed','conflict') DEFAULT 'pending',
    `sync_attempts` INT DEFAULT 0,
    `sync_error` TEXT,
    `sync_version` BIGINT DEFAULT 1, -- For optimistic locking
    
    -- Local audit (YOUR system)
    `local_created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `local_updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Indexes for YOUR query patterns
    INDEX `idx_ma_users_public_id` (`public_id`),
    INDEX `idx_ma_users_email` (`email`),
    INDEX `idx_ma_users_status` (`status`),
    INDEX `idx_ma_users_sync_status` (`sync_status`),
    INDEX `idx_ma_users_last_synced` (`last_synced_at`),
    INDEX `idx_ma_users_deleted` (`deleted`)
    
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;