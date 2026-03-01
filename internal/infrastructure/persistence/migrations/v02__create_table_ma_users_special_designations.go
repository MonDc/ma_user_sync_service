-- =====================================================
-- MA_SPECIAL_DESIGNATIONS: Who qualifies for sync
-- This table drives the sync logic - only users here
-- get copied/synced to ma_users
-- =====================================================

CREATE TABLE `ma_users_special_designations` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `user_id` BIGINT UNSIGNED NOT NULL, -- References ma_users.id logically
    
    `designation_type` VARCHAR(50) NOT NULL, -- 'merchant', 'premium', 'vip', etc.
    `designation_reason` VARCHAR(255), -- Why they qualified
    `designation_source` VARCHAR(50), -- 'event', 'manual', 'migration'
    
    -- Event tracking
    `event_id` VARCHAR(255), -- Original event ID that triggered this
    `event_type` VARCHAR(100), -- 'UserPromoted', 'MerchantApproved', etc.
    `event_time` TIMESTAMP NOT NULL, -- When event occurred
    
    -- Immutability guarantee
    `granted_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `granted_by` VARCHAR(255), -- System or user who granted
    
    -- Revocation (rare, but tracked)
    `revoked_at` TIMESTAMP NULL,
    `revoked_reason` VARCHAR(255),
    
    -- Uniqueness: one user can have multiple designations, but only one active per type
    UNIQUE KEY `uniq_active_designation` (`user_id`, `designation_type`, `revoked_at`),
    
    INDEX `idx_designations_user` (`user_id`),
    INDEX `idx_designations_type` (`designation_type`),
    INDEX `idx_designations_event` (`event_id`)
    
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;