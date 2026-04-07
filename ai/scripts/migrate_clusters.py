"""One-time migration: import clusters.json into yarr DB AI cluster tables.

Run this after deploying the new Go binary (which runs migration m12 on startup)
to preserve your existing topic clusters without re-running the full clustering pipeline.

Usage:
    python -m ai.scripts.migrate_clusters /path/to/clusters.json /path/to/yarr.db

Example (on the server):
    python -m ai.scripts.migrate_clusters \\
        /mnt/seagate/yarr/data/clusters.json \\
        /mnt/seagate/yarr/data/yarr.db
"""

import argparse
import base64
import json
import sqlite3
import sys
from pathlib import Path


def open_db(path: str) -> sqlite3.Connection:
    conn = sqlite3.connect(path)
    conn.row_factory = sqlite3.Row
    conn.execute("PRAGMA journal_mode=WAL")
    conn.execute("PRAGMA foreign_keys=ON")
    return conn


def check_tables_exist(conn: sqlite3.Connection) -> bool:
    """Verify that m12 migration has already run (Go binary was started at least once)."""
    row = conn.execute(
        "SELECT name FROM sqlite_master WHERE type='table' AND name='ai_cluster_runs'"
    ).fetchone()
    return row is not None


def migrate(clusters_json_path: str, yarr_db_path: str) -> None:
    clusters_path = Path(clusters_json_path)
    db_path = Path(yarr_db_path)

    if not clusters_path.exists():
        print(f"Error: clusters.json not found at {clusters_path}", file=sys.stderr)
        sys.exit(1)

    if not db_path.exists():
        print(f"Error: yarr.db not found at {db_path}", file=sys.stderr)
        sys.exit(1)

    print(f"Loading {clusters_path} ({clusters_path.stat().st_size / 1024 / 1024:.1f} MB)...")
    with open(clusters_path) as f:
        data = json.load(f)

    clusters = data.get("clusters", [])
    if not clusters:
        print("No clusters found in clusters.json — nothing to migrate.")
        sys.exit(0)

    print(f"Found {len(clusters)} cluster entries.")

    conn = open_db(yarr_db_path)

    if not check_tables_exist(conn):
        print(
            "Error: AI cluster tables not found in yarr.db.\n"
            "Start the new yarr binary at least once to run migration m12, then re-run this script.",
            file=sys.stderr,
        )
        conn.close()
        sys.exit(1)

    # Check if a run already exists
    existing = conn.execute("SELECT COUNT(*) FROM ai_cluster_runs").fetchone()[0]
    if existing > 0:
        print(f"Warning: {existing} cluster run(s) already exist in the DB.")
        answer = input("Overwrite? This will delete existing cluster data. [y/N] ").strip().lower()
        if answer != "y":
            print("Aborted.")
            conn.close()
            sys.exit(0)
        conn.execute("DELETE FROM ai_cluster_runs")
        conn.commit()
        print("Existing cluster data deleted.")

    generated_at = data.get("generated_at", "")
    algorithm = data.get("algorithm", "hdbscan")
    n_articles = data.get("n_articles", 0)
    n_noise = data.get("n_noise", 0)

    print(f"Inserting run: generated_at={generated_at}, algorithm={algorithm}, "
          f"n_articles={n_articles}, n_noise={n_noise}")

    cur = conn.cursor()

    # Insert run
    cur.execute(
        "INSERT INTO ai_cluster_runs (generated_at, algorithm, n_articles, n_noise) VALUES (?, ?, ?, ?)",
        (generated_at, algorithm, n_articles, n_noise),
    )
    run_id = cur.lastrowid

    # Insert labels and centroids
    label_rows = []
    centroid_rows = []
    skipped_centroids = 0

    for c in clusters:
        cid = c.get("id")
        label = c.get("label", f"Cluster {cid}")
        article_count = c.get("article_count", 0)
        centroid = c.get("centroid", [])

        label_rows.append((run_id, label, article_count))

        if centroid and cid is not None:
            try:
                import numpy as np
                vec = np.array(centroid, dtype=np.float64)
                blob = vec.tobytes()
                centroid_rows.append((run_id, int(cid), label, blob))
            except Exception as e:
                skipped_centroids += 1
                print(f"  Warning: could not encode centroid for cluster {cid}: {e}")

    cur.executemany(
        "INSERT INTO ai_cluster_labels (run_id, label, article_count) VALUES (?, ?, ?)",
        label_rows,
    )
    if centroid_rows:
        cur.executemany(
            "INSERT INTO ai_cluster_centroids (run_id, cluster_id, label, centroid) VALUES (?, ?, ?, ?)",
            centroid_rows,
        )

    conn.commit()
    conn.close()

    centroid_bytes = sum(len(r[3]) for r in centroid_rows)
    print(f"\nMigration complete:")
    print(f"  {len(label_rows)} cluster labels inserted")
    print(f"  {len(centroid_rows)} centroids stored ({centroid_bytes / 1024:.1f} KB)")
    if skipped_centroids:
        print(f"  {skipped_centroids} centroids skipped (encoding errors)")
    print(f"\nYou can now delete clusters.json to reclaim disk space:")
    print(f"  rm {clusters_path}")


def main():
    parser = argparse.ArgumentParser(
        description="Migrate clusters.json to yarr DB AI cluster tables."
    )
    parser.add_argument("clusters_json", help="Path to clusters.json")
    parser.add_argument("yarr_db", help="Path to yarr.db")
    args = parser.parse_args()
    migrate(args.clusters_json, args.yarr_db)


if __name__ == "__main__":
    main()
