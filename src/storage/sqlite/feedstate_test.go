package storage

import (
	"testing"
	"time"
)

func TestUpdateFeedState_Full(t *testing.T) {
	s := testDB()
	defer s.Close()

	f := s.CreateFeed(CreateFeedParams{Title: "Test", FeedLink: "http://example.com"})

	now := time.Now().UTC().Truncate(time.Second)
	errMsg := "error"
	lmod := "today"
	etag := "v1"

	ok, err := s.UpdateFeedState(f.Id, UpdateFeedStateParams{
		LastRefreshed:    &now,
		LastError:        &errMsg,
		HTTPLastModified: &lmod,
		HTTPEtag:         &etag,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("expected true")
	}

	state, err := s.GetFeedState(f.Id)
	if err != nil {
		t.Fatal(err)
	}
	if state == nil {
		t.Fatal("expected state, got nil")
	}
	if !state.LastRefreshed.Equal(now) {
		t.Errorf("expected %v, got %v", now, state.LastRefreshed)
	}
	if state.LastError != errMsg {
		t.Errorf("expected %s, got %v", errMsg, state.LastError)
	}
	if state.HTTPLastModified != lmod {
		t.Errorf("expected %s, got %s", lmod, state.HTTPLastModified)
	}
	if state.HTTPEtag != etag {
		t.Errorf("expected %s, got %s", etag, state.HTTPEtag)
	}
}

func TestUpdateFeedState_Partial(t *testing.T) {
	s := testDB()
	defer s.Close()

	f := s.CreateFeed(CreateFeedParams{Title: "Test", FeedLink: "http://example.com"})
	etag := "v1"
	s.UpdateFeedState(f.Id, UpdateFeedStateParams{HTTPEtag: &etag})

	newErr := "new error"
	_, err := s.UpdateFeedState(f.Id, UpdateFeedStateParams{
		LastError: &newErr,
	})
	if err != nil {
		t.Fatal(err)
	}

	state, err := s.GetFeedState(f.Id)
	if err != nil {
		t.Fatal(err)
	}
	if state.LastError != newErr {
		t.Errorf("expected %s, got %v", newErr, state.LastError)
	}
	if state.HTTPEtag != etag {
		t.Errorf("etag should be unchanged, got %s", state.HTTPEtag)
	}
}

func TestUpdateFeedState_ClearError(t *testing.T) {
	s := testDB()
	defer s.Close()

	f := s.CreateFeed(CreateFeedParams{Title: "Test", FeedLink: "http://example.com"})
	errMsg := "error"
	s.UpdateFeedState(f.Id, UpdateFeedStateParams{LastError: &errMsg})

	empty := ""
	_, err := s.UpdateFeedState(f.Id, UpdateFeedStateParams{
		LastError: &empty,
	})
	if err != nil {
		t.Fatal(err)
	}

	state, err := s.GetFeedState(f.Id)
	if err != nil {
		t.Fatal(err)
	}
	if state.LastError != "" {
		t.Errorf("expected empty error string, got %v", state.LastError)
	}
}

func TestListFeedStates(t *testing.T) {
	s := testDB()
	defer s.Close()

	f1 := s.CreateFeed(CreateFeedParams{Title: "F1", FeedLink: "L1"})
	f2 := s.CreateFeed(CreateFeedParams{Title: "F2", FeedLink: "L2"})

	errMsg := "fail"
	s.UpdateFeedState(f1.Id, UpdateFeedStateParams{LastError: &errMsg})
	s.UpdateFeedState(f2.Id, UpdateFeedStateParams{HTTPEtag: ptr("e")})

	states, err := s.ListFeedStates()
	if err != nil {
		t.Fatal(err)
	}

	if len(states) != 2 {
		t.Errorf("expected 2 states, got %d", len(states))
	}
}

func ptr[T any](v T) *T {
	return &v
}
