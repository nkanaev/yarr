"""Read-only access to yarr's SQLite database."""

import logging
import sqlite3

log = logging.getLogger(__name__)


def open_db(path: str) -> sqlite3.Connection:
    """Open yarr's SQLite DB in read-only mode."""
    conn = sqlite3.connect(f"file:{path}?mode=ro", uri=True)
    conn.row_factory = sqlite3.Row
    conn.execute("PRAGMA journal_mode=WAL")
    return conn


def list_items(
    conn: sqlite3.Connection,
    feed_id: int | None = None,
    folder_id: int | None = None,
    since_id: int | None = None,
    limit: int = 500,
) -> list[dict]:
    """List items with feed and folder metadata."""
    query = """
        SELECT i.id, i.title, i.link, i.content, i.date, i.status,
               i.feed_id, f.title as feed_title,
               COALESCE(fo.title, 'uncategorized') as folder_name
        FROM items i
        JOIN feeds f ON f.id = i.feed_id
        LEFT JOIN folders fo ON fo.id = f.folder_id
        WHERE 1=1
    """
    args: list = []
    if feed_id is not None:
        query += " AND i.feed_id = ?"
        args.append(feed_id)
    if folder_id is not None:
        query += " AND f.folder_id = ?"
        args.append(folder_id)
    if since_id is not None:
        query += " AND i.id > ?"
        args.append(since_id)
    query += " ORDER BY i.date DESC LIMIT ?"
    args.append(limit)

    rows = conn.execute(query, args).fetchall()
    return [dict(r) for r in rows]


def get_item(conn: sqlite3.Connection, item_id: int) -> dict | None:
    """Get a single item with full content."""
    row = conn.execute(
        """
        SELECT i.id, i.title, i.link, i.content, i.date, i.status,
               i.feed_id, f.title as feed_title,
               COALESCE(fo.title, 'uncategorized') as folder_name
        FROM items i
        JOIN feeds f ON f.id = i.feed_id
        LEFT JOIN folders fo ON fo.id = f.folder_id
        WHERE i.id = ?
        """,
        (item_id,),
    ).fetchone()
    return dict(row) if row else None


def get_items_by_ids(conn: sqlite3.Connection, ids: list[int]) -> list[dict]:
    """Batch fetch items by ID."""
    if not ids:
        return []
    placeholders = ",".join("?" for _ in ids)
    rows = conn.execute(
        f"""
        SELECT i.id, i.title, i.link, i.content, i.date, i.status,
               i.feed_id, f.title as feed_title,
               COALESCE(fo.title, 'uncategorized') as folder_name
        FROM items i
        JOIN feeds f ON f.id = i.feed_id
        LEFT JOIN folders fo ON fo.id = f.folder_id
        WHERE i.id IN ({placeholders})
        """,
        ids,
    ).fetchall()
    return [dict(r) for r in rows]


def list_feeds(conn: sqlite3.Connection) -> list[dict]:
    """List all feeds with folder names."""
    rows = conn.execute(
        """
        SELECT f.id, f.title, f.feed_link, f.link,
               COALESCE(fo.title, 'uncategorized') as folder_name
        FROM feeds f
        LEFT JOIN folders fo ON fo.id = f.folder_id
        ORDER BY f.title COLLATE NOCASE
        """
    ).fetchall()
    return [dict(r) for r in rows]


def list_all_items(conn: sqlite3.Connection) -> list[dict]:
    """Fetch ALL items with feed and folder metadata. No limit, no pagination."""
    rows = conn.execute(
        """
        SELECT i.id, i.title, i.link, i.content, i.date, i.status,
               i.feed_id, f.title as feed_title,
               COALESCE(fo.title, 'uncategorized') as folder_name
        FROM items i
        JOIN feeds f ON f.id = i.feed_id
        LEFT JOIN folders fo ON fo.id = f.folder_id
        ORDER BY i.id
        """
    ).fetchall()
    return [dict(r) for r in rows]


def get_feed_folder_map(conn: sqlite3.Connection) -> dict[int, dict]:
    """Build feed_id -> {folder, feed_name} mapping."""
    feeds = list_feeds(conn)
    return {
        f["id"]: {"folder": f["folder_name"], "feed_name": f["title"]}
        for f in feeds
    }
