-- =====================================================
-- MA_SYNC_AUDIT: Track all sync operations
-- =====================================================

CREATE TABLE `ma_users_sync_audit` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `sync_type` ENUM('full','incremental','single') NOT NULL,
    `started_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `completed_at` TIMESTAMP NULL,
    `status` ENUM('running','completed','failed') DEFAULT 'running',
    
    `users_processed` INT DEFAULT 0,
    `users_succeeded` INT DEFAULT 0,
    `users_failed` INT DEFAULT 0,
    
    `error_summary` TEXT,
    `triggered_by` VARCHAR(255), -- 'schedule', 'manual', 'event'
    
    INDEX `idx_sync_audit_status` (`status`),
    INDEX `idx_sync_audit_started` (`started_at`)
    
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- VIEW: Users eligible for sync
-- =====================================================

CREATE VIEW `vw_eligible_users` AS
SELECT DISTINCT d.user_id
FROM `ma_special_designations` d
WHERE d.revoked_at IS NULL;

-- =====================================================
-- VIEW: Complete user data with designations
-- =====================================================

CREATE VIEW `vw_users_with_designations` AS
SELECT 
    u.*,
    GROUP_CONCAT(d.designation_type) as active_designations
FROM `ma_users` u
LEFT JOIN `ma_special_designations` d ON u.id = d.user_id AND d.revoked_at IS NULL
GROUP BY u.id;