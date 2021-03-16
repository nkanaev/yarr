package server

import "github.com/nkanaev/yarr/src/storage"

type ItemUpdateForm struct {
	Status *storage.ItemStatus `json:"status,omitempty"`
}

type FolderCreateForm struct {
	Title string `json:"title"`
}

type FolderUpdateForm struct {
	Title      *string `json:"title,omitempty"`
	IsExpanded *bool   `json:"is_expanded,omitempty"`
}

type FeedCreateForm struct {
	Url      string `json:"url"`
	FolderID *int64 `json:"folder_id,omitempty"`
}
