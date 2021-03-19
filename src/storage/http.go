package storage

import (
	"log"
	"time"
)

type HTTPState struct {
	FeedID        int64
	LastRefreshed time.Time

	LastModified string
	Etag         string
}

func (s *Storage) ListHTTPStates() map[int64]HTTPState {
	result := make(map[int64]HTTPState)
	rows, err := s.db.Query(`select feed_id, last_refreshed, last_modified, etag from http_states`)
	if err != nil {
		log.Print(err)
		return result
	}
	for rows.Next() {
		var state HTTPState
		err = rows.Scan(
			&state.FeedID,
			&state.LastRefreshed,
			&state.LastModified,
			&state.Etag,
		)
		if err != nil {
			log.Print(err)
			return result
		}
		result[state.FeedID] = state
	}
	return result
}

func (s *Storage) GetHTTPState(feedID int64) *HTTPState {
	row := s.db.QueryRow(`
		select feed_id, last_refreshed, last_modified, etag
		from http_states where feed_id = ?
	`, feedID)

	if row == nil {
		return nil
	}

	var state HTTPState
	row.Scan(
		&state.FeedID,
		&state.LastRefreshed,
		&state.LastModified,
		&state.Etag,
	)
	return &state
}

func (s *Storage) SetHTTPState(feedID int64, lastModified, etag string) {
	_, err := s.db.Exec(`
		insert into http_states (feed_id, last_modified, etag, last_refreshed)
		values (?, ?, ?, datetime())
		on conflict (feed_id) do update set last_modified = ?, etag = ?, last_refreshed = datetime()`,
		// insert
		feedID, lastModified, etag,
		// upsert
		lastModified, etag,
	)
	if err != nil {
		log.Print(err)
	}
}
