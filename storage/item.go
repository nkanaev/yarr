package storage

import (
	"fmt"
	"time"
	"encoding/json"
)

type ItemStatus int

const (
	UNREAD  ItemStatus = 0
	READ    ItemStatus = 1
	STARRED ItemStatus = 2
)

var StatusRepresentations = map[ItemStatus]string {
	UNREAD: "unread",
	READ: "read",
	STARRED: "starred",
}

var StatusValues = map[string]ItemStatus {
	"unread": UNREAD,
	"read": READ,
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

type Item struct {
	Id int64 `json:"id"`
	GUID string `json:"guid"`
	FeedId int64 `json:"feed_id"`
	Title string `json:"title"`
	Link string `json:"link"`
	Description string `json:"description"`
	Content string `json:"content"`
	Author string `json:"author"`
	Date *time.Time `json:"date"`
	DateUpdated *time.Time `json:"date_updated"`
	Status ItemStatus `json:"status"`
	Image string `json:"image"`
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
				guid, feed_id, title, link, description,
				content, author, date, date_updated, status, image
			)
			values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			on conflict (guid) do update set date_updated=?`,
			item.GUID, item.FeedId, item.Title, item.Link, item.Description,
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

func itemQuery(s *Storage, cond string, v ...interface{}) []Item {
	result := make([]Item, 0, 0)
	query := fmt.Sprintf(`
		select
			id, guid, feed_id, title, link, description,
			content, author, date, date_updated, status, image
		from items
		where %s`, cond)
	rows, err := s.db.Query(query, v...)
	if err != nil {
		s.log.Print(err)
		return result
	}
	for rows.Next() {
		var x Item
		err = rows.Scan(
			&x.Id,
			&x.GUID,
			&x.FeedId,
			&x.Title,
			&x.Link,
			&x.Description,
			&x.Content,
			&x.Author,
			&x.Date,
			&x.DateUpdated,
			&x.Status,
			&x.Image,
		)
		if err != nil {
			s.log.Print(err)
			return result
		}
		result = append(result, x)
	}
	return result
}

func (s *Storage) ListItems() []Item {
	return itemQuery(s, `1`)
}

func (s *Storage) ListFolderItems(folder_id int64) []Item {
	return itemQuery(s, `folder_id = ?`, folder_id)
}

func (s *Storage) ListFolderItemsFiltered(folder_id int64, status ItemStatus) []Item {
	return itemQuery(s, `folder_id = ? and status = ?`, folder_id, status)
}

func (s *Storage) ListFeedItems(feed_id int64) []Item {
	return itemQuery(s, `feed_id = ?`, feed_id)
}

func (s *Storage) ListFeedItemsFiltered(feed_id int64, status ItemStatus) []Item {
	return itemQuery(s, `feed_id = ? and status = ?`, feed_id, status)
}

func (s *Storage) UpdateItemStatus(item_id int64, status ItemStatus) bool {
	_, err := s.db.Exec(`update items set status = ? where id = ?`, status, item_id)
	return err == nil
}
