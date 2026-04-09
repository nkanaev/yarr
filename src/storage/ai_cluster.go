package storage

import (
	"database/sql"
	"strings"
)

// ClusterLabel represents a topic label with its article count.
type ClusterLabel struct {
	Label        string `json:"label"`
	ArticleCount int    `json:"article_count"`
}

// ClusterCentroid holds the centroid vector for a cluster, used for label stability.
type ClusterCentroid struct {
	ClusterID int    `json:"cluster_id"`
	Label     string `json:"label"`
	Centroid  []byte `json:"centroid"` // raw float64 bytes (numpy .tobytes())
}

// ClusterSummary is the response for GET /api/ai/clusters.
type ClusterSummary struct {
	GeneratedAt string         `json:"generated_at"`
	NClusters   int            `json:"n_clusters"`
	NArticles   int            `json:"n_articles"`
	Clusters    []ClusterLabel `json:"clusters"`
}

// ArticleTag maps an article URL to a cluster tag label.
type ArticleTag struct {
	URL string `json:"url"`
	Tag string `json:"tag"`
}

// ArticleResult is returned by GetArticlesByTag.
type ArticleResult struct {
	ID        int64      `json:"id"`
	URL       string     `json:"url"`
	Title     string     `json:"title"`
	Published string     `json:"published"`
	Folder    string     `json:"folder"`
	FeedName  string     `json:"feed_name"`
	Tag       string     `json:"tag"`
	Status    ItemStatus `json:"status"`
}

func m13_add_ai_article_tags(tx *sql.Tx) error {
	sql := `
		create table ai_article_tags (
			id  integer primary key autoincrement,
			url text    not null,
			tag text    not null
		);

		create index idx_ai_article_tags_tag on ai_article_tags(tag);
		create index idx_ai_article_tags_url on ai_article_tags(url);
	`
	_, err := tx.Exec(sql)
	return err
}

func m12_add_ai_cluster_tables(tx *sql.Tx) error {
	sql := `
		create table ai_cluster_runs (
			id           integer primary key autoincrement,
			generated_at text    not null,
			algorithm    text    not null default 'hdbscan',
			n_articles   integer not null default 0,
			n_noise      integer not null default 0
		);

		create table ai_cluster_labels (
			id            integer primary key autoincrement,
			run_id        integer not null references ai_cluster_runs(id) on delete cascade,
			label         text    not null,
			article_count integer not null default 0
		);

		create index idx_ai_cluster_labels_run on ai_cluster_labels(run_id);

		create table ai_cluster_centroids (
			id         integer primary key autoincrement,
			run_id     integer not null references ai_cluster_runs(id) on delete cascade,
			cluster_id integer not null,
			label      text    not null,
			centroid   blob    not null
		);

		create index idx_ai_cluster_centroids_run on ai_cluster_centroids(run_id);
	`
	_, err := tx.Exec(sql)
	return err
}

// SaveClusterRun inserts a new cluster run with its labels and centroids,
// deleting all previous runs (keeps only the latest).
func (s *Storage) SaveClusterRun(
	generatedAt string,
	algorithm string,
	nArticles int,
	nNoise int,
	labels []ClusterLabel,
	centroids []ClusterCentroid,
) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Insert run metadata
	res, err := tx.Exec(
		`insert into ai_cluster_runs (generated_at, algorithm, n_articles, n_noise) values (?, ?, ?, ?)`,
		generatedAt, algorithm, nArticles, nNoise,
	)
	if err != nil {
		return 0, err
	}
	runID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Batch-insert labels
	labelStmt, err := tx.Prepare(`insert into ai_cluster_labels (run_id, label, article_count) values (?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer labelStmt.Close()
	for _, l := range labels {
		if _, err := labelStmt.Exec(runID, l.Label, l.ArticleCount); err != nil {
			return 0, err
		}
	}

	// Batch-insert centroids
	centStmt, err := tx.Prepare(`insert into ai_cluster_centroids (run_id, cluster_id, label, centroid) values (?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer centStmt.Close()
	for _, c := range centroids {
		if _, err := centStmt.Exec(runID, c.ClusterID, c.Label, c.Centroid); err != nil {
			return 0, err
		}
	}

	// Delete all previous runs (cascade deletes labels + centroids)
	if _, err := tx.Exec(`delete from ai_cluster_runs where id != ?`, runID); err != nil {
		return 0, err
	}

	return runID, tx.Commit()
}

// GetClusterSummary returns the merged cluster summary from the latest run.
// Returns nil if no run exists.
func (s *Storage) GetClusterSummary(status int64, since string) (*ClusterSummary, error) {
	row := s.db.QueryRow(`
		select r.id, r.generated_at, r.n_articles
		from ai_cluster_runs r
		order by r.id desc
		limit 1
	`)

	var runID int64
	var generatedAt string
	var runArticles int
	if err := row.Scan(&runID, &generatedAt, &runArticles); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var (
		clusters      []ClusterLabel
		totalArticles int
	)

	if status < 0 && since == "" {
		// Fast path: default topics view uses precomputed cluster counts.
		rows, err := s.db.Query(`
			select label, sum(article_count) as article_count
			from ai_cluster_labels
			where run_id = ?
			group by label
			order by sum(article_count) desc
		`, runID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var l ClusterLabel
			if err := rows.Scan(&l.Label, &l.ArticleCount); err != nil {
				return nil, err
			}
			clusters = append(clusters, l)
			totalArticles += l.ArticleCount
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		if totalArticles == 0 {
			totalArticles = runArticles
		}
	} else {
		query := `
			select lbl.label, count(distinct at.url) as article_count
			from (
				select distinct label
				from ai_cluster_labels
				where run_id = ?
			) lbl
			join ai_article_tags at on at.tag = lbl.label
			join items i on i.link = at.url
		`
		args := []interface{}{runID}
		clauses := make([]string, 0, 2)
		if status >= 0 {
			clauses = append(clauses, "i.status = ?")
			args = append(args, status)
		}
		if since != "" {
			clauses = append(clauses, "i.date >= ?")
			args = append(args, since)
		}
		if len(clauses) > 0 {
			query += " where " + strings.Join(clauses, " and ")
		}
		query += " group by lbl.label order by article_count desc"

		rows, err := s.db.Query(query, args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var l ClusterLabel
			if err := rows.Scan(&l.Label, &l.ArticleCount); err != nil {
				return nil, err
			}
			clusters = append(clusters, l)
			totalArticles += l.ArticleCount
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return &ClusterSummary{
		GeneratedAt: generatedAt,
		NClusters:   len(clusters),
		NArticles:   totalArticles,
		Clusters:    clusters,
	}, nil
}

// GetClusterCentroids returns all centroids from the latest run for label stability.
// Returns nil if no run exists.
func (s *Storage) GetClusterCentroids() ([]ClusterCentroid, error) {
	row := s.db.QueryRow(`select id from ai_cluster_runs order by id desc limit 1`)
	var runID int64
	if err := row.Scan(&runID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	rows, err := s.db.Query(`
		select cluster_id, label, centroid
		from ai_cluster_centroids
		where run_id = ?
	`, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var centroids []ClusterCentroid
	for rows.Next() {
		var c ClusterCentroid
		if err := rows.Scan(&c.ClusterID, &c.Label, &c.Centroid); err != nil {
			return nil, err
		}
		centroids = append(centroids, c)
	}
	return centroids, rows.Err()
}

// SaveArticleTags replaces all article-tag mappings in a single transaction.
func (s *Storage) SaveArticleTags(tags []ArticleTag) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`delete from ai_article_tags`); err != nil {
		return err
	}

	stmt, err := tx.Prepare(`insert into ai_article_tags (url, tag) values (?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, t := range tags {
		if _, err := stmt.Exec(t.URL, t.Tag); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetArticlesByTag returns articles for a given tag label, joined with item metadata.
func (s *Storage) GetArticlesByTag(tag string, limit int, status int64, since string) ([]ArticleResult, error) {
	query := `
		select
			max(i.id),
			at.url,
			max(i.title),
			coalesce(max(i.date), '') as published,
			coalesce(max(fo.title), 'uncategorized') as folder,
			max(f.title) as feed_name,
			at.tag,
			max(i.status) as status
		from ai_article_tags at
		join items i on i.link = at.url
		join feeds f on f.id = i.feed_id
		left join folders fo on fo.id = f.folder_id
		where at.tag = ?`
	args := []interface{}{tag}
	if status >= 0 {
		query += " and i.status = ?"
		args = append(args, status)
	}
	if since != "" {
		query += " and i.date >= ?"
		args = append(args, since)
	}
	query += `
		group by at.url
		order by published desc
		limit ?
	`
	args = append(args, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ArticleResult
	for rows.Next() {
		var a ArticleResult
		if err := rows.Scan(&a.ID, &a.URL, &a.Title, &a.Published, &a.Folder, &a.FeedName, &a.Tag, &a.Status); err != nil {
			return nil, err
		}
		results = append(results, a)
	}
	return results, rows.Err()
}

// UpsertArticleTags inserts or replaces article-tag mappings for the given URLs
// without touching rows for URLs not in the payload (append-only semantics).
// Each URL in the payload has its existing row(s) deleted before re-inserting,
// so calling this twice for the same URL is idempotent.
func (s *Storage) UpsertArticleTags(tags []ArticleTag) error {
	if len(tags) == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	delStmt, err := tx.Prepare(`delete from ai_article_tags where url = ?`)
	if err != nil {
		return err
	}
	defer delStmt.Close()

	insStmt, err := tx.Prepare(`insert into ai_article_tags (url, tag) values (?, ?)`)
	if err != nil {
		return err
	}
	defer insStmt.Close()

	for _, t := range tags {
		if _, err := delStmt.Exec(t.URL); err != nil {
			return err
		}
		if _, err := insStmt.Exec(t.URL, t.Tag); err != nil {
			return err
		}
	}

	return tx.Commit()
}
