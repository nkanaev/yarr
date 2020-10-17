package storage

type HTTPState struct {
	LastModified string
	Etag         string
}

func (s *Storage) GetHTTPState(url string) *HTTPState {
	row := s.db.QueryRow(`
		select last_modified, etag
		from http_state where url = ?
	`, url)

	if row == nil {
		return nil
	}

	var state HTTPState
	row.Scan(&state.LastModified, &state.Etag)
	return &state
}

func (s *Storage) SetHTTPState(url string, state HTTPState) {
	_, err := s.db.Exec(`
		insert into http_state (url, last_modified, etag)
		values (?, ?, ?)
		on conflict (url) do update set last_modified = ?, etag = ?`,
		url, state.LastModified, state.Etag,
		state.LastModified, state.Etag,
	)
	if err != nil {
		s.log.Print(err)
	}
}
