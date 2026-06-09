package sqlite

import (
	"database/sql"
	"time"
)

func (s *SQLiteStorage) ListFeedStates() ([]FeedState, error) {
	rows, err := s.db.Query(`
		select
			feed_id
			, last_refreshed
			, last_error
			, http_lmod
			, http_etag
		from feed_states
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	states := make([]FeedState, 0)
	for rows.Next() {
		var state FeedState
		err := rows.Scan(
			&state.FeedID,
			&state.LastRefreshed,
			&state.LastError,
			&state.HTTPLastModified,
			&state.HTTPEtag,
		)
		if err != nil {
			return nil, err
		}
		states = append(states, state)
	}
	return states, nil
}

func (s *SQLiteStorage) GetFeedState(feedID int64) (*FeedState, error) {
	var state FeedState
	err := s.db.QueryRow(`
		select
			feed_id
			, last_refreshed
			, last_error
			, http_lmod
			, http_etag
		from feed_states where feed_id = :id
	`, sql.Named("id", feedID)).Scan(
		&state.FeedID,
		&state.LastRefreshed,
		&state.LastError,
		&state.HTTPLastModified,
		&state.HTTPEtag,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (s *SQLiteStorage) UpdateFeedState(feedID int64, params UpdateFeedStateParams) (bool, error) {
	lastError := params.LastError
	if lastError != nil && *lastError == "" {
		lastError = nil
	}

	_, err := s.db.Exec(`
		insert into feed_states (
			feed_id
			, last_refreshed
			, last_error
			, http_lmod
			, http_etag
		)
		values (
			:id
			, coalesce(:last_refreshed, 0)
			, coalesce(:last_error, '')
			, coalesce(:http_lmod, '')
			, coalesce(:http_etag, '')
		)
		on conflict (feed_id) do update set
			last_refreshed = coalesce(:last_refreshed, last_refreshed),
			last_error     = coalesce(:last_error, last_error),
			http_lmod      = coalesce(:http_lmod, http_lmod),
			http_etag      = coalesce(:http_etag, http_etag)
	`,
		sql.Named("id", feedID),
		sql.Named("last_refreshed", params.LastRefreshed),
		sql.Named("last_error", params.LastError),
		sql.Named("http_lmod", params.HTTPLastModified),
		sql.Named("http_etag", params.HTTPEtag),
	)
	if err != nil {
		return false, err
	}
	return true, nil
}
