package storage

type ItemStatus int

const (
	UNREAD  ItemStatus = 0
	READ    ItemStatus = 1
	STARRED ItemStatus = 2
)

type Item struct {
	Id int64
	FeedId int64
	Title string
	Link string
	Description string
	Content string
	Author string
	Date int64
	DateUpdated int64
	Status ItemStatus
	Image string
}

func (s *Storage) CreateItems(items []Item) bool {
	return true
}
