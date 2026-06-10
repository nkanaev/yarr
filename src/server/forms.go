package server

import "github.com/nkanaev/yarr/src/storage/model"

type ItemUpdateForm struct {
	Status *model.ItemStatus `json:"status,omitempty"`
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
