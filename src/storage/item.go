package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nkanaev/yarr/src/content/htmlutil"
)

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

type Item struct {
	Id       int64      `json:"id"`
	GUID     string     `json:"guid"`
	FeedId   int64      `json:"feed_id"`
	Title    string     `json:"title"`
	Link     string     `json:"link"`
	Content  string     `json:"content,omitempty"`
	Date     time.Time  `json:"date"`
	Status   ItemStatus `json:"status"`
	ImageURL *string    `json:"image"`
	AudioURL *string    `json:"podcast_url"`
}

type ItemFilter struct {
	FolderID *int64
	FeedID   *int64
	Status   *ItemStatus
	Search   *string
	After  *int64
}

type MarkFilter struct {
	FolderID *int64
	FeedID   *int64
}

func (s *Storage) CreateItems(items []Item) bool {
	tx, err := s.db.Begin()
	if err != nil {
		log.Print(err)
		return false
	}

	now := time.Now()

	for _, item := range items {
		_, err = tx.Exec(`
			insert into items (
				guid, feed_id, title, link, date,
				content, image, podcast_url,
				date_arrived, status
			)
			values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			on conflict (feed_id, guid) do nothing`,
			item.GUID, item.FeedId, item.Title, item.Link, item.Date,
			item.Content, item.ImageURL, item.AudioURL,
			now, UNREAD,
		)
		if err != nil {
			log.Print(err)
			if err = tx.Rollback(); err != nil {
				log.Print(err)
				return false
			}
			return false
		}
	}
	if err = tx.Commit(); err != nil {
		log.Print(err)
		return false
	}
	return true
}

func listQueryPredicate(filter ItemFilter, newestFirst bool) (string, []interface{}) {
	cond := make([]string, 0)
	args := make([]interface{}, 0)
	if filter.FolderID != nil {
		cond = append(cond, "i.feed_id in (select id from feeds where folder_id = ?)")
		args = append(args, *filter.FolderID)
	}
	if filter.FeedID != nil {
		cond = append(cond, "i.feed_id = ?")
		args = append(args, *filter.FeedID)
	}
	if filter.Status != nil {
		cond = append(cond, "i.status = ?")
		args = append(args, *filter.Status)
	}
	if filter.Search != nil {
		words := strings.Fields(*filter.Search)
		terms := make([]string, len(words))
		for idx, word := range words {
			terms[idx] = word + "*"
		}

		cond = append(cond, "i.search_rowid in (select rowid from search where search match ?)")
		args = append(args, strings.Join(terms, " "))
	}
	if filter.After != nil {
		compare := ">"
		if newestFirst {
			compare = "<"
		}
		cond = append(cond, fmt.Sprintf("(i.date, i.id) %s (select date, id from items where id = ?)", compare))
		args = append(args, *filter.After)
	}

	predicate := "1"
	if len(cond) > 0 {
		predicate = strings.Join(cond, " and ")
	}

	return predicate, args
}

func (s *Storage) ListItems(filter ItemFilter, limit int, newestFirst bool) []Item {
	predicate, args := listQueryPredicate(filter, newestFirst)
	result := make([]Item, 0, 0)

	order := "date desc, id desc"
	if !newestFirst {
		order = "date asc, id asc"
	}

	query := fmt.Sprintf(`
		select
			i.id, i.guid, i.feed_id,
			i.title, i.link, i.date,
			i.status, i.image, i.podcast_url
		from items i
		where %s
		order by %s
		limit %d
		`, predicate, order, limit)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		log.Print(err)
		return result
	}
	for rows.Next() {
		var x Item
		err = rows.Scan(
			&x.Id, &x.GUID, &x.FeedId,
			&x.Title, &x.Link, &x.Date,
			&x.Status, &x.ImageURL, &x.AudioURL,
		)
		if err != nil {
			log.Print(err)
			return result
		}
		result = append(result, x)
	}
	return result
}

func (s *Storage) GetItem(id int64) *Item {
	i := &Item{}
	err := s.db.QueryRow(`
		select
			i.id, i.guid, i.feed_id, i.title, i.link, i.content,
			i.date, i.status, i.image, i.podcast_url
		from items i
		where i.id = ?
	`, id).Scan(
		&i.Id, &i.GUID, &i.FeedId, &i.Title, &i.Link, &i.Content,
		&i.Date, &i.Status, &i.ImageURL, &i.AudioURL,
	)
	if err != nil {
		log.Print(err)
		return nil
	}
	return i
}

func (s *Storage) UpdateItemStatus(item_id int64, status ItemStatus) bool {
	_, err := s.db.Exec(`update items set status = ? where id = ?`, status, item_id)
	return err == nil
}

func (s *Storage) MarkItemsRead(filter MarkFilter) bool {
	predicate, args := listQueryPredicate(ItemFilter{FolderID: filter.FolderID, FeedID: filter.FeedID}, false)
	query := fmt.Sprintf(`
		update items as i set status = %d
		where %s and i.status != %d
		`, READ, predicate, STARRED)
	_, err := s.db.Exec(query, args...)
	if err != nil {
		log.Print(err)
	}
	return err == nil
}

type FeedStat struct {
	FeedId       int64 `json:"feed_id"`
	UnreadCount  int64 `json:"unread"`
	StarredCount int64 `json:"starred"`
}

func (s *Storage) FeedStats() []FeedStat {
	result := make([]FeedStat, 0)
	rows, err := s.db.Query(fmt.Sprintf(`
		select
			feed_id,
			sum(case status when %d then 1 else 0 end),
			sum(case status when %d then 1 else 0 end)
		from items
		group by feed_id
	`, UNREAD, STARRED))
	if err != nil {
		log.Print(err)
		return result
	}
	for rows.Next() {
		stat := FeedStat{}
		rows.Scan(&stat.FeedId, &stat.UnreadCount, &stat.StarredCount)
		result = append(result, stat)
	}
	return result
}

func (s *Storage) SyncSearch() {
	rows, err := s.db.Query(`
		select id, title, content
		from items
		where search_rowid is null;
	`)
	if err != nil {
		log.Print(err)
		return
	}

	items := make([]Item, 0)
	for rows.Next() {
		var item Item
		rows.Scan(&item.Id, &item.Title, &item.Content)
		items = append(items, item)
	}

	for _, item := range items {
		result, err := s.db.Exec(`
			insert into search (title, description, content) values (?, "", ?)`,
			item.Title, htmlutil.ExtractText(item.Content),
		)
		if err != nil {
			log.Print(err)
			return
		}
		if numrows, err := result.RowsAffected(); err == nil && numrows == 1 {
			if rowId, err := result.LastInsertId(); err == nil {
				s.db.Exec(
					`update items set search_rowid = ? where id = ?`,
					rowId, item.Id,
				)
			}
		}
	}
}


// TODO: better naming
var (
	itemsKeepSize = 100
	itemsKeepDays = 90
)

func (s *Storage) DeleteOldItems() {
	rows, err := s.db.Query(fmt.Sprintf(`
		select feed_id, count(*) as num_items
		from items
		where status != %d
		group by feed_id
		having num_items > 50
	`, STARRED))

	if err != nil {
		log.Print(err)
		return
	}

	feedIds := make([]int64, 0)
	for rows.Next() {
		var id int64
		rows.Scan(&id, nil)
		feedIds = append(feedIds, id)
	}

	for _, feedId := range feedIds {
		result, err := s.db.Exec(`
			delete from items where feed_id = ? and status != ? and date_arrived < ?`,
			feedId,
			STARRED,
			time.Now().Add(-time.Hour*time.Duration(24*itemsKeepDays)),
		)
		if err != nil {
			log.Print(err)
			return
		}
		num, err := result.RowsAffected()
		if err != nil {
			log.Print(err)
			return
		}
		if num > 0 {
			log.Printf("Deleted %d old items (%d)", num, feedId)
		}
	}
}
