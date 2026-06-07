package storage

type IStorage interface {
	Close() error
	CountItems() int
	CreateFeed(params CreateFeedParams) *Feed
	CreateFolder(title string) *Folder
	CreateItems(items []Item) bool
	DeleteFeed(feedId int64) bool
	DeleteFolder(folderId int64) bool
	DeleteOldItems()
	FeedStats() []FeedStat
	GetFeed(id int64) *Feed
	GetFeedState(feedID int64) (*FeedState, error)
	GetItem(id int64) *Item
	GetSettings() Settings
	ListFeedStates() ([]FeedState, error)
	ListFeeds() []Feed
	ListFolders() []Folder
	ListItems(filter ItemFilter, limit int, newestFirst bool, withContent bool) []Item
	MarkItemsRead(filter MarkFilter) bool
	UpdateFeed(feedId int64, params UpdateFeedParams) (bool, error)
	UpdateFeedState(feedID int64, params UpdateFeedStateParams) (bool, error)
	UpdateFolder(folderId int64, params UpdateFolderParams) (bool, error)
	UpdateItemStatus(item_id int64, status ItemStatus) bool
	UpdateSettings(params UpdateSettingsParams) bool
}
