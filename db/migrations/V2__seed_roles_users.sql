-- V2__seed_roles_permissions_users.sql

INSERT INTO users (id, email, display_name)
VALUES
  ('00000000-0000-0000-0000-000000000001', 'thang.demo@gmail.com', 'Thang Bui'),
  ('00000000-0000-0000-0000-000000000002', 'hung.demo@gmail.com', 'Hung Nguyen Anh'),
  ('00000000-0000-0000-0000-00000000002b', 'minh.demo@gmail.com', 'Minh Pham'),
  ('00000000-0000-0000-0000-000000000003', 'linh.demo@gmail.com', 'Linh Tran'),
  ('00000000-0000-0000-0000-00000000003b', 'mai.demo@gmail.com', 'Mai Le'),
  ('00000000-0000-0000-0000-00000000003c', 'quan.demo@gmail.com', 'Quan Do');

INSERT INTO roles (code, name, description)
VALUES
  ('trainer_admin', 'Trainer Admin', 'Can manage all folders, metadata, and permissions'),
  ('editor', 'Editor', 'Can read and write folders and metadata'),
  ('viewer', 'Viewer', 'Can read folders and metadata only');

INSERT INTO permission_actions (code, description)
VALUES
  ('read', 'Read resource'),
  ('write', 'Create or update resource'),
  ('manage_permissions', 'Grant or revoke object-level permissions');

-- trainer_admin: all permissions on folders and metadata items
INSERT INTO role_permissions (role_id, action_id, resource_type)
SELECT r.id, a.id, rt.resource_type::resource_type
FROM roles r
CROSS JOIN permission_actions a
CROSS JOIN (
  VALUES
    ('folder'),
    ('metadata_item')
) AS rt(resource_type)
WHERE r.code = 'trainer_admin';

-- editor: read/write folders and metadata items
INSERT INTO role_permissions (role_id, action_id, resource_type)
SELECT r.id, a.id, rt.resource_type::resource_type
FROM roles r
JOIN permission_actions a ON a.code IN ('read', 'write')
CROSS JOIN (
  VALUES
    ('folder'),
    ('metadata_item')
) AS rt(resource_type)
WHERE r.code = 'editor';

-- viewer: read-only folders and metadata items
INSERT INTO role_permissions (role_id, action_id, resource_type)
SELECT r.id, a.id, rt.resource_type::resource_type
FROM roles r
JOIN permission_actions a ON a.code = 'read'
CROSS JOIN (
  VALUES
    ('folder'),
    ('metadata_item')
) AS rt(resource_type)
WHERE r.code = 'viewer';

-- Assign roles
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
JOIN roles r ON r.code = 'trainer_admin'
WHERE u.email = 'thang.demo@gmail.com';

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
JOIN roles r ON r.code = 'editor'
WHERE u.email IN (
  'hung.demo@gmail.com',
  'minh.demo@gmail.com'
);

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
JOIN roles r ON r.code = 'viewer'
WHERE u.email IN (
  'linh.demo@gmail.com',
  'mai.demo@gmail.com',
  'quan.demo@gmail.com'
);