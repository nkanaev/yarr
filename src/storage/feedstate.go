package storage

import "time"

type FeedState struct {
	LastRefreshed time.Time
	LastError     string

	HTTPLastModified string
	HTTPEtag         string
}

func (s *Storage) ListFeedStates() ([]FeedState, error) {
	// TODO: implement
}

func (s *Storage) GetFeedState() (FeedState, error) {
	// TODO: implement
}

type UpdateFeedStateParams struct {
	LastRefreshed *time.Time
	LastError *string

	HTTPLastModified *string
	HTTPEtag *string
}

func (s *Storage) UpdateFeedState(params UpdateFeedStateParams) (bool, error) {
	// TODO: implement
}
