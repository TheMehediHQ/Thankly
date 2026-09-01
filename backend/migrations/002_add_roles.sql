ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) NOT NULL DEFAULT 'user';

-- Roles: user, admin, super_admin
-- super_admin can do everything
-- admin can manage users and content
-- user is the default
