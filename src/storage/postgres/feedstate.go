package postgres

import (
	"database/sql"

	"github.com/nkanaev/yarr/src/storage/model"
)

func (s *PostgresStorage) ListFeedStates() ([]model.FeedState, error) {
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

	states := make([]model.FeedState, 0)
	for rows.Next() {
		var state model.FeedState
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

func (s *PostgresStorage) GetFeedState(feedID int64) (*model.FeedState, error) {
	var state model.FeedState
	err := s.db.QueryRow(`
		select
			feed_id
			, last_refreshed
			, last_error
			, http_lmod
			, http_etag
		from feed_states where feed_id = $1
	`, feedID).Scan(
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

func (s *PostgresStorage) UpdateFeedState(feedID int64, params model.UpdateFeedStateParams) (bool, error) {
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
			$1
			, coalesce($2, '1970-01-01 00:00:00+00'::timestamptz)
			, coalesce($3, '')
			, coalesce($4, '')
			, coalesce($5, '')
		)
		on conflict (feed_id) do update set
			last_refreshed = coalesce($2, feed_states.last_refreshed),
			last_error     = coalesce($3, feed_states.last_error),
			http_lmod      = coalesce($4, feed_states.http_lmod),
			http_etag      = coalesce($5, feed_states.http_etag)
	`,
		feedID,
		params.LastRefreshed,
		params.LastError,
		params.HTTPLastModified,
		params.HTTPEtag,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}
