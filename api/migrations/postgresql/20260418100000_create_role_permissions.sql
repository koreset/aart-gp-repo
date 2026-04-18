-- Migration: explicit role_permissions join table for the GPUserRole <->
-- GPPermission many2many relationship (gorm:"many2many:role_permissions;").
-- Previously relied on AutoMigrate; this makes the table self-documenting
-- and survivable on non-AutoMigrate deployments.

CREATE TABLE IF NOT EXISTS role_permissions (
    gp_user_role_id  INTEGER NOT NULL,
    gp_permission_id INTEGER NOT NULL,
    PRIMARY KEY (gp_user_role_id, gp_permission_id)
);

CREATE INDEX IF NOT EXISTS idx_role_permissions_permission
    ON role_permissions (gp_permission_id);
