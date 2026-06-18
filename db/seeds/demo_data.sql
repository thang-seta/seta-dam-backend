-- demo_data.sql
-- Seed demo folders, metadata, and permissions

-- Seed folders
INSERT INTO folders (id, name, parent_id) VALUES
('11111111-1111-1111-1111-111111111111', 'Dataset A', NULL),
('22222222-2222-2222-2222-222222222222', 'Images', '11111111-1111-1111-1111-111111111111'),
('33333333-3333-3333-3333-333333333333', 'Secret Dataset', NULL)
ON CONFLICT (id) DO NOTHING;

-- Seed metadata items
INSERT INTO metadata (id, folder_id, title, description, labels, category, source_url, notes) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '22222222-2222-2222-2222-222222222222', 'Cat Image', 'A close-up shot of a cute cat.', ARRAY['animal', 'cat', 'cute'], 'Pets', 'https://example.com/cat.jpg', 'Dataset A -> Images'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '22222222-2222-2222-2222-222222222222', 'Dog Image', 'A dog playing in the park.', ARRAY['animal', 'dog', 'outdoor'], 'Pets', 'https://example.com/dog.jpg', 'Dataset A -> Images'),
('cccccccc-cccc-cccc-cccc-cccccccccccc', '33333333-3333-3333-3333-333333333333', 'Classified Satellite', 'Top secret satellite imagery.', ARRAY['military', 'satellite', 'aerial'], 'Classified', 'https://example.com/secret.jpg', 'Secret Dataset')
ON CONFLICT (id) DO NOTHING;

-- Seed object-level permissions
-- viewer_user (00000000-0000-0000-0000-000000000003) can READ Dataset A (11111111...)
INSERT INTO object_permissions (user_id, object_type, object_id, action) VALUES
('00000000-0000-0000-0000-000000000003', 'folder', '11111111-1111-1111-1111-111111111111', 'read')
ON CONFLICT ON CONSTRAINT unique_user_object_action DO NOTHING;

-- editor_user (00000000-0000-0000-0000-000000000002) can WRITE and READ Dataset A
INSERT INTO object_permissions (user_id, object_type, object_id, action) VALUES
('00000000-0000-0000-0000-000000000002', 'folder', '11111111-1111-1111-1111-111111111111', 'read'),
('00000000-0000-0000-0000-000000000002', 'folder', '11111111-1111-1111-1111-111111111111', 'write')
ON CONFLICT ON CONSTRAINT unique_user_object_action DO NOTHING;

-- editor_user can READ Secret Dataset (33333333...) but cannot write to it
INSERT INTO object_permissions (user_id, object_type, object_id, action) VALUES
('00000000-0000-0000-0000-000000000002', 'folder', '33333333-3333-3333-3333-333333333333', 'read')
ON CONFLICT ON CONSTRAINT unique_user_object_action DO NOTHING;
