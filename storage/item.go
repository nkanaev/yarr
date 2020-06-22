package storage

type ItemStatus int

const (
	UNREAD  ItemStatus = 0
	READ    ItemStatus = 1
	STARRED ItemStatus = 2
)

type Item struct {
	Id string
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
	tx, err := s.db.Begin()
	if err != nil {
		s.log.Print(err)
		return false
	}
	for _, item := range items {
		_, err = tx.Exec(`
			insert into items (
				id, feed_id, title, link, description,
				content, author, date, date_updated, status, image
			)
			values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			on conflict (id) do update set date_updated=?`,
			item.Id, item.FeedId, item.Title, item.Link, item.Description,
			item.Content, item.Author, item.Date, item.DateUpdated, UNREAD, item.Image,
			// upsert values
			item.DateUpdated,
		)
		if err != nil {
			s.log.Print(err)
			if err = tx.Rollback(); err != nil {
				s.log.Print(err)
				return false
			}
			return false
		}
	}
	if err = tx.Commit(); err != nil {
		s.log.Print(err)
		return false
	}
	return true
}
