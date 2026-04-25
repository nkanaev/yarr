package storage

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"sort"
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

type MediaLink struct {
	URL         string `json:"url"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

type MediaLinks []MediaLink

func (m *MediaLinks) Scan(src any) error {
	switch data := src.(type) {
	case []byte:
		return json.Unmarshal(data, m)
	case string:
		return json.Unmarshal([]byte(data), m)
	default:
		return nil
	}
}

func (m MediaLinks) Value() (driver.Value, error) {
	return json.Marshal(m)
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

type ItemList []Item

func (list ItemList) Len() int {
	return len(list)
}

func (list ItemList) SortKey(i int) string {
	return list[i].Date.Format(time.RFC3339) + "::" + list[i].GUID
}

func (list ItemList) Less(i, j int) bool {
	return list.SortKey(i) < list.SortKey(j)
}

func (list ItemList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (s *Storage) CreateItems(items []Item) bool {
	tx, err := s.db.Begin()
	if err != nil {
		log.Print(err)
		return false
	}

	now := time.Now().UTC()

	itemsSorted := ItemList(items)
	sort.Sort(itemsSorted)

	for _, item := range itemsSorted {
		_, err = tx.Exec(`
			insert into items (
				guid, feed_id, title, link, date,
				content, media_links,
				date_arrived, status
			)
			values (
				:guid, :feed_id, :title, :link, strftime('%Y-%m-%d %H:%M:%f', :date),
				:content, :media_links,
				:date_arrived, :status
			)
			on conflict (feed_id, guid) do nothing`,
			sql.Named("guid", item.GUID),
			sql.Named("feed_id", item.FeedId),
			sql.Named("title", item.Title),
			sql.Named("link", item.Link),
			sql.Named("date", item.Date),
			sql.Named("content", item.Content),
			sql.Named("media_links", item.MediaLinks),
			sql.Named("date_arrived", now),
			sql.Named("status", UNREAD),
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
		cond = append(cond, "i.feed_id in (select id from feeds where folder_id = :folder_id)")
		args = append(args, sql.Named("folder_id", *filter.FolderID))
	}
	if filter.FeedID != nil {
		cond = append(cond, "i.feed_id = :feed_id")
		args = append(args, sql.Named("feed_id", *filter.FeedID))
	}
	if filter.Status != nil {
		cond = append(cond, "i.status = :status")
		args = append(args, sql.Named("status", *filter.Status))
	}
	if filter.Search != nil {
		words := strings.Fields(*filter.Search)
		terms := make([]string, len(words))
		for idx, word := range words {
			terms[idx] = word + "*"
		}

		cond = append(cond, "i.search_rowid in (select rowid from search where search match :search)")
		args = append(args, sql.Named("search", strings.Join(terms, " ")))
	}
	if filter.After != nil {
		compare := ">"
		if newestFirst {
			compare = "<"
		}
		cond = append(cond, fmt.Sprintf("(i.date, i.id) %s (select date, id from items where id = :after_id)", compare))
		args = append(args, sql.Named("after_id", *filter.After))
	}
	if filter.IDs != nil && len(*filter.IDs) > 0 {
		qmarks := make([]string, len(*filter.IDs))
		for i, id := range *filter.IDs {
			name := fmt.Sprintf("id%d", i)
			qmarks[i] = ":" + name
			args = append(args, sql.Named(name, id))
		}
		cond = append(cond, "i.id in ("+strings.Join(qmarks, ",")+")")
	}
	if filter.SinceID != nil {
		cond = append(cond, "i.id > :since_id")
		args = append(args, sql.Named("since_id", filter.SinceID))
	}
	if filter.MaxID != nil {
		cond = append(cond, "i.id < :max_id")
		args = append(args, sql.Named("max_id", filter.MaxID))
	}
	if filter.Before != nil {
		cond = append(cond, "i.date < :before")
		args = append(args, sql.Named("before", filter.Before))
	}

	predicate := "1"
	if len(cond) > 0 {
		predicate = strings.Join(cond, " and ")
	}

	return predicate, args
}

func (s *Storage) CountItems(filter ItemFilter) int {
	predicate, args := listQueryPredicate(filter, false)

	var count int
	query := fmt.Sprintf(`
		select count(*)
		from items
		where %s
		`, predicate)
	err := s.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		log.Print(err)
		return 0
	}
	return count
}

func (s *Storage) ListItems(filter ItemFilter, limit int, newestFirst bool, withContent bool) []Item {
	predicate, args := listQueryPredicate(filter, newestFirst)
	result := make([]Item, 0)

	order := "date desc, id desc"
	if !newestFirst {
		order = "date asc, id asc"
	}
	if filter.IDs != nil || filter.SinceID != nil {
		order = "i.id asc"
	}
	if filter.MaxID != nil {
		order = "i.id desc"
	}

	selectCols := "i.id, i.guid, i.feed_id, i.title, i.link, i.date, i.status, i.media_links"
	if withContent {
		selectCols += ", i.content"
	} else {
		selectCols += ", '' as content"
	}
	query := fmt.Sprintf(`
		select %s
		from items i
		where %s
		order by %s
		limit %d
		`, selectCols, predicate, order, limit)
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
			&x.Status, &x.MediaLinks, &x.Content,
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
			i.date, i.status, i.media_links
		from items i
		where i.id = :id
	`, sql.Named("id", id)).Scan(
		&i.Id, &i.GUID, &i.FeedId, &i.Title, &i.Link, &i.Content,
		&i.Date, &i.Status, &i.MediaLinks,
	)
	if err != nil {
		log.Print(err)
		return nil
	}
	return i
}

func (s *Storage) UpdateItemStatus(item_id int64, status ItemStatus) bool {
	_, err := s.db.Exec(`update items set status = :status where id = :id`,
		sql.Named("status", status),
		sql.Named("id", item_id),
	)
	return err == nil
}

func (s *Storage) MarkItemsRead(filter MarkFilter) bool {
	predicate, args := listQueryPredicate(ItemFilter{
		FolderID: filter.FolderID,
		FeedID:   filter.FeedID,
		Before:   filter.Before,
	}, false)
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
			insert into search (title, description, content) values (:title, "", :content)`,
			sql.Named("title", item.Title),
			sql.Named("content", htmlutil.ExtractText(item.Content)),
		)
		if err != nil {
			log.Print(err)
			return
		}
		if numrows, err := result.RowsAffected(); err == nil && numrows == 1 {
			if rowId, err := result.LastInsertId(); err == nil {
				s.db.Exec(
					`update items set search_rowid = :search_rowid where id = :id`,
					sql.Named("search_rowid", rowId),
					sql.Named("id", item.Id),
				)
			}
		}
	}
}

var (
	itemsKeepSize = 50
	itemsKeepDays = 90
)

// Delete old articles from the database to cleanup space.
//
// The rules:
//   - Never delete starred entries.
//   - Keep at least the same amount of articles the feed provides (default: 50).
//     This prevents from deleting items for rarely updated and/or ever-growing
//     feeds which might eventually reappear as unread.
//   - Keep entries for a certain period (default: 90 days).
func (s *Storage) DeleteOldItems() {
	rows, err := s.db.Query(`
		select
			i.feed_id,
			max(coalesce(s.size, 0), :keep_size) as max_items,
			count(*) as num_items
		from items i
		left outer join feed_sizes s on s.feed_id = i.feed_id
		where status != :starred_status
		group by i.feed_id
	`,
		sql.Named("keep_size", itemsKeepSize),
		sql.Named("starred_status", STARRED),
	)
	if err != nil {
		log.Print(err)
		return
	}

	feedLimits := make(map[int64]int64, 0)
	for rows.Next() {
		var feedId, limit int64
		rows.Scan(&feedId, &limit, nil)
		feedLimits[feedId] = limit
	}

	for feedId, limit := range feedLimits {
		result, err := s.db.Exec(`
			delete from items
			where id in (
				select i.id
				from items i
				where i.feed_id = :feed_id and status != :starred_status
				order by date desc
				limit -1 offset :limit
			) and date_arrived < :date_limit
			`,
			sql.Named("feed_id", feedId),
			sql.Named("starred_status", STARRED),
			sql.Named("limit", limit),
			sql.Named("date_limit", time.Now().UTC().Add(-time.Hour*time.Duration(24*itemsKeepDays))),
		)
		if err != nil {
			log.Print(err)
			return
		}
		numDeleted, err := result.RowsAffected()
		if err != nil {
			log.Print(err)
			return
		}
		if numDeleted > 0 {
			log.Printf("Deleted %d old items (feed: %d)", numDeleted, feedId)
		}
	}
}
