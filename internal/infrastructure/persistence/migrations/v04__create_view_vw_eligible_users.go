-- =====================================================
-- VIEW: Users eligible for sync
-- =====================================================

CREATE VIEW `vw_eligible_users` AS
SELECT DISTINCT d.user_id
FROM `ma_users_special_designations` d
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