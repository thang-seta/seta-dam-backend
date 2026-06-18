-- V2__seed_roles_users.sql
-- Seed roles and base users

INSERT INTO roles (name) VALUES
('trainer_admin'),
('editor'),
('viewer')
ON CONFLICT (name) DO NOTHING;

-- Seed default users with deterministic UUIDs for testability
INSERT INTO users (id, username) VALUES
('00000000-0000-0000-0000-000000000001', 'admin_user'),
('00000000-0000-0000-0000-000000000002', 'editor_user'),
('00000000-0000-0000-0000-000000000003', 'viewer_user')
ON CONFLICT (id) DO NOTHING;

INSERT INTO user_roles (user_id, role_name) VALUES
('00000000-0000-0000-0000-000000000001', 'trainer_admin'),
('00000000-0000-0000-0000-000000000002', 'editor'),
('00000000-0000-0000-0000-000000000003', 'viewer')
ON CONFLICT (user_id, role_name) DO NOTHING;
