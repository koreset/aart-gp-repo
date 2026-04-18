-- Migration: explicit role_permissions join table for the GPUserRole <->
-- GPPermission many2many relationship (gorm:"many2many:role_permissions;").
-- Previously relied on AutoMigrate; this makes the table self-documenting
-- and survivable on non-AutoMigrate deployments.

CREATE TABLE IF NOT EXISTS role_permissions (
    gp_user_role_id  INT NOT NULL,
    gp_permission_id INT NOT NULL,
    PRIMARY KEY (gp_user_role_id, gp_permission_id)
);

-- Index on the permission side for reverse lookups ("which roles have
-- this permission?"); the primary key already indexes the role side.
SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'role_permissions' AND INDEX_NAME = 'idx_role_permissions_permission'),
        'SELECT 1',
        'CREATE INDEX idx_role_permissions_permission ON role_permissions (gp_permission_id)'
    )
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
