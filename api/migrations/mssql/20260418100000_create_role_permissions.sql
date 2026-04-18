-- Migration: explicit role_permissions join table for the GPUserRole <->
-- GPPermission many2many relationship (gorm:"many2many:role_permissions;").
-- Previously relied on AutoMigrate; this makes the table self-documenting
-- and survivable on non-AutoMigrate deployments.

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'role_permissions')
BEGIN
    CREATE TABLE role_permissions (
        gp_user_role_id  INT NOT NULL,
        gp_permission_id INT NOT NULL,
        CONSTRAINT pk_role_permissions PRIMARY KEY (gp_user_role_id, gp_permission_id)
    );
END;

IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = 'idx_role_permissions_permission' AND object_id = OBJECT_ID('role_permissions'))
BEGIN
    CREATE INDEX idx_role_permissions_permission ON role_permissions (gp_permission_id);
END;
