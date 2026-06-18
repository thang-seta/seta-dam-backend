#!/usr/bin/env python3
import json
import os
import sys
import time
import urllib.error
import urllib.request


GRAPHQL_URL = os.getenv("GRAPHQL_URL", "http://localhost:4000/graphql")
CORE_HEALTH_URL = os.getenv("CORE_HEALTH_URL", "http://localhost:8080/health")
NODE_HEALTH_URL = os.getenv("NODE_HEALTH_URL", "http://localhost:4000/health")


def http_get_json(url):
    with urllib.request.urlopen(url, timeout=10) as response:
        return json.loads(response.read().decode("utf-8"))


def gql(query, variables=None, auth="admin_user", expect_error=False):
    body = json.dumps({"query": query, "variables": variables or {}}).encode("utf-8")
    request = urllib.request.Request(
        GRAPHQL_URL,
        data=body,
        headers={
            "Content-Type": "application/json",
            "Authorization": auth,
        },
        method="POST",
    )
    try:
        with urllib.request.urlopen(request, timeout=15) as response:
            payload = json.loads(response.read().decode("utf-8"))
    except urllib.error.HTTPError as exc:
        payload = json.loads(exc.read().decode("utf-8"))

    has_errors = bool(payload.get("errors"))
    if expect_error:
        if not has_errors:
            raise AssertionError(f"Expected GraphQL error but got: {payload}")
        return payload

    if has_errors:
        raise AssertionError(f"GraphQL returned errors: {json.dumps(payload, indent=2)}")
    return payload["data"]


def main():
    stamp = int(time.time())
    root_id = None
    child_id = None
    metadata_id = None

    print("Checking health endpoints...")
    assert http_get_json(CORE_HEALTH_URL)["status"] == "OK"
    assert http_get_json(NODE_HEALTH_URL)["status"] == "OK"

    try:
        print("Querying admin folder tree...")
        gql("query { folderTree { id name parent_id } }")

        print("Creating demo root and child folders...")
        data = gql(
            """
            mutation CreateFolder($name: String!) {
              createFolder(name: $name, description: "Smoke test root") {
                id
                name
              }
            }
            """,
            {"name": f"Smoke Test {stamp}"},
        )
        root_id = data["createFolder"]["id"]

        data = gql(
            """
            mutation CreateFolder($name: String!, $parentId: ID!) {
              createFolder(name: $name, parentId: $parentId, description: "Smoke test child") {
                id
                parent_id
              }
            }
            """,
            {"name": f"Smoke Child {stamp}", "parentId": root_id},
            auth="editor_user",
        )
        child_id = data["createFolder"]["id"]

        print("Creating and reading metadata...")
        data = gql(
            """
            mutation CreateMetadata($folderId: ID!, $title: String!, $metadataJson: String) {
              createMetadata(
                folderId: $folderId
                title: $title
                description: "Smoke metadata item"
                labels: ["smoke", "phase1"]
                category: "demo"
                sourceUrl: "https://example.com/smoke"
                metadataJson: $metadataJson
              ) {
                id
                folder_id
                title
              }
            }
            """,
            {
                "folderId": child_id,
                "title": f"Smoke Metadata {stamp}",
                "metadataJson": json.dumps({"source": "smoke-test"}),
            },
            auth="editor_user",
        )
        metadata_id = data["createMetadata"]["id"]

        gql(
            """
            query MetadataList($folderId: ID!) {
              metadataList(folderId: $folderId) {
                id
                title
              }
            }
            """,
            {"folderId": child_id},
            auth="viewer_user",
        )

        print("Verifying viewer write is denied...")
        gql(
            """
            mutation UpdateMetadata($id: ID!, $folderId: ID!) {
              updateMetadata(id: $id, folderId: $folderId, title: "Viewer should fail") {
                id
              }
            }
            """,
            {"id": metadata_id, "folderId": child_id},
            auth="viewer_user",
            expect_error=True,
        )

        print("Checking effective permissions...")
        data = gql(
            """
            query EffectivePermissions($userId: ID!, $objectId: ID!) {
              effectivePermissions(userId: $userId, objectType: "metadata_item", objectId: $objectId) {
                actions
              }
            }
            """,
            {
                "userId": "00000000-0000-0000-0000-000000000003",
                "objectId": metadata_id,
            },
        )
        if "read" not in data["effectivePermissions"]["actions"]:
            raise AssertionError("viewer_user should have read permission")

        print("Smoke test passed.")
        return 0
    finally:
        print("Cleaning up smoke-test data...")
        if metadata_id:
            gql(
                "mutation DeleteMetadata($id: ID!) { deleteMetadata(id: $id) }",
                {"id": metadata_id},
            )
        if child_id:
            gql(
                "mutation DeleteFolder($id: ID!) { deleteFolder(id: $id) }",
                {"id": child_id},
            )
        if root_id:
            gql(
                "mutation DeleteFolder($id: ID!) { deleteFolder(id: $id) }",
                {"id": root_id},
            )


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:
        print(f"Smoke test failed: {exc}", file=sys.stderr)
        raise SystemExit(1)
