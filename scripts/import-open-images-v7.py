#!/usr/bin/env python3
"""
Import a small Open Images V7 metadata sample into the Docker PostgreSQL database.

The script downloads CSV metadata only. It does not download image files.
"""

from __future__ import annotations

import argparse
import csv
import io
import json
import subprocess
import sys
import tempfile
import urllib.request
from collections import defaultdict
from pathlib import Path


DEFAULT_IMAGE_INFO_URL = (
    "https://storage.googleapis.com/openimages/2018_04/validation/"
    "validation-images-with-rotation.csv"
)
DEFAULT_LABELS_URL = (
    "https://storage.googleapis.com/openimages/v5/"
    "validation-annotations-human-imagelabels-boxable.csv"
)
DEFAULT_CLASSES_URL = (
    "https://storage.googleapis.com/openimages/v7/"
    "oidv7-class-descriptions-boxable.csv"
)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Import Open Images V7 validation metadata into Docker Postgres."
    )
    parser.add_argument("--limit", type=int, default=100, help="metadata rows to import")
    parser.add_argument(
        "--root-folder",
        default="Open Images V7 Import",
        help="root folder created under the DAM folder tree",
    )
    parser.add_argument(
        "--created-by-email",
        default="thang.demo@gmail.com",
        help="existing user email used as created_by for folders and metadata",
    )
    parser.add_argument(
        "--min-confidence",
        type=float,
        default=1.0,
        help="minimum image-label confidence to include",
    )
    parser.add_argument(
        "--cache-dir",
        default=".cache/open-images-v7",
        help="local cache directory for downloaded CSV files",
    )
    parser.add_argument("--db-service", default="db", help="Docker Compose database service")
    parser.add_argument("--db-user", default="postgres", help="PostgreSQL user")
    parser.add_argument("--db-name", default="setadam", help="PostgreSQL database")
    parser.add_argument("--image-info-url", default=DEFAULT_IMAGE_INFO_URL)
    parser.add_argument("--labels-url", default=DEFAULT_LABELS_URL)
    parser.add_argument("--classes-url", default=DEFAULT_CLASSES_URL)
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="download and prepare rows without writing to Postgres",
    )
    return parser.parse_args()


def download(url: str, cache_dir: Path) -> Path:
    cache_dir.mkdir(parents=True, exist_ok=True)
    filename = url.rstrip("/").split("/")[-1]
    path = cache_dir / filename
    if path.exists() and path.stat().st_size > 0:
        return path

    print(f"Downloading {url}", file=sys.stderr)
    with urllib.request.urlopen(url, timeout=120) as response:
        with tempfile.NamedTemporaryFile(delete=False, dir=cache_dir) as tmp:
            tmp.write(response.read())
            tmp_path = Path(tmp.name)
    tmp_path.replace(path)
    return path


def load_classes(path: Path) -> dict[str, str]:
    with path.open(newline="", encoding="utf-8") as handle:
        return {
            row["LabelName"]: row["DisplayName"].strip()
            for row in csv.DictReader(handle)
            if row.get("LabelName") and row.get("DisplayName")
        }


def load_labels(path: Path, class_names: dict[str, str], min_confidence: float) -> dict[str, list[dict[str, str]]]:
    labels: dict[str, list[dict[str, str]]] = defaultdict(list)
    with path.open(newline="", encoding="utf-8") as handle:
        for row in csv.DictReader(handle):
            try:
                confidence = float(row.get("Confidence") or 0)
            except ValueError:
                continue

            if confidence < min_confidence:
                continue

            label_mid = row.get("LabelName", "")
            labels[row["ImageID"]].append(
                {
                    "mid": label_mid,
                    "name": class_names.get(label_mid, label_mid),
                    "source": row.get("Source", ""),
                    "confidence": confidence,
                }
            )
    return labels


def build_rows(
    image_info_path: Path,
    labels_by_image: dict[str, list[dict[str, str]]],
    limit: int,
    source_urls: dict[str, str],
) -> list[dict[str, str]]:
    rows: list[dict[str, str]] = []
    with image_info_path.open(newline="", encoding="utf-8") as handle:
        for image in csv.DictReader(handle):
            image_id = image["ImageID"]
            label_items = labels_by_image.get(image_id, [])
            if not label_items:
                continue

            labels = sorted({item["name"] for item in label_items})
            label_mids = sorted({item["mid"] for item in label_items})
            category = labels[0] if labels else "Uncategorized"
            title = (image.get("Title") or "").strip() or f"Open Images image {image_id}"
            source_url = image.get("OriginalLandingURL") or image.get("OriginalURL") or ""
            thumbnail_url = image.get("Thumbnail300KURL") or image.get("OriginalURL") or ""
            metadata_json = {
                "source": "open_images_v7",
                "split": image.get("Subset", "validation"),
                "image_id": image_id,
                "original_url": image.get("OriginalURL", ""),
                "original_landing_url": image.get("OriginalLandingURL", ""),
                "author_profile_url": image.get("AuthorProfileURL", ""),
                "original_size": image.get("OriginalSize", ""),
                "original_md5": image.get("OriginalMD5", ""),
                "rotation": image.get("Rotation", ""),
                "label_mids": label_mids,
                "label_sources": sorted({item["source"] for item in label_items if item["source"]}),
                "source_files": source_urls,
            }

            rows.append(
                {
                    "external_id": image_id,
                    "title": title,
                    "description": "Open Images V7 validation metadata with labels: "
                    + ", ".join(labels[:8]),
                    "labels_json": json.dumps(labels, ensure_ascii=False),
                    "category": category,
                    "source_url": source_url,
                    "thumbnail_url": thumbnail_url,
                    "license": image.get("License", ""),
                    "author": image.get("Author", ""),
                    "metadata_json": json.dumps(metadata_json, ensure_ascii=False, sort_keys=True),
                }
            )

            if len(rows) >= limit:
                break
    return rows


def sql_literal(value: str) -> str:
    return "'" + value.replace("'", "''") + "'"


def build_sql(rows: list[dict[str, str]], args: argparse.Namespace) -> str:
    output = io.StringIO()
    writer = csv.DictWriter(output, fieldnames=list(rows[0].keys()), lineterminator="\n")
    writer.writeheader()
    writer.writerows(rows)

    return f"""\\set ON_ERROR_STOP on
BEGIN;

CREATE TEMP TABLE tmp_open_images_import (
  external_id text NOT NULL,
  title text NOT NULL,
  description text NOT NULL,
  labels_json jsonb NOT NULL,
  category text NOT NULL,
  source_url text,
  thumbnail_url text,
  license text,
  author text,
  metadata_json jsonb NOT NULL
) ON COMMIT DROP;

COPY tmp_open_images_import (
  external_id,
  title,
  description,
  labels_json,
  category,
  source_url,
  thumbnail_url,
  license,
  author,
  metadata_json
) FROM STDIN WITH (FORMAT csv, HEADER true);
{output.getvalue()}\\.

DO $$
DECLARE
  admin_user_id uuid;
  root_folder_id uuid;
BEGIN
  SELECT id INTO admin_user_id
  FROM users
  WHERE email = {sql_literal(args.created_by_email)}
    AND deleted_at IS NULL;

  IF admin_user_id IS NULL THEN
    RAISE EXCEPTION 'No active user found for email %', {sql_literal(args.created_by_email)};
  END IF;

  SELECT id INTO root_folder_id
  FROM folders
  WHERE parent_id IS NULL
    AND name = {sql_literal(args.root_folder)}
    AND deleted_at IS NULL
  ORDER BY created_at
  LIMIT 1;

  IF root_folder_id IS NULL THEN
    INSERT INTO folders (name, description, created_by)
    VALUES (
      {sql_literal(args.root_folder)},
      'Imported text metadata from Open Images V7 validation CSV files',
      admin_user_id
    )
    RETURNING id INTO root_folder_id;
  END IF;

  INSERT INTO folders (parent_id, name, description, created_by)
  SELECT DISTINCT
    root_folder_id,
    import_rows.category,
    'Imported Open Images V7 metadata category',
    admin_user_id
  FROM tmp_open_images_import import_rows
  WHERE NOT EXISTS (
    SELECT 1
    FROM folders existing
    WHERE existing.parent_id = root_folder_id
      AND existing.name = import_rows.category
      AND existing.deleted_at IS NULL
  );
END $$;

WITH admin_user AS (
  SELECT id
  FROM users
  WHERE email = {sql_literal(args.created_by_email)}
    AND deleted_at IS NULL
),
root_folder AS (
  SELECT id
  FROM folders
  WHERE parent_id IS NULL
    AND name = {sql_literal(args.root_folder)}
    AND deleted_at IS NULL
  ORDER BY created_at
  LIMIT 1
),
category_folder AS (
  SELECT folders.id, folders.name
  FROM folders
  JOIN root_folder ON folders.parent_id = root_folder.id
  WHERE folders.deleted_at IS NULL
)
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
  notes,
  created_by,
  updated_by,
  created_at,
  updated_at,
  deleted_at
)
SELECT
  category_folder.id,
  import_rows.title,
  import_rows.description,
  ARRAY(SELECT jsonb_array_elements_text(import_rows.labels_json)),
  import_rows.category,
  'open_images_v7_validation',
  import_rows.external_id,
  import_rows.source_url,
  import_rows.thumbnail_url,
  import_rows.license,
  import_rows.author,
  import_rows.metadata_json,
  'Imported by scripts/import-open-images-v7.py',
  admin_user.id,
  admin_user.id,
  NOW(),
  NOW(),
  NULL
FROM tmp_open_images_import import_rows
JOIN category_folder ON category_folder.name = import_rows.category
CROSS JOIN admin_user
ON CONFLICT (external_source, external_id)
WHERE external_source IS NOT NULL
  AND external_id IS NOT NULL
DO UPDATE SET
  folder_id = EXCLUDED.folder_id,
  title = EXCLUDED.title,
  description = EXCLUDED.description,
  labels = EXCLUDED.labels,
  category = EXCLUDED.category,
  source_url = EXCLUDED.source_url,
  thumbnail_url = EXCLUDED.thumbnail_url,
  license = EXCLUDED.license,
  author = EXCLUDED.author,
  metadata_json = EXCLUDED.metadata_json,
  notes = EXCLUDED.notes,
  updated_by = EXCLUDED.updated_by,
  updated_at = NOW(),
  deleted_at = NULL;

COMMIT;
"""


def run_psql(sql: str, args: argparse.Namespace) -> None:
    command = [
        "docker",
        "compose",
        "exec",
        "-T",
        args.db_service,
        "psql",
        "-U",
        args.db_user,
        "-d",
        args.db_name,
    ]
    subprocess.run(command, input=sql.encode("utf-8"), check=True)


def main() -> int:
    args = parse_args()
    if args.limit < 1:
        print("--limit must be greater than 0", file=sys.stderr)
        return 2

    cache_dir = Path(args.cache_dir)
    class_path = download(args.classes_url, cache_dir)
    labels_path = download(args.labels_url, cache_dir)
    image_info_path = download(args.image_info_url, cache_dir)

    class_names = load_classes(class_path)
    labels_by_image = load_labels(labels_path, class_names, args.min_confidence)
    rows = build_rows(
        image_info_path,
        labels_by_image,
        args.limit,
        {
            "image_info": args.image_info_url,
            "labels": args.labels_url,
            "classes": args.classes_url,
        },
    )

    if not rows:
        print("No rows matched the selected import options.", file=sys.stderr)
        return 1

    print(f"Prepared {len(rows)} Open Images metadata rows.")
    if args.dry_run:
        preview = json.dumps(rows[:3], indent=2, ensure_ascii=False)
        print(preview)
        return 0

    run_psql(build_sql(rows, args), args)
    print(
        f"Imported {len(rows)} rows into metadata_items "
        f"under '{args.root_folder}'."
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
