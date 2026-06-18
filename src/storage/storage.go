package storage

import (
	"strings"

	"github.com/nkanaev/yarr/src/storage/model"
	"github.com/nkanaev/yarr/src/storage/postgres"
	"github.com/nkanaev/yarr/src/storage/sqlite"
)

type Storage interface {
	Close() error
	CountItems() int
	CreateFeed(params model.CreateFeedParams) *model.Feed
	CreateFolder(title string) *model.Folder
	CreateItems(items []model.Item) bool
	DeleteFeed(feedId int64) bool
	DeleteItem(id int64) bool
	DeleteFolder(folderId int64) bool
	DeleteOldItems()
	FeedStats() []model.FeedStat
	GetFeed(id int64) *model.Feed
	GetFeedState(feedID int64) (*model.FeedState, error)
	GetItem(id int64) *model.Item
	GetSettings() model.Settings
	ListFeedStates() ([]model.FeedState, error)
	ListFeeds() []model.Feed
	ListFolders() []model.Folder
	ListItems(filter model.ItemFilter, limit int, newestFirst bool, withContent bool) []model.Item
	MarkItemsRead(filter model.MarkFilter) bool
	UpdateFeed(feedId int64, params model.UpdateFeedParams) (bool, error)
	UpdateFeedState(feedID int64, params model.UpdateFeedStateParams) (bool, error)
	UpdateFolder(folderId int64, params model.UpdateFolderParams) (bool, error)
	UpdateItem(id int64, params model.UpdateItemParams) bool
	UpdateItemStatus(item_id int64, status model.ItemStatus) bool
	UpdateSettings(params model.UpdateSettingsParams) bool
}

func New(path string) (Storage, error) {
	if strings.HasPrefix(path, "postgres://") || strings.HasPrefix(path, "postgresql://") {
		return postgres.New(path)
	}
	return sqlite.New(path)
}
