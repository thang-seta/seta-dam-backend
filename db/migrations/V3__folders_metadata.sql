-- V3__seed_demo_folders.sql

INSERT INTO folders (name, description, created_by)
SELECT 'Open Images Sample', 'Sample imported image text metadata', u.id
FROM users u
WHERE u.email = 'thang.demo@gmail.com';

INSERT INTO folders (name, description, created_by)
SELECT 'Internal Review', 'Folders used for review workflow demo', u.id
FROM users u
WHERE u.email = 'thang.demo@gmail.com';

INSERT INTO folders (name, description, parent_id, created_by)
SELECT child.name, child.description, parent.id, u.id
FROM folders parent
JOIN users u ON u.email = 'thang.demo@gmail.com'
CROSS JOIN (
  VALUES
    ('Animals', 'Animal-related image metadata'),
    ('Vehicles', 'Vehicle-related image metadata'),
    ('Food', 'Food-related image metadata'),
    ('People', 'People and portrait-related metadata'),
    ('Indoor Scenes', 'Indoor environment metadata'),
    ('Outdoor Scenes', 'Outdoor environment metadata')
) AS child(name, description)
WHERE parent.name = 'Open Images Sample';

INSERT INTO folders (name, description, parent_id, created_by)
SELECT child.name, child.description, parent.id, u.id
FROM folders parent
JOIN users u ON u.email = 'thang.demo@gmail.com'
CROSS JOIN (
  VALUES
    ('Need Review', 'Metadata items waiting for review'),
    ('Approved Metadata', 'Metadata items already reviewed and approved')
) AS child(name, description)
WHERE parent.name = 'Internal Review';