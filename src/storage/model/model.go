package model

import (
	"encoding/json"
	"time"
)

type Feed struct {
	Id          int64   `json:"id"`
	FolderId    *int64  `json:"folder_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Link        string  `json:"link"`
	FeedLink    string  `json:"feed_link"`
	Icon        *[]byte `json:"icon,omitempty"`
	HasIcon     bool    `json:"has_icon"`
}

type CreateFeedParams struct {
	Title       string
	Description string
	Link        string
	FeedLink    string
	FolderID    *int64
}

type Item struct {
	Id         int64      `json:"id"`
	GUID       string     `json:"guid"`
	FeedId     int64      `json:"feed_id"`
	Title      string     `json:"title"`
	Link       string     `json:"link"`
	Content    string     `json:"content,omitempty"`
	Date       time.Time  `json:"date"`
	Status     ItemStatus `json:"status"`
	MediaLinks MediaLinks `json:"media_links"`
}

type ItemStatus int

const (
	UNREAD  ItemStatus = 0
	READ    ItemStatus = 1
	STARRED ItemStatus = 2
)


var StatusRepresentations = map[ItemStatus]string{
	UNREAD:  "unread",
	READ:    "read",
	STARRED: "starred",
}

var StatusValues = map[string]ItemStatus{
	"unread":  UNREAD,
	"read":    READ,
	"starred": STARRED,
}

func (s ItemStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(StatusRepresentations[s])
}

func (s *ItemStatus) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = StatusValues[str]
	return nil
}

type MediaLink struct {
	URL         string `json:"url"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

type MediaLinks []MediaLink

type ItemFilter struct {
	FolderID *int64
	FeedID   *int64
	Status   *ItemStatus
	Search   *string
	After    *int64
	IDs      *[]int64
	SinceID  *int64
	MaxID    *int64
	Before   *time.Time
}

type MarkFilter struct {
	FolderID *int64
	FeedID   *int64

	Before *time.Time
}

type Folder struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	IsExpanded bool   `json:"is_expanded"`
}

type UpdateFolderParams struct {
	Title      *string
	IsExpanded *bool
}

type FeedStat struct {
	FeedId       int64 `json:"feed_id"`
	UnreadCount  int64 `json:"unread"`
	StarredCount int64 `json:"starred"`
}


type Settings struct {
	Filter          string `json:"filter"`
	Feed            string `json:"feed"`
	FeedListWidth   int    `json:"feed_list_width"`
	ItemListWidth   int    `json:"item_list_width"`
	SortNewestFirst bool   `json:"sort_newest_first"`
	ThemeName       string `json:"theme_name"`
	ThemeFont       string `json:"theme_font"`
	ThemeSize       int    `json:"theme_size"`
	RefreshRate     int64  `json:"refresh_rate"`
	Language        string `json:"language"`
}

type UpdateSettingsParams struct {
	Filter          *string `json:"filter"`
	Feed            *string `json:"feed"`
	FeedListWidth   *int    `json:"feed_list_width"`
	ItemListWidth   *int    `json:"item_list_width"`
	SortNewestFirst *bool   `json:"sort_newest_first"`
	ThemeName       *string `json:"theme_name"`
	ThemeFont       *string `json:"theme_font"`
	ThemeSize       *int    `json:"theme_size"`
	RefreshRate     *int64  `json:"refresh_rate"`
	Language        *string `json:"language"`
}

func (s Settings) Map() map[string]any {
	return map[string]any{
		"filter":            s.Filter,
		"feed":              s.Feed,
		"feed_list_width":   s.FeedListWidth,
		"item_list_width":   s.ItemListWidth,
		"sort_newest_first": s.SortNewestFirst,
		"theme_name":        s.ThemeName,
		"theme_font":        s.ThemeFont,
		"theme_size":        s.ThemeSize,
		"refresh_rate":      s.RefreshRate,
		"language":          s.Language,
	}
}

type FeedState struct {
	FeedID           int64
	LastRefreshed    time.Time
	LastError        string
	HTTPLastModified string
	HTTPEtag         string
}

type UpdateFeedStateParams struct {
	LastRefreshed    *time.Time
	LastError        *string
	HTTPLastModified *string
	HTTPEtag         *string
}

type UpdateFeedParams struct {
	Title    *string
	FeedLink *string
	FolderID Nullable[int64]
	Icon     Nullable[[]byte]
}

type Nullable[T any] struct {
	Set   bool
	Value *T
}

func SetNullable[T any](v *T) Nullable[T] {
	return Nullable[T]{Set: true, Value: v}
}
