INSERT INTO metadata_items (
  folder_id,
  title,
  description,
  labels,
  category,
  external_source,
  external_id,
  source_url,
  thumbnail_url,
  license,
  author,
  metadata_json,
  created_by
)
SELECT
  f.id,
  'Dog standing on grass',
  'Sample text metadata imported from Open Images V7',
  ARRAY['Dog', 'Animal', 'Grass'],
  'animal',
  'open_images_v7',
  '000002b66c9c498e',
  'https://example.com/source-image',
  'https://example.com/thumbnail-image',
  'CC BY 2.0',
  'Open Images contributor',
  '{"source": "open_images_v7", "split": "validation"}'::jsonb,
  u.id
FROM folders f
JOIN users u ON u.email = 'thang.demo@gmail.com'
WHERE f.name = 'Animals';