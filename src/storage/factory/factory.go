package factory

import (
	"github.com/nkanaev/yarr/src/storage/model"
)

type Storage interface {
	Close() error
	Migrate() error
	CountItems() int
	CreateFeed(params model.CreateFeedParams) *model.Feed
	CreateFolder(title string) *model.Folder
	CreateItems(items []model.Item) bool
	DeleteFeed(feedId int64) bool
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
	UpdateItemStatus(item_id int64, status model.ItemStatus) bool
	UpdateSettings(params model.UpdateSettingsParams) bool
}
