package storage

import (
	"database/sql"
	"log"
	"math"
	"strings"
	"time"
	"unicode"
)

// ── Types ──────────────────────────────────────────────────────────────────────

// RankedItem extends Item with ranking metadata.
type RankedItem struct {
	Id              int64      `json:"id"`
	GUID            string     `json:"guid"`
	FeedId          int64      `json:"feed_id"`
	Title           string     `json:"title"`
	Link            string     `json:"link"`
	Date            time.Time  `json:"date"`
	Status          ItemStatus `json:"status"`
	MediaLinks      MediaLinks `json:"media_links"`
	Score           float64    `json:"score"`
	RelevanceReason string     `json:"relevance_reason"`
	Reaction        string     `json:"reaction"`
}

// FeedAffinity represents a feed's affinity score for the user.
type FeedAffinity struct {
	FeedID int64   `json:"feed_id"`
	Title  string  `json:"title"`
	Score  float64 `json:"score"`
}

// TopicAffinity represents a topic's affinity score for the user.
type TopicAffinity struct {
	Topic string  `json:"topic"`
	Score float64 `json:"score"`
}

// PreferenceStats holds aggregate preference data for the settings UI.
type PreferenceStats struct {
	TotalLikes     int             `json:"total_likes"`
	TotalDislikes  int             `json:"total_dislikes"`
	TotalClicks    int             `json:"total_clicks"`
	TotalReadHeres int             `json:"total_read_heres"`
	TopFeeds       []FeedAffinity  `json:"top_feeds"`
	TopTopics      []TopicAffinity `json:"top_topics"`
}

// ── Migration ──────────────────────────────────────────────────────────────────

func m14_add_ranking_tables(tx *sql.Tx) error {
	sql := `
		CREATE TABLE reactions (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			item_id    INTEGER NOT NULL UNIQUE REFERENCES items(id) ON DELETE CASCADE,
			reaction   TEXT    NOT NULL CHECK(reaction IN ('like','dislike')),
			created_at TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f','now'))
		);
		CREATE INDEX idx_reactions_created_at ON reactions(created_at);

		CREATE TABLE click_throughs (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			item_id    INTEGER NOT NULL REFERENCES items(id) ON DELETE CASCADE,
			created_at TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f','now'))
		);
		CREATE INDEX idx_click_throughs_item ON click_throughs(item_id);

		CREATE TABLE keyword_index (
			id      INTEGER PRIMARY KEY AUTOINCREMENT,
			item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE CASCADE,
			keyword TEXT    NOT NULL
		);
		CREATE INDEX idx_keyword_index_keyword ON keyword_index(keyword);
		CREATE INDEX idx_keyword_index_item    ON keyword_index(item_id);
	`
	_, err := tx.Exec(sql)
	return err
}

// ── Keyword Extraction ─────────────────────────────────────────────────────────

var stopWords = map[string]bool{
	"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
	"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
	"with": true, "by": true, "from": true, "as": true, "is": true, "was": true,
	"are": true, "were": true, "been": true, "be": true, "have": true, "has": true,
	"had": true, "do": true, "does": true, "did": true, "will": true, "would": true,
	"could": true, "should": true, "may": true, "might": true, "can": true,
	"this": true, "that": true, "these": true, "those": true, "its": true,
	"you": true, "he": true, "she": true, "it": true, "we": true, "they": true,
	"what": true, "which": true, "who": true, "whom": true, "whose": true,
	"where": true, "when": true, "why": true, "how": true, "all": true,
	"each": true, "every": true, "both": true, "few": true, "more": true,
	"most": true, "other": true, "some": true, "such": true, "no": true,
	"nor": true, "not": true, "only": true, "own": true, "same": true,
	"so": true, "than": true, "too": true, "very": true, "just": true,
	"new": true, "about": true, "into": true, "over": true, "after": true,
	"also": true, "back": true, "use": true, "two": true, "way": true,
	"your": true, "our": true, "out": true, "up": true, "one": true,
}

func extractKeywords(title string) []string {
	title = strings.ToLower(title)
	var result []string
	var word []rune
	for _, r := range title {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			word = append(word, r)
		} else {
			if len(word) > 2 {
				w := string(word)
				if !stopWords[w] {
					result = append(result, w)
				}
			}
			word = word[:0]
		}
	}
	if len(word) > 2 {
		w := string(word)
		if !stopWords[w] {
			result = append(result, w)
		}
	}
	return result
}

// IndexItemKeywords extracts keywords from title and inserts into keyword_index.
func (s *Storage) IndexItemKeywords(itemID int64, title string) error {
	keywords := extractKeywords(title)
	if len(keywords) == 0 {
		return nil
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO keyword_index (item_id, keyword) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, kw := range keywords {
		if _, err := stmt.Exec(itemID, kw); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// IndexNewItemKeywords indexes keywords for recently created items by looking up
// their IDs via feed_id + guid, skipping items that already have keywords indexed.
func (s *Storage) IndexNewItemKeywords(items []Item) {
	for _, item := range items {
		var id int64
		err := s.db.QueryRow(
			`SELECT id FROM items WHERE feed_id = ? AND guid = ?`,
			item.FeedId, item.GUID,
		).Scan(&id)
		if err != nil {
			continue
		}
		// Skip if already indexed
		var count int
		s.db.QueryRow(`SELECT COUNT(*) FROM keyword_index WHERE item_id = ?`, id).Scan(&count)
		if count > 0 {
			continue
		}
		if err := s.IndexItemKeywords(id, item.Title); err != nil {
			log.Printf("IndexNewItemKeywords error for item %d: %v", id, err)
		}
	}
}

// BackfillKeywordIndex populates keyword_index for all items that have no entries.
func (s *Storage) BackfillKeywordIndex() error {
	rows, err := s.db.Query(`
		SELECT i.id, i.title FROM items i
		WHERE NOT EXISTS (SELECT 1 FROM keyword_index ki WHERE ki.item_id = i.id)
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	type entry struct {
		id    int64
		title string
	}
	var batch []entry
	for rows.Next() {
		var e entry
		if err := rows.Scan(&e.id, &e.title); err != nil {
			return err
		}
		batch = append(batch, e)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if len(batch) == 0 {
		return nil
	}

	log.Printf("Backfilling keyword_index for %d items", len(batch))

	const batchSize = 1000
	for i := 0; i < len(batch); i += batchSize {
		end := i + batchSize
		if end > len(batch) {
			end = len(batch)
		}
		tx, err := s.db.Begin()
		if err != nil {
			return err
		}
		stmt, err := tx.Prepare(`INSERT INTO keyword_index (item_id, keyword) VALUES (?, ?)`)
		if err != nil {
			tx.Rollback()
			return err
		}
		for _, e := range batch[i:end] {
			for _, kw := range extractKeywords(e.title) {
				if _, err := stmt.Exec(e.id, kw); err != nil {
					stmt.Close()
					tx.Rollback()
					return err
				}
			}
		}
		stmt.Close()
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	log.Printf("Keyword backfill complete")
	return nil
}

// EnsureKeywordIndex checks if keyword_index is empty and backfills if needed.
func (s *Storage) EnsureKeywordIndex() {
	var count int64
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM keyword_index`).Scan(&count); err != nil {
		log.Print("EnsureKeywordIndex count error:", err)
		return
	}
	if count == 0 {
		if err := s.BackfillKeywordIndex(); err != nil {
			log.Print("EnsureKeywordIndex backfill error:", err)
		}
	}
}

// ── Reaction CRUD ──────────────────────────────────────────────────────────────

// SetReaction sets or clears a reaction for an item.
// If reaction is empty, the reaction is deleted.
func (s *Storage) SetReaction(itemID int64, reaction string) error {
	if reaction == "" {
		_, err := s.db.Exec(`DELETE FROM reactions WHERE item_id = ?`, itemID)
		return err
	}
	_, err := s.db.Exec(`
		INSERT INTO reactions (item_id, reaction) VALUES (?, ?)
		ON CONFLICT(item_id) DO UPDATE SET reaction = excluded.reaction, created_at = strftime('%Y-%m-%d %H:%M:%f','now')
	`, itemID, reaction)
	return err
}

// GetReaction returns the reaction for a single item ("like", "dislike", or "").
func (s *Storage) GetReaction(itemID int64) (string, error) {
	var reaction string
	err := s.db.QueryRow(`SELECT reaction FROM reactions WHERE item_id = ?`, itemID).Scan(&reaction)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return reaction, err
}

// LogClickThrough records a click-through event for an item.
func (s *Storage) LogClickThrough(itemID int64) error {
	_, err := s.db.Exec(`INSERT INTO click_throughs (item_id) VALUES (?)`, itemID)
	return err
}

// LogReadHere records or refreshes a read-here event for an item.
// Upserts so that re-reading the same article updates the timestamp (fresher = higher decay weight).
func (s *Storage) LogReadHere(itemID int64) error {
	_, err := s.db.Exec(`
		INSERT INTO read_heres (item_id) VALUES (?)
		ON CONFLICT(item_id) DO UPDATE SET created_at = strftime('%Y-%m-%d %H:%M:%f','now')
	`, itemID)
	return err
}

// DeleteAllPreferences clears all reactions, click-throughs, and read-heres.
func (s *Storage) DeleteAllPreferences() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`DELETE FROM reactions`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM click_throughs`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM read_heres`); err != nil {
		return err
	}
	return tx.Commit()
}

// ── Preference Stats ───────────────────────────────────────────────────────────

// GetPreferenceStats returns aggregate preference data.
func (s *Storage) GetPreferenceStats() (PreferenceStats, error) {
	var stats PreferenceStats

	// Counts
	s.db.QueryRow(`SELECT COUNT(*) FROM reactions WHERE reaction = 'like'`).Scan(&stats.TotalLikes)
	s.db.QueryRow(`SELECT COUNT(*) FROM reactions WHERE reaction = 'dislike'`).Scan(&stats.TotalDislikes)
	s.db.QueryRow(`SELECT COUNT(*) FROM click_throughs`).Scan(&stats.TotalClicks)
	s.db.QueryRow(`SELECT COUNT(*) FROM read_heres`).Scan(&stats.TotalReadHeres)

	// Top feeds by affinity
	feedRows, err := s.db.Query(`
		SELECT i.feed_id, f.title,
			SUM(CASE r.reaction WHEN 'like' THEN 1.0 WHEN 'dislike' THEN -1.0 END) as score
		FROM reactions r
		JOIN items i ON i.id = r.item_id
		JOIN feeds f ON f.id = i.feed_id
		GROUP BY i.feed_id
		ORDER BY score DESC
		LIMIT 5
	`)
	if err == nil {
		defer feedRows.Close()
		for feedRows.Next() {
			var fa FeedAffinity
			if err := feedRows.Scan(&fa.FeedID, &fa.Title, &fa.Score); err == nil {
				stats.TopFeeds = append(stats.TopFeeds, fa)
			}
		}
	}

	// Top topics by affinity
	topicRows, err := s.db.Query(`
		SELECT at.tag,
			SUM(CASE r.reaction WHEN 'like' THEN 1.0 WHEN 'dislike' THEN -1.0 END) as score
		FROM reactions r
		JOIN items i ON i.id = r.item_id
		JOIN ai_article_tags at ON at.url = i.link
		GROUP BY at.tag
		ORDER BY score DESC
		LIMIT 5
	`)
	if err == nil {
		defer topicRows.Close()
		for topicRows.Next() {
			var ta TopicAffinity
			if err := topicRows.Scan(&ta.Topic, &ta.Score); err == nil {
				stats.TopTopics = append(stats.TopTopics, ta)
			}
		}
	}

	if stats.TopFeeds == nil {
		stats.TopFeeds = []FeedAffinity{}
	}
	if stats.TopTopics == nil {
		stats.TopTopics = []TopicAffinity{}
	}

	return stats, nil
}

// ── Ranked Items ───────────────────────────────────────────────────────────────

type rankedCandidate struct {
	item     RankedItem
	keywords []string
}

const decayRate = 0.015 // ~10% per week

func decayWeight(ageDays float64) float64 {
	return math.Exp(-decayRate * ageDays)
}

// GetRankedItems returns items scored by the personalized ranking algorithm.
// Items are fetched from the last 30 days (or all if fewer), scored in Go, sorted, and paginated.
func (s *Storage) GetRankedItems(filter ItemFilter, limit int, offset int) ([]RankedItem, bool, error) {
	now := time.Now().UTC()

	// Step 1: Compute feed affinity scores from reactions + stars
	feedAffinity := make(map[int64]float64)
	feedReactionCount := make(map[int64]int)
	{
		rows, err := s.db.Query(`
			SELECT i.feed_id, r.reaction, r.created_at
			FROM reactions r
			JOIN items i ON i.id = r.item_id
		`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var feedID int64
				var reaction string
				var createdAt time.Time
				if err := rows.Scan(&feedID, &reaction, &createdAt); err != nil {
					log.Print("GetRankedItems step1 reactions scan:", err)
					continue
				}
				age := now.Sub(createdAt).Hours() / 24
				w := decayWeight(age)
				if reaction == "like" {
					feedAffinity[feedID] += w
				} else {
					feedAffinity[feedID] -= w
				}
				feedReactionCount[feedID]++
			}
		}
		// Stars as implicit strong likes (2x weight)
		starRows, err := s.db.Query(`
			SELECT feed_id, date FROM items WHERE status = ?
		`, STARRED)
		if err == nil {
			defer starRows.Close()
			for starRows.Next() {
				var feedID int64
				var date time.Time
				if err := starRows.Scan(&feedID, &date); err != nil {
					log.Print("GetRankedItems step1 stars scan:", err)
					continue
				}
				age := now.Sub(date).Hours() / 24
				feedAffinity[feedID] += 2.0 * decayWeight(age)
				feedReactionCount[feedID]++
			}
		}
	}

	// Step 2: Compute topic affinity scores
	topicAffinity := make(map[string]float64)
	topicReactionCount := make(map[string]int)
	{
		rows, err := s.db.Query(`
			SELECT at.tag, r.reaction, r.created_at
			FROM reactions r
			JOIN items i ON i.id = r.item_id
			JOIN ai_article_tags at ON at.url = i.link
		`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var tag, reaction string
				var createdAt time.Time
				if err := rows.Scan(&tag, &reaction, &createdAt); err != nil {
					log.Print("GetRankedItems step2 scan:", err)
					continue
				}
				age := now.Sub(createdAt).Hours() / 24
				w := decayWeight(age)
				if reaction == "like" {
					topicAffinity[tag] += w
				} else {
					topicAffinity[tag] -= w
				}
				topicReactionCount[tag]++
			}
		}
	}

	// Step 3: Compute keyword affinity scores
	keywordAffinity := make(map[string]float64)
	keywordCount := make(map[string]int)
	{
		rows, err := s.db.Query(`
			SELECT ki.keyword, r.reaction, r.created_at
			FROM reactions r
			JOIN keyword_index ki ON ki.item_id = r.item_id
		`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var keyword, reaction string
				var createdAt time.Time
				if err := rows.Scan(&keyword, &reaction, &createdAt); err != nil {
					log.Print("GetRankedItems step3 scan:", err)
					continue
				}
				age := now.Sub(createdAt).Hours() / 24
				w := decayWeight(age)
				if reaction == "like" {
					keywordAffinity[keyword] += w
				} else {
					keywordAffinity[keyword] -= w
				}
				keywordCount[keyword]++
			}
		}
	}

	// Step 4: Compute click-through feed boost
	clickFeedBoost := make(map[int64]float64)
	{
		rows, err := s.db.Query(`
			SELECT i.feed_id, ct.created_at
			FROM click_throughs ct
			JOIN items i ON i.id = ct.item_id
		`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var feedID int64
				var createdAt time.Time
				if err := rows.Scan(&feedID, &createdAt); err != nil {
					log.Print("GetRankedItems step4 scan:", err)
					continue
				}
				age := now.Sub(createdAt).Hours() / 24
				clickFeedBoost[feedID] += 0.3 * decayWeight(age)
			}
		}
	}

	// Step 4b: Compute read-here feed boost (0.5x weight — stronger than click-through)
	readHereFeedBoost := make(map[int64]float64)
	{
		rows, err := s.db.Query(`
			SELECT i.feed_id, rh.created_at
			FROM read_heres rh
			JOIN items i ON i.id = rh.item_id
		`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var feedID int64
				var createdAt time.Time
				if err := rows.Scan(&feedID, &createdAt); err != nil {
					log.Print("GetRankedItems step4b scan:", err)
					continue
				}
				age := now.Sub(createdAt).Hours() / 24
				readHereFeedBoost[feedID] += 0.5 * decayWeight(age)
			}
		}
	}

	// Step 5: Build item-to-topic mapping
	itemTopics := make(map[int64][]string)
	{
		rows, err := s.db.Query(`
			SELECT i.id, at.tag
			FROM ai_article_tags at
			JOIN items i ON i.link = at.url
		`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var id int64
				var tag string
				if err := rows.Scan(&id, &tag); err != nil {
					continue
				}
				itemTopics[id] = append(itemTopics[id], tag)
			}
		}
	}

	// Step 6: Build item-to-reaction mapping
	itemReactions := make(map[int64]string)
	{
		rows, err := s.db.Query(`SELECT item_id, reaction FROM reactions`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var id int64
				var r string
				if err := rows.Scan(&id, &r); err != nil {
					continue
				}
				itemReactions[id] = r
			}
		}
	}

	// Step 7: Fetch candidate items (recent, with filters)
	cond := make([]string, 0)
	args := make([]interface{}, 0)
	if filter.FolderID != nil {
		cond = append(cond, "i.feed_id IN (SELECT id FROM feeds WHERE folder_id = ?)")
		args = append(args, *filter.FolderID)
	}
	if filter.FeedID != nil {
		cond = append(cond, "i.feed_id = ?")
		args = append(args, *filter.FeedID)
	}
	if filter.Status != nil {
		cond = append(cond, "i.status = ?")
		args = append(args, *filter.Status)
	} else {
		// By default, exclude read articles — they've already been consumed.
		// Starred articles (status=2) are kept: they're high-value and users may re-read them.
		cond = append(cond, "i.status != ?")
		args = append(args, READ)
	}

	where := ""
	if len(cond) > 0 {
		where = "WHERE " + strings.Join(cond, " AND ")
	}

	// Fetch a generous candidate pool to rank from
	query := `
		SELECT i.id, i.guid, i.feed_id, i.title, i.link, i.date, i.status, i.media_links
		FROM items i
		` + where + `
		ORDER BY i.date DESC
		LIMIT 1000
	`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var candidates []rankedCandidate
	for rows.Next() {
		var ri RankedItem
		if err := rows.Scan(&ri.Id, &ri.GUID, &ri.FeedId, &ri.Title, &ri.Link, &ri.Date, &ri.Status, &ri.MediaLinks); err != nil {
			log.Print("GetRankedItems step7 scan:", err)
			continue
		}
		ri.Reaction = itemReactions[ri.Id]
		kw := extractKeywords(ri.Title)
		candidates = append(candidates, rankedCandidate{item: ri, keywords: kw})
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	// Step 8: Score each candidate
	for i := range candidates {
		c := &candidates[i]
		item := &c.item
		ageDays := now.Sub(item.Date).Hours() / 24

		// Feed affinity (0.25 weight) — reactions + click-throughs + read-heres
		fa := feedAffinity[item.FeedId] + clickFeedBoost[item.FeedId] + readHereFeedBoost[item.FeedId]
		feedScore := clamp(fa*10, -30, 30) * 0.25

		// Topic affinity (0.25 weight)
		var topicScore float64
		topics := itemTopics[item.Id]
		if len(topics) > 0 {
			var sum float64
			for _, t := range topics {
				sum += topicAffinity[t]
			}
			avg := sum / float64(len(topics))
			topicScore = clamp(avg*10, -30, 30) * 0.25
		}

		// Keyword affinity (0.15 weight)
		var kwScore float64
		if len(c.keywords) > 0 {
			var kwSum float64
			var kwN int
			for _, kw := range c.keywords {
				if _, ok := keywordAffinity[kw]; ok {
					// Only count keywords with >= 5 total interactions
					if keywordCount[kw] >= 5 {
						kwSum += keywordAffinity[kw]
						kwN++
					}
				}
			}
			if kwN > 0 {
				avg := kwSum / float64(kwN)
				kwScore = clamp(avg*10, -20, 20) * 0.15
			}
		}

		// Recency (0.20 weight) — linear decay: 20 points at 0 days, 0 at 20+ days
		recencyScore := math.Max(0, 20.0-ageDays) * 0.20

		// Exploration bonus (0.15 weight) — boost under-explored feeds/topics
		var explorationScore float64
		if feedReactionCount[item.FeedId] < 5 {
			explorationScore = 15.0 * 0.15
		} else if len(topics) > 0 {
			underExplored := false
			for _, t := range topics {
				if topicReactionCount[t] < 5 {
					underExplored = true
					break
				}
			}
			if underExplored {
				explorationScore = 10.0 * 0.15
			}
		}

		// Base score of 50 + all components
		score := 50.0 + feedScore + topicScore + kwScore + recencyScore + explorationScore
		item.Score = clamp(score, 0, 100)

		// Generate relevance reason
		item.RelevanceReason = generateRelevanceReason(fa, topicAffinity, topics, keywordAffinity, c.keywords, feedReactionCount[item.FeedId])
	}

	// Step 9: Sort by score DESC, date DESC
	sortCandidates(candidates)

	// Step 10: Paginate
	total := len(candidates)
	start := offset
	if start > total {
		start = total
	}
	end := start + limit + 1 // fetch one extra to check has_more
	if end > total {
		end = total
	}

	result := make([]RankedItem, 0, limit)
	hasMore := false
	for idx, c := range candidates[start:end] {
		if idx >= limit {
			hasMore = true
			break
		}
		result = append(result, c.item)
	}

	return result, hasMore, nil
}

func clamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func sortCandidates(candidates []rankedCandidate) {
	// Simple insertion sort — stable and fine for <= 1000 items
	for i := 1; i < len(candidates); i++ {
		for j := i; j > 0; j-- {
			a, b := candidates[j], candidates[j-1]
			if a.item.Score > b.item.Score || (a.item.Score == b.item.Score && a.item.Date.After(b.item.Date)) {
				candidates[j], candidates[j-1] = b, a
			} else {
				break
			}
		}
	}
}

func generateRelevanceReason(feedAff float64, topicAff map[string]float64, topics []string, kwAff map[string]float64, keywords []string, feedReactionCount int) string {
	// Feed signal
	if feedAff > 2.0 {
		return "From a source you like"
	}

	// Topic signal
	for _, t := range topics {
		if topicAff[t] > 2.0 {
			return "Matches your interests"
		}
	}

	// Keyword signal
	var likedKW []string
	for _, kw := range keywords {
		if kwAff[kw] > 1.0 {
			likedKW = append(likedKW, kw)
		}
	}
	if len(likedKW) > 0 {
		if len(likedKW) > 2 {
			likedKW = likedKW[:2]
		}
		return "Related to: " + strings.Join(likedKW, ", ")
	}

	// Exploration
	if feedReactionCount < 5 {
		return "Discover new content"
	}

	return ""
}
