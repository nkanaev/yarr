package storage

import (
	"database/sql"
	"time"
)

type FeedState struct {
	FeedID           int64
	LastRefreshed    time.Time
	LastError        *string
	HTTPLastModified string
	HTTPEtag         string
}

func (s *Storage) ListFeedStates() ([]FeedState, error) {
	rows, err := s.db.Query(`
		select feed_id, last_refreshed, last_modified, etag, last_error
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
			&state.HTTPLastModified,
			&state.HTTPEtag,
			&state.LastError,
		)
		if err != nil {
			return nil, err
		}
		states = append(states, state)
	}
	return states, nil
}

func (s *Storage) GetFeedState(feedID int64) (*FeedState, error) {
	var state FeedState
	err := s.db.QueryRow(`
		select feed_id, last_refreshed, last_modified, etag, last_error
		from feed_states where feed_id = :id
	`, sql.Named("id", feedID)).Scan(
		&state.FeedID,
		&state.LastRefreshed,
		&state.HTTPLastModified,
		&state.HTTPEtag,
		&state.LastError,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &state, nil
}

type UpdateFeedStateParams struct {
	LastRefreshed    *time.Time
	LastError        *string
	HTTPLastModified *string
	HTTPEtag         *string
}

func (s *Storage) UpdateFeedState(feedID int64, params UpdateFeedStateParams) (bool, error) {
	lastError := params.LastError
	if lastError != nil && *lastError == "" {
		lastError = nil
	}

	_, err := s.db.Exec(`
		insert into feed_states (
			feed_id
			, last_refreshed
			, last_modified
			, etag
			, last_error
		)
		values (
			:id
			, coalesce(:refreshed, 0)
			, coalesce(:last_modified, '')
			, coalesce(:etag, '')
			, coalesce(:last_error, '')
		)
		on conflict (feed_id) do update set
			last_refreshed = coalesce(:refreshed, last_refreshed),
			last_modified  = coalesce(:last_modified, last_modified),
			etag           = coalesce(:etag, etag),
			last_error     = coalesce(:last_error, last_error)
	`,
		sql.Named("id", feedID),
		sql.Named("refreshed", params.LastRefreshed),
		sql.Named("last_modified", params.HTTPLastModified),
		sql.Named("etag", params.HTTPEtag),
		sql.Named("last_error", params.LastError),
	)
	if err != nil {
		return false, err
	}
	return true, nil
}
