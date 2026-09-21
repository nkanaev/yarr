package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/storage/model"
)

type feverTestStorage struct {
	storage.Storage
	markedRead   *model.MarkFilter
	itemStatusID int64
	itemStatus   *model.ItemStatus
}

func (s *feverTestStorage) MarkItemsRead(filter model.MarkFilter) bool {
	s.markedRead = &filter
	return true
}

func (s *feverTestStorage) UpdateItem(id int64, params model.UpdateItemParams) bool {
	s.itemStatusID = id
	s.itemStatus = params.Status
	return true
}

type feverTest struct {
	storage *feverTestStorage
	handler http.Handler
	apiKey  string
}

func newFeverTest(t *testing.T) *feverTest {
	t.Helper()
	db := &feverTestStorage{}
	auth := NewLocalAuthProvider("u", "p", "")
	server := NewServer("127.0.0.1:8000")
	server.Storage = NewLocalStorage(db)
	server.Auth = auth
	return &feverTest{
		storage: db,
		handler: server.Handler(),
		apiKey:  auth.FeverAPIKey(nil),
	}
}

func (f *feverTest) mark(form url.Values) *httptest.ResponseRecorder {
	form.Set("api_key", f.apiKey)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/fever/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	f.handler.ServeHTTP(rec, req)
	return rec
}

func TestFeverMarkFeedInvalidAs(t *testing.T) {
	f := newFeverTest(t)

	rec := f.mark(url.Values{"mark": {"feed"}, "as": {"unread"}, "id": {"1"}})

	if rec.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Result().StatusCode)
	}
	if f.storage.markedRead != nil {
		t.Fatalf("rejected request had a side effect (filter=%#v)", *f.storage.markedRead)
	}
}

func TestFeverMarkGroupInvalidAs(t *testing.T) {
	f := newFeverTest(t)

	rec := f.mark(url.Values{"mark": {"group"}, "as": {"unread"}, "id": {"0"}})

	if rec.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Result().StatusCode)
	}
	if f.storage.markedRead != nil {
		t.Fatalf("rejected request had a side effect (filter=%#v)", *f.storage.markedRead)
	}
}

func TestFeverMarkInvalidAsValues(t *testing.T) {
	for _, mark := range []string{"feed", "group"} {
		for _, as := range []string{"unread", "saved", "unsaved", ""} {
			t.Run(mark+"/"+as, func(t *testing.T) {
				f := newFeverTest(t)

				rec := f.mark(url.Values{"mark": {mark}, "as": {as}, "id": {"1"}})

				if rec.Result().StatusCode != http.StatusBadRequest {
					t.Fatalf("expected 400, got %d", rec.Result().StatusCode)
				}
				if f.storage.markedRead != nil {
					t.Fatalf("rejected request had a side effect (filter=%#v)", *f.storage.markedRead)
				}
			})
		}
	}
}

func TestFeverMarkFeedRead(t *testing.T) {
	f := newFeverTest(t)

	rec := f.mark(url.Values{"mark": {"feed"}, "as": {"read"}, "id": {"7"}})

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Result().StatusCode)
	}
	if f.storage.markedRead == nil || f.storage.markedRead.FeedID == nil {
		t.Fatal("expected MarkItemsRead with a FeedID")
	}
	if *f.storage.markedRead.FeedID != 7 {
		t.Fatalf("expected FeedID 7, got %d", *f.storage.markedRead.FeedID)
	}
}

func TestFeverMarkGroupRead(t *testing.T) {
	f := newFeverTest(t)

	rec := f.mark(url.Values{"mark": {"group"}, "as": {"read"}, "id": {"3"}})

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Result().StatusCode)
	}
	if f.storage.markedRead == nil || f.storage.markedRead.FolderID == nil {
		t.Fatal("expected MarkItemsRead with a FolderID")
	}
	if *f.storage.markedRead.FolderID != 3 {
		t.Fatalf("expected FolderID 3, got %d", *f.storage.markedRead.FolderID)
	}
}

func TestFeverMarkGroupReadAll(t *testing.T) {
	f := newFeverTest(t)

	rec := f.mark(url.Values{"mark": {"group"}, "as": {"read"}, "id": {"0"}})

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Result().StatusCode)
	}
	if f.storage.markedRead == nil {
		t.Fatal("expected MarkItemsRead to be invoked")
	}
	if f.storage.markedRead.FolderID != nil {
		t.Fatalf("expected no FolderID for id=0, got %d", *f.storage.markedRead.FolderID)
	}
}

func TestFeverMarkBefore(t *testing.T) {
	for _, mark := range []string{"feed", "group"} {
		t.Run(mark, func(t *testing.T) {
			f := newFeverTest(t)

			rec := f.mark(url.Values{
				"mark":   {mark},
				"as":     {"read"},
				"id":     {"1"},
				"before": {"1700000000"},
			})

			if rec.Result().StatusCode != http.StatusOK {
				t.Fatalf("expected 200, got %d", rec.Result().StatusCode)
			}
			if f.storage.markedRead == nil || f.storage.markedRead.Before == nil {
				t.Fatal("expected MarkItemsRead with a Before time")
			}
			want := time.Unix(1700000000, 0).UTC()
			if !f.storage.markedRead.Before.Equal(want) {
				t.Fatalf("expected Before %v, got %v", want, *f.storage.markedRead.Before)
			}
		})
	}
}

func TestFeverMarkNoBefore(t *testing.T) {
	for _, before := range []string{"0", "abc", ""} {
		t.Run(before, func(t *testing.T) {
			f := newFeverTest(t)

			rec := f.mark(url.Values{
				"mark":   {"group"},
				"as":     {"read"},
				"id":     {"1"},
				"before": {before},
			})

			if rec.Result().StatusCode != http.StatusOK {
				t.Fatalf("expected 200, got %d", rec.Result().StatusCode)
			}
			if f.storage.markedRead == nil {
				t.Fatal("expected MarkItemsRead to be invoked")
			}
			if f.storage.markedRead.Before != nil {
				t.Fatalf("expected no Before, got %v", *f.storage.markedRead.Before)
			}
		})
	}
}

func TestFeverMarkItem(t *testing.T) {
	cases := map[string]model.ItemStatus{
		"read":    model.READ,
		"unread":  model.UNREAD,
		"saved":   model.STARRED,
		"unsaved": model.READ,
	}
	for as, want := range cases {
		t.Run(as, func(t *testing.T) {
			f := newFeverTest(t)

			rec := f.mark(url.Values{"mark": {"item"}, "as": {as}, "id": {"42"}})

			if rec.Result().StatusCode != http.StatusOK {
				t.Fatalf("expected 200, got %d", rec.Result().StatusCode)
			}
			if f.storage.itemStatus == nil {
				t.Fatal("expected UpdateItem to be invoked")
			}
			if f.storage.itemStatusID != 42 {
				t.Fatalf("expected item id 42, got %d", f.storage.itemStatusID)
			}
			if *f.storage.itemStatus != want {
				t.Fatalf("expected status %v, got %v", want, *f.storage.itemStatus)
			}
		})
	}
}

func TestFeverMarkItemInvalidAs(t *testing.T) {
	f := newFeverTest(t)

	rec := f.mark(url.Values{"mark": {"item"}, "as": {"bogus"}, "id": {"42"}})

	if rec.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Result().StatusCode)
	}
	if f.storage.itemStatus != nil {
		t.Fatalf("rejected request had a side effect (status=%v)", *f.storage.itemStatus)
	}
}

func TestFeverMarkUnknownMark(t *testing.T) {
	f := newFeverTest(t)

	rec := f.mark(url.Values{"mark": {"foo"}, "as": {"read"}, "id": {"1"}})

	if rec.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Result().StatusCode)
	}
	if f.storage.markedRead != nil || f.storage.itemStatus != nil {
		t.Fatal("rejected request had a side effect")
	}
}

func TestFeverMarkInvalidID(t *testing.T) {
	f := newFeverTest(t)

	rec := f.mark(url.Values{"mark": {"group"}, "as": {"read"}, "id": {"abc"}})

	if rec.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Result().StatusCode)
	}
	if f.storage.markedRead != nil || f.storage.itemStatus != nil {
		t.Fatal("rejected request had a side effect")
	}
}
