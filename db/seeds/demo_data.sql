-- demo_data.sql
-- Seed demo folders, metadata, and permissions

-- Seed folders
INSERT INTO folders (id, name, parent_id, created_by) VALUES
('11111111-1111-1111-1111-111111111111', 'Dataset A', NULL, '00000000-0000-0000-0000-000000000001'),
('22222222-2222-2222-2222-222222222222', 'Images', '11111111-1111-1111-1111-111111111111', '00000000-0000-0000-0000-000000000001'),
('33333333-3333-3333-3333-333333333333', 'Secret Dataset', NULL, '00000000-0000-0000-0000-000000000001')
ON CONFLICT (id) DO NOTHING;

-- Seed metadata items
INSERT INTO metadata_items (id, folder_id, title, description, labels, category, source_url, notes, created_by) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '22222222-2222-2222-2222-222222222222', 'Cat Image', 'A close-up shot of a cute cat.', ARRAY['animal', 'cat', 'cute'], 'Pets', 'https://example.com/cat.jpg', 'Dataset A -> Images', '00000000-0000-0000-0000-000000000001'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '22222222-2222-2222-2222-222222222222', 'Dog Image', 'A dog playing in the park.', ARRAY['animal', 'dog', 'outdoor'], 'Pets', 'https://example.com/dog.jpg', 'Dataset A -> Images', '00000000-0000-0000-0000-000000000001'),
('cccccccc-cccc-cccc-cccc-cccccccccccc', '33333333-3333-3333-3333-333333333333', 'Classified Satellite', 'Top secret satellite imagery.', ARRAY['military', 'satellite', 'aerial'], 'Classified', 'https://example.com/secret.jpg', 'Secret Dataset', '00000000-0000-0000-0000-000000000001')
ON CONFLICT (id) DO NOTHING;

-- Seed object-level permissions (folder_permissions)
-- viewer_user (00000000-0000-0000-0000-000000000003) can READ Dataset A (11111111...)
INSERT INTO folder_permissions (folder_id, grantee_user_id, action_id, granted_by)
SELECT 
  '11111111-1111-1111-1111-111111111111',
  '00000000-0000-0000-0000-000000000003',
  a.id,
  '00000000-0000-0000-0000-000000000001'
FROM permission_actions a
WHERE a.code = 'read'::permission_action_code
ON CONFLICT (folder_id, grantee_user_id, action_id) WHERE grantee_user_id IS NOT NULL DO NOTHING;

-- editor_user (00000000-0000-0000-0000-000000000002) can WRITE and READ Dataset A
INSERT INTO folder_permissions (folder_id, grantee_user_id, action_id, granted_by)
SELECT 
  '11111111-1111-1111-1111-111111111111',
  '00000000-0000-0000-0000-000000000002',
  a.id,
  '00000000-0000-0000-0000-000000000001'
FROM permission_actions a
WHERE a.code IN ('read'::permission_action_code, 'write'::permission_action_code)
ON CONFLICT (folder_id, grantee_user_id, action_id) WHERE grantee_user_id IS NOT NULL DO NOTHING;

-- editor_user can READ Secret Dataset (33333333...) but cannot write to it
INSERT INTO folder_permissions (folder_id, grantee_user_id, action_id, granted_by)
SELECT 
  '33333333-3333-3333-3333-333333333333',
  '00000000-0000-0000-0000-000000000002',
  a.id,
  '00000000-0000-0000-0000-000000000001'
FROM permission_actions a
WHERE a.code = 'read'::permission_action_code
ON CONFLICT (folder_id, grantee_user_id, action_id) WHERE grantee_user_id IS NOT NULL DO NOTHING;
